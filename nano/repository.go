package nano

import (
	"encoding/json"
	"errors"

	"github.com/docker/libchan"
)

type ThingeyRepository interface {
	Create(*Thingey) error
	Delete(*Thingey) error
	Get(string) (*Thingey, error)
	List() ([]*Thingey, error)
}

type thingeyRepository struct {
	senderFunc   SenderFunc
	receiver     libchan.Receiver
	remoteSender libchan.Sender
}

func NewThingeyRepository(senderFunc SenderFunc, receiver libchan.Receiver, remoteSender libchan.Sender) ThingeyRepository {
	return &thingeyRepository{
		senderFunc:   senderFunc,
		receiver:     receiver,
		remoteSender: remoteSender,
	}
}

func (repo *thingeyRepository) Create(thingey *Thingey) error {
	req := &ThingeyCreateRequest{
		Thingey: thingey,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request := &Request{
		Payload:      data,
		Type:         "ThingeyCreateRequest",
		ResponseChan: repo.remoteSender,
	}
	response, err := repo.dispatch(request)
	if err != nil {
		return err
	}
	return response.Err
}

func (repo *thingeyRepository) Delete(thingey *Thingey) error {
	req := &ThingeyDeleteRequest{
		Thingey: thingey,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request := &Request{
		Payload:      data,
		Type:         "ThingeyDeleteRequest",
		ResponseChan: repo.remoteSender,
	}
	response, err := repo.dispatch(request)
	if err != nil {
		return err
	}
	return response.Err
}

func (repo *thingeyRepository) Get(thingeyID string) (*Thingey, error) {
	req := &ThingeyGetRequest{
		ThingeyID: thingeyID,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &Request{
		Payload:      data,
		Type:         "ThingeyGetRequest",
		ResponseChan: repo.remoteSender,
	}
	response, err := repo.dispatch(request)
	if err != nil {
		return nil, err
	}

	if response.Type == "ThingeyGetResponse" {
		resp := &ThingeyGetResponse{}
		err := json.Unmarshal(response.Payload, resp)
		if err != nil {
			return nil, err
		}
		return resp.Thingey, nil
	}
	return nil, errors.New("Unknown response type")
}

func (repo *thingeyRepository) List() ([]*Thingey, error) {
	req := &ThingeyListRequest{}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &Request{
		Payload:      data,
		Type:         "ThingeyListRequest",
		ResponseChan: repo.remoteSender,
	}
	response, err := repo.dispatch(request)
	if err != nil {
		return nil, err
	}
	if response.Type == "ThingeyListResponse" {
		resp := &ThingeyListResponse{}
		err := json.Unmarshal(response.Payload, resp)
		if err != nil {
			return nil, err
		}
		return resp.Thingeys, nil
	}
	return nil, errors.New("Unknown response type")
}

func (repo *thingeyRepository) dispatch(req *Request) (*Response, error) {
	sender, err := repo.senderFunc()
	if err != nil {
		return nil, err
	}

	if err := sender.Send(req); err != nil {
		return nil, err
	}
	response := &Response{}
	if err := repo.receiver.Receive(response); err != nil {
		return nil, err
	}
	return response, nil
}
