package api

import (
	"encoding/json"
	"net/http"
)

func (d Dependencies) handleSendMail(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req OutgoingMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.From == "" {
		req.From = userID + "@example.com" // or your domain logic
	}

	if err := d.Queue.EnqueueSend(req); err != nil {
		d.Logger.Error("enqueue send:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d.Analytics.SendEvent("email_send", map[string]any{
		"user":    userID,
		"to":      req.To,
		"subject": req.Subject,
	})

	w.WriteHeader(http.StatusAccepted)
}

func (d Dependencies) handleInbox(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msgs, err := d.Store.ListInbox(userID)
	if err != nil {
		d.Logger.Error("list inbox:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func (d Dependencies) handleGetMessage(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := d.Store.GetMessage(userID, id)
	if err != nil {
		d.Logger.Error("get message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if msg == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (d Dependencies) handleDeleteMessage(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := d.Store.DeleteMessage(userID, id); err != nil {
		d.Logger.Error("delete message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
