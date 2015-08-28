package nano

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
)

func NewLocalRepository() ThingeyRepository {
	receiver, remoteSender := libchan.Pipe()
	remoteReceiver, sender := libchan.Pipe()
	repo := NewThingeyRepository(sender, receiver, remoteSender)
	adapter := NewThingeyAdapter(remoteReceiver)
	go func() {
		for {
			adapter.Listen()
		}
	}()
	return repo
}

func NewRemoteRepository(remoteURL string) ThingeyRepository {
	receiver, remoteSender := libchan.Pipe()
	var client net.Conn
	var err error
	client, err = tls.Dial("tcp", remoteURL, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	transport, err := spdy.NewClientTransport(client)
	if err != nil {
		log.Fatal(err)
	}
	sender, err := transport.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}
	return NewThingeyRepository(sender, receiver, remoteSender)
}
