package sender

import (
	"errors"
	"fmt"
	"log"
	"math"
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

type SyncSender struct {
	client      *messagebird.Client
	rateLimiter <-chan time.Time
}

func (s *SyncSender) send_part(originator string, recipients []string, text string, partNo byte, totalParts byte, ref byte) error {
	<-s.rateLimiter
	s.rateLimiter = time.After(time.Second)

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

	mbMsg, err := s.client.NewMessage(originator, recipients, text, msgParams)

	if err != nil {
		errStr := fmt.Sprintf("Error sending message: %s:", err)
		if mbMsg != nil {
			errStr += fmt.Sprintf("\n%v", mbMsg.Errors)
		}

		return errors.New(errStr)
	}

	return nil
}

func (s *SyncSender) Send(msg message.Message) error {
	recipients := []string{fmt.Sprintf("%d", msg.Recipient)}
	ref := byte(time.Now().Second())

	length := len(msg.Message)

	if length <= MaxSinglePartMessageSize {
		return s.send_part(msg.Originator, recipients, msg.Message, 1, 1, ref)
	}

	if length > MaxMessageSize {
		return errors.New("Message too long")
	}

	totalParts := byte(math.Ceil(float64(length) / float64(PartSize)))
	var partNo byte
	for partNo = 1; partNo <= totalParts; partNo++ {
		start := int(partNo-1) * PartSize
		end := start + PartSize
		log.Printf("start=%d, end=%d", start, end)

		if end > length {
			end = length
		}

		txt := msg.Message[start:end]

		if err := s.send_part(msg.Originator, recipients, txt, partNo, totalParts, ref); err != nil {
			return err
		}
	}

	return nil
}

func New(apiKey string) *SyncSender {
	client := messagebird.New(apiKey)
	client.DebugLog = log.New(os.Stderr, "client", 0)

	return &SyncSender{client: client, rateLimiter: time.After(time.Second)}
}
