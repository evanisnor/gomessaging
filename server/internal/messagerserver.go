package internal

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	"github.com/evanisnor/gomessaging/server/pkg/api"
)

const (
	senderID = "serverman"
)

// MessagerServer ok
type MessagerServer struct {
	senders []string
}

// Register ok
func (s *MessagerServer) Register(ctx context.Context, in *api.Registration) (*empty.Empty, error) {
	log.Printf("Registered sender: %s", in.GetSenderId())
	s.registerSender(in.GetSenderId())
	return &empty.Empty{}, nil
}

// Messager ok
func (s *MessagerServer) Messager(stream api.Messaging_MessagerServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			log.Println("Sender dropped")
			return nil
		} else if err != nil {
			log.Printf("Stream error: %v\n", err)
			return err
		}

		log.Printf("[%s] %s", msg.GetSenderId(), msg.GetText())
		now, _ := ptypes.TimestampProto(time.Now())
		sent := ptypes.TimestampString(msg.GetTimestamp())
		stream.Send(&api.Message{
			Text:      fmt.Sprintf("hello %s, you said \"%s\" at %v", msg.GetSenderId(), msg.GetText(), sent),
			Timestamp: now,
			SenderId:  senderID,
		})
	}
}

func (s *MessagerServer) registerSender(senderID string) {
	if s.senders == nil {
		s.senders = make([]string, 0, 0)
	}

	s.senders = append(s.senders, senderID)
}
