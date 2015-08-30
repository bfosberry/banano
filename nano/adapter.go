package nano

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/docker/libchan"
)

type SenderFunc func() (libchan.Sender, error)

type Thingey struct {
	ID   string
	Data string
}

type Request struct {
	ResponseChan libchan.Sender
	Payload      []byte
	Type         string
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
	Payload []byte
	Type    string
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
	Listen(libchan.Receiver) error
}

type thingeyAdapter struct {
	thingeys map[string]*Thingey
}

func NewThingeyAdapter() ThingeyAdapter {
	return &thingeyAdapter{
		thingeys: map[string]*Thingey{},
	}
}

func (adapter *thingeyAdapter) Listen(receiver libchan.Receiver) error {
	request := &Request{}
	response := &Response{}
	err := receiver.Receive(request)
	if err != nil {
		log.Print(err)
		return err
	}

	switch request.Type {
	case "ThingeyCreateRequest":
		payload := &ThingeyCreateRequest{}
		if err = json.Unmarshal(request.Payload, payload); err != nil {
			return err
		}
		adapter.thingeys[payload.Thingey.ID] = payload.Thingey
		response.Payload, err = json.Marshal(&ThingeyCreateResponse{})
		response.Type = "ThingeyCreateResponse"
		if err != nil {
			return err
		}
	case "ThingeyDeleteRequest":
		payload := &ThingeyDeleteRequest{}
		if err = json.Unmarshal(request.Payload, payload); err != nil {
			return err
		}
		delete(adapter.thingeys, payload.Thingey.ID)
		response.Payload, err = json.Marshal(&ThingeyDeleteResponse{})
		response.Type = "ThingeyDeleteResponse"
		if err != nil {
			return err
		}
	case "ThingeyGetRequest":
		payload := &ThingeyGetRequest{}
		if err = json.Unmarshal(request.Payload, payload); err != nil {
			return err
		}
		resp := &ThingeyGetResponse{
			Thingey: adapter.thingeys[payload.ThingeyID],
		}
		response.Payload, err = json.Marshal(resp)
		response.Type = "ThingeyGetResponse"
		if err != nil {
			return err
		}
	case "ThingeyListRequest":
		payload := &ThingeyListRequest{}
		if err = json.Unmarshal(request.Payload, payload); err != nil {
			return err
		}
		thingeys := []*Thingey{}
		for _, t := range adapter.thingeys {
			thingeys = append(thingeys, t)
		}
		resp := &ThingeyListResponse{
			Thingeys: thingeys,
		}
		response.Payload, err = json.Marshal(resp)
		response.Type = "ThingeyListResponse"
		if err != nil {
			return err
		}
	default:
		response.Err = errors.New("unknown request type")
	}
	return request.ResponseChan.Send(response)
}
