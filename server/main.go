package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bfosberry/banano/nano"
	"github.com/docker/libchan/spdy"
)

func main() {
	log.Println("Starting Server..")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening")

	tl, err := spdy.NewTransportListener(listener, spdy.NoAuthenticator)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("Accept loop")
		t, err := tl.AcceptTransport()
		if err != nil {
			log.Print(err)
			break
		}

		go func() {
			for {
				fmt.Println("Receive loop")
				receiver, err := t.WaitReceiveChannel()
				if err != nil {
					log.Print(err)
					break
				}

				nano.NewThingeyAdapter(receiver).Listen()
			}
		}()
	}

}
