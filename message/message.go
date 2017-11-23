package message

import (
	"encoding/json"
	"errors"
	"io"
)

type Message struct {
	Recipient  int64  `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

func (m Message) Validate() error {
	err := ""
	if m.Recipient == 0 {
		//TODO: Come up with a better validation
		err += "empty recipient"
	}

	if len(m.Originator) == 0 {
		if len(err) > 0 {
			err += "; "
		}
		err += "empty originator"
	}

	if len(m.Message) == 0 {
		if len(err) > 0 {
			err += "; "
		}
		err += "empty message"
	}

	if len(err) > 0 {
		return errors.New(err)
	}

	return nil
}

func NewFrom(r io.Reader) (Message, error) {
	msg := Message{}

	if err := json.NewDecoder(r).Decode(&msg); err != nil {
		return msg, err
	}

	return msg, nil
}
