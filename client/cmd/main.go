package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/evanisnor/gomessaging/server/pkg/api"
)

const (
	serverAddress = "localhost:5558"
	exitCommand   = "!q"
	senderID      = "clientman"
)

func main() {
	client := connect()
	ctx := context.Background()

	register(ctx, &client)
	s := connectMessagerStream(ctx, &client)
	go readMessages(ctx, &s)

	readInput(func(text string) {
		message := createMessage(text)
		sendMessage(ctx, &s, message)
	}, func() {
		s.CloseSend()
	})
}

func connect() api.MessagingClient {
	log.Println("Connecting to", serverAddress)
	connection, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}
	client := api.NewMessagingClient(connection)
	return client
}

func connectMessagerStream(ctx context.Context, client *api.MessagingClient) api.Messaging_MessagerClient {
	log.Println("Connecting to messager stream")
	c := *client
	stream, err := c.Messager(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to Messager stream:\n%v", err)
	}
	return stream
}

func register(ctx context.Context, client *api.MessagingClient) {
	log.Println("Registering sender")

	c := *client
	_, err := c.Register(ctx, &api.Registration{
		SenderId: senderID,
	})
	if err != nil {
		log.Fatalf("Failed to register sender: %v", err)
	}

	log.Printf("Connected as %s\n", senderID)
}

func createMessage(text string) *api.Message {
	now, _ := ptypes.TimestampProto(time.Now())
	return &api.Message{
		Text:      text,
		SenderId:  senderID,
		Timestamp: now,
	}
}

func readInput(handler func(text string), exitHandler func()) {
	log.Println("Chat Ready")
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for text != exitCommand {
		scanner.Scan()
		text = scanner.Text()
		if text != exitCommand {
			handler(text)
		}
	}

	exitHandler()
}

func sendMessage(ctx context.Context, stream *api.Messaging_MessagerClient, message *api.Message) {
	s := *stream
	s.Send(message)
}

func readMessages(ctx context.Context, stream *api.Messaging_MessagerClient) {
	for {
		s := *stream
		in, err := s.Recv()
		if err == io.EOF {
			log.Println("Server closed the messager stream")
		} else if err != nil {
			log.Fatalf("Error receiving from messager stream: %v", err)
		}

		log.Printf("[%s] %s\n", in.SenderId, in.Text)
	}
}
