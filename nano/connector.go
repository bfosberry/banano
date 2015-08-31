package nano

import (
	"log"
	"net"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
)

func NewLocalRepository() ThingeyRepository {
	receiver, remoteSender := libchan.Pipe()
	remoteReceiver, sender := libchan.Pipe()
	senderFunc := func() (libchan.Sender, error) {
		return sender, nil
	}
	repo := NewThingeyRepository(senderFunc, receiver, remoteSender)
	adapter := NewThingeyAdapter()
	go func() {
		for {
			adapter.Listen(remoteReceiver)
		}
	}()
	return repo
}

func NewRemoteRepository(remoteURL string) ThingeyRepository {
	receiver, remoteSender := libchan.Pipe()
	var client net.Conn
	var err error
	client, err = net.Dial("tcp", remoteURL)
	if err != nil {
		log.Fatal(err)
	}

	transport, err := spdy.NewClientTransport(client)
	if err != nil {
		log.Fatal(err)
	}
	return NewThingeyRepository(transport.NewSendChannel, receiver, remoteSender)
}
