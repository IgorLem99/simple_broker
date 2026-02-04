package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"broker/internal/broker"
)

type Handler struct {
	broker *broker.Broker
}

func New(b *broker.Broker) *Handler {
	return &Handler{broker: b}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.NotFound(w, r)
		return
	}

	queueName := parts[3]
	action := parts[4]

	switch {
	case r.Method == http.MethodPost && action == "messages":
		h.postMessage(w, r, queueName)
	case r.Method == http.MethodPost && action == "subscriptions":
		h.postSubscription(w, r, queueName)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request, queueName string) {
	q, err := h.broker.GetQueue(queueName)
	if err != nil {
		if err == broker.ErrQueueNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var msg broker.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := q.Send(msg); err != nil {
		if err == broker.ErrQueueFull {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) postSubscription(w http.ResponseWriter, r *http.Request, queueName string) {
	q, err := h.broker.GetQueue(queueName)
	if err != nil {
		if err == broker.ErrQueueNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sub, err := q.Subscribe()
	if err != nil {
		if err == broker.ErrTooManySub {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for msg := range sub {
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			q.Unsubscribe(sub)
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
