package srv

import (
	"fmt"
	"net/http"

	"github.com/gluke77/mbapi/message"
	syncsender "github.com/gluke77/mbapi/sender/sync"
)

type Sender interface {
	Send(msg message.Message) error
}

type Server struct {
	Sender
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%d -- Method is not supported: %s", http.StatusBadRequest, r.Method)
		return
	}

	msg, err := message.NewFrom(r.Body)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing the request: %s", err), http.StatusBadRequest)
		return
	}

	if err := msg.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Error processing the request: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.Send(msg); err != nil {
		http.Error(w, fmt.Sprintf("Error sending the message: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func New(apiKey string) *Server {
	sender := syncsender.New(apiKey)
	return &Server{Sender: sender}
}
