package async

import (
	"errors"
	"log"

	"github.com/gluke77/mbapi/message"
	"github.com/gluke77/mbapi/sender/sync"
)

var (
	ErrRateLimited = errors.New("Ratelimited")
)

type Sender struct {
	syncsender *sync.Sender
	messages   chan message.Message
}

func (s *Sender) Send(msg message.Message) error {
	select {
	case s.messages <- msg:
		return nil
	default:
		return ErrRateLimited
	}
	return nil
}

func (s *Sender) send() {
	for m := range s.messages {
		if err := s.syncsender.Send(m); err != nil {
			log.Println(err)
		}
	}
}

func (s *Sender) Close() {
	close(s.messages)
}

func New(apiKey string, queueSize int) *Sender {
	s := &Sender{syncsender: sync.New(apiKey), messages: make(chan message.Message, queueSize)}
	go s.send()
	return s
}
