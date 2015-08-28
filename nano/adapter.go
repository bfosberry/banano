package nano

import (
	"errors"
	"log"

	"github.com/docker/libchan"
)

type Thingey struct {
	ID   string
	Data string
}

type Request struct {
	ResponseChan libchan.Sender
	Payload      interface{}
}

type ThingeyCreateRequest struct {
	Thingey *Thingey
}

type ThingeyDeleteRequest struct {
	Thingey *Thingey
}

type ThingeyGetRequest struct {
	ThingeyID string
}

type ThingeyListRequest struct {
}

type Response struct {
	Payload interface{}
	Err     error
}

type ThingeyCreateResponse struct {
}

type ThingeyDeleteResponse struct {
}

type ThingeyGetResponse struct {
	Thingey *Thingey
}

type ThingeyListResponse struct {
	Thingeys []*Thingey
}

type ThingeyAdapter interface {
	Listen()
}

type thingeyAdapter struct {
	receiver libchan.Receiver
	thingeys map[string]*Thingey
}

func NewThingeyAdapter(receiver libchan.Receiver) ThingeyAdapter {
	return &thingeyAdapter{
		receiver: receiver,
		thingeys: map[string]*Thingey{},
	}
}

func (adapter *thingeyAdapter) Listen() {
	request := &Request{}
	response := &Response{}
	err := adapter.receiver.Receive(request)
	if err != nil {
		log.Print(err)
		return
	}

	payload := request.Payload
	switch payload := payload.(type) {
	case *ThingeyCreateRequest:
		log.Println("received create request")
		adapter.thingeys[payload.Thingey.ID] = payload.Thingey
		response.Payload = &ThingeyCreateResponse{}
	case *ThingeyDeleteRequest:
		log.Println("received delete request")
		delete(adapter.thingeys, payload.Thingey.ID)
		response.Payload = &ThingeyDeleteResponse{}
	case *ThingeyGetRequest:
		log.Println("received get request")
		response.Payload = &ThingeyGetResponse{
			Thingey: adapter.thingeys[payload.ThingeyID],
		}
	case *ThingeyListRequest:
		log.Println("received list request")
		thingeys := []*Thingey{}
		for _, t := range adapter.thingeys {
			thingeys = append(thingeys, t)
		}
		response.Payload = &ThingeyListResponse{
			Thingeys: thingeys,
		}
	default:
		log.Println("unknown request type")
		response.Err = errors.New("unknown request type")
	}
	request.ResponseChan.Send(response)
}
