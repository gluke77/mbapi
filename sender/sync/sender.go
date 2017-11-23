package sync

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	messagebird "github.com/messagebird/go-rest-api"

	"github.com/gluke77/mbapi/message"
)

const (
	MaxSinglePartMessageSize = 160
	PartSize                 = 153
	MaxMessageSize           = PartSize * 255
)

type Sender struct {
	client      *messagebird.Client
	rateLimiter <-chan time.Time
}

func split(m message.Message) []message.Message {
	messages := []message.Message{}

	text := m.Message

	if len(text) <= MaxSinglePartMessageSize {
		messages := append(messages, m)
		return messages
	}

	for len(text) > 0 {
		end := PartSize
		if end > len(text) {
			end = len(text)
		}

		txt := text[:end]
		text = text[end:]

		msg := message.Message{
			Originator: m.Originator,
			Recipient:  m.Recipient,
			Message:    txt,
		}

		messages = append(messages, msg)
	}

	return messages
}

func (s *Sender) send_part(m message.Message, partNo byte, totalParts byte, ref byte) error {
	<-s.rateLimiter

	var msgParams *messagebird.MessageParams

	if totalParts > 1 {
		udh := fmt.Sprintf("%02d%02d%02d%02d%02d%02d", 5, 0, 3, ref, totalParts, partNo)

		//udh[0] = 5 // Length of User Data Handler
		//udh[1] = 0 // Information Element Identifier, equal to 00 (Concatenated short messages, 8-bit reference number)
		//udh[2] = 3 // Length of the header, excluding the first two fields
		//udh[3] = ref        // CSMS reference number
		//udh[4] = totalParts // Total number of parts
		//udh[5] = partNo

		typeDetails := messagebird.TypeDetails{"udh": udh}
		msgParams = &messagebird.MessageParams{Type: "binary", TypeDetails: typeDetails, DataCoding: "auto"}
	}

	recipients := []string{m.Recipient}
	mbMsg, err := s.client.NewMessage(m.Originator, recipients, m.Message, msgParams)

	if err != nil {
		errStr := fmt.Sprintf("Error sending message: %s:", err)
		if mbMsg != nil {
			errStr += fmt.Sprintf("\n%v", mbMsg.Errors)
		}

		return errors.New(errStr)
	}

	return nil
}

func (s *Sender) Send(msg message.Message) error {
	if len(msg.Message) > MaxMessageSize {
		return errors.New("Message too long")
	}

	messages := split(msg)
	totalParts := byte(len(messages))
	ref := byte(time.Now().Second())

	for partNo, m := range messages {
		if err := s.send_part(m, byte(partNo+1), totalParts, ref); err != nil {
			return err
		}
	}

	return nil
}

func New(apiKey string) *Sender {
	client := messagebird.New(apiKey)
	client.DebugLog = log.New(os.Stderr, "client", 0)

	return &Sender{client: client, rateLimiter: time.Tick(time.Second)}
}
