package nano

import (
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
	sender       libchan.Sender
	receiver     libchan.Receiver
	remoteSender libchan.Sender
}

func NewThingeyRepository(sender libchan.Sender, receiver libchan.Receiver, remoteSender libchan.Sender) ThingeyRepository {
	return &thingeyRepository{
		sender:       sender,
		receiver:     receiver,
		remoteSender: remoteSender,
	}
}

func (repo *thingeyRepository) Create(thingey *Thingey) error {
	request := &Request{
		Payload: &ThingeyCreateRequest{
			Thingey: thingey,
		},
		ResponseChan: repo.remoteSender,
	}
	if err := repo.sender.Send(request); err != nil {
		return err
	}
	response := &Response{}
	if err := repo.receiver.Receive(response); err != nil {
		return err
	}
	return response.Err
}

func (repo *thingeyRepository) Delete(thingey *Thingey) error {
	request := &Request{
		Payload: &ThingeyDeleteRequest{
			Thingey: thingey,
		},
		ResponseChan: repo.remoteSender,
	}
	if err := repo.sender.Send(request); err != nil {
		return err
	}
	response := &Response{}
	if err := repo.receiver.Receive(response); err != nil {
		return err
	}
	return response.Err
}

func (repo *thingeyRepository) Get(thingeyID string) (*Thingey, error) {
	request := &Request{
		Payload: &ThingeyGetRequest{
			ThingeyID: thingeyID,
		},
		ResponseChan: repo.remoteSender,
	}
	if err := repo.sender.Send(request); err != nil {
		return nil, err
	}
	response := &Response{}
	if err := repo.receiver.Receive(response); err != nil {
		return nil, err
	}
	if response.Err != nil {
		return nil, response.Err
	}
	switch thingeyGetResponse := response.Payload.(type) {
	case *ThingeyGetResponse:
		return thingeyGetResponse.Thingey, nil
	}
	return nil, errors.New("Unknown response type")
}

func (repo *thingeyRepository) List() ([]*Thingey, error) {
	request := &Request{
		Payload:      &ThingeyListRequest{},
		ResponseChan: repo.remoteSender,
	}
	if err := repo.sender.Send(request); err != nil {
		return nil, err
	}
	response := &Response{}
	if err := repo.receiver.Receive(response); err != nil {
		return nil, err
	}
	if response.Err != nil {
		return nil, response.Err
	}
	switch thingeyListResponse := response.Payload.(type) {
	case *ThingeyListResponse:
		return thingeyListResponse.Thingeys, nil
	}
	return nil, errors.New("Unknown response type")
}
