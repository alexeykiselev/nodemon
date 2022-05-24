package pubsub

import (
	"context"
	"log"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
	"nodemon/pkg/entities"
	"nodemon/pkg/messaging"
)

func subscribeToAlerts(socket protocol.Socket) error {
	for alertType := range entities.AlertTypes {
		err := socket.SetOption(mangos.OptionSubscribe, []byte{byte(alertType)})
		if err != nil {
			return err
		}
	}
	return nil
}

func StartMessagingClient(ctx context.Context, nanomsgURL string, bot messaging.Bot) error {
	socket, err := sub.NewSocket()
	if err != nil {
		log.Printf("failed to get new sub socket: %v", err)
		return err
	}
	defer func(socketPair protocol.Socket) {
		if err := socketPair.Close(); err != nil {
			log.Printf("Failed to close pair socket: %v", err)
		}
	}(socket)
	if err := socket.Dial(nanomsgURL); err != nil {
		log.Printf("failed to dial on sub socket: %v", err)
		return err
	}
	err = subscribeToAlerts(socket)
	if err != nil {
		log.Printf("failed to subscribe on empty topic: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := socket.Recv()
				if err != nil {
					log.Printf("failed to receive message: %v", err)
					return
				}
				bot.SendMessage(msg)
			}
		}
	}()

	<-ctx.Done()
	log.Println("pubsub messaging service finished")
	return nil
}
