package messaging

import "go.nanomsg.org/mangos/v3/protocol"

type Bot interface {
	SendAlertMessage(msg []byte)
	Start()
	SubscribeToAllAlerts() error
	SetPubSubSocket(pubSubSocket protocol.Socket)
}
