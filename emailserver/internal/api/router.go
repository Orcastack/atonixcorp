package api

import (
	"net/http"
)

type Dependencies struct {
	Store     Store
	Blob      BlobStore
	Auth      AuthService
	Domain    DomainService
	Queue     QueueService
	Analytics AnalyticsClient
	Logger    Logger
}

type Store interface {
	ListInbox(userID string) ([]Message, error)
	GetMessage(userID, id string) (*Message, error)
	DeleteMessage(userID, id string) error
}

type BlobStore interface {
	GetAttachment(userID, msgID, attachmentID string) ([]byte, error)
}

type AuthService interface {
	AuthenticateToken(token string) (userID string, err error)
}

type DomainService interface {
	RegisterDomain(name string) error
	ListDomains() ([]Domain, error)
}

type QueueService interface {
	EnqueueSend(msg OutgoingMessage) error
}

type AnalyticsClient interface {
	SendEvent(kind string, payload map[string]any)
}

type Logger interface {
	Info(...any)
	Error(...any)
}

type Message struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
}

type Domain struct {
	Name string `json:"name"`
}

type OutgoingMessage struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func NewRouter(dep Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/mail/send", dep.withAuth(dep.handleSendMail))
	mux.HandleFunc("/mail/inbox", dep.withAuth(dep.handleInbox))
	mux.HandleFunc("/mail/message", dep.withAuth(dep.handleGetMessage))
	mux.HandleFunc("/mail/delete", dep.withAuth(dep.handleDeleteMessage))

	mux.HandleFunc("/domain/register", dep.withAuth(dep.handleRegisterDomain))
	mux.HandleFunc("/domain/list", dep.withAuth(dep.handleListDomains))

	return mux
}

func (d Dependencies) withAuth(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := d.Auth.AuthenticateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r, userID)
	}
}
