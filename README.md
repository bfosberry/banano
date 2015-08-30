## Banano

Banano is a simple "nano-service" prototype, implementing transparent libchan connections between layers of an application. The purpose here is to provide a simplified example of how applications can be composed of local and remote components, utilizing a common interface. 

In this way, applications can be broken into small, isolated chunks, and whether those components are deployed independently, or within the same binary, becomes an infrastructure level decision. This can be influenced by design, or by monitoring, resulting in a potentially highly and granularly scalable application. 

## Design

Libchan allows a Go style channel to connect a listener and a sender across a number of different transports including a standard Go channel, and HTTP. By wrapping connections between applications components in generic libchan handlers, we can easily support remote connections as well as local ones.

In this example we have a `repository` connecting to an `adapter`. These communicate using a set of specific request and response objects (e.g. `ThingeyCreateRequest`) marhsalled inside standard Request and Response objects. `sender` and `receiver` libchan handlers are used to dispatch and receive messages. The client side request response cycle looks like this:

```
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
	ResponseChan: remoteSender,
}

if err := sender.Send(req); err != nil {
	return nil, err
}
response := &Response{}
if err := repo.receiver.Receive(response); err != nil {
	return nil, err
}
return response, nil
```

The server side looks like this:

```
request := &Request{}
response := &Response{}
err := receiver.Receive(request)
if err != nil {
	log.Print(err)
	return err
}

if request.Type == "ThingeyCreateRequest" {
	payload := &ThingeyCreateRequest{}
	if err = json.Unmarshal(request.Payload, payload); err != nil {
		return err
	}
	//handle create request
	response.Payload, err = json.Marshal(&ThingeyCreateResponse{})
	response.Type = "ThingeyCreateResponse"
	if err != nil {
		return err
	}
}
```

Much of the work done here is around marhsalling and unmarhsalling data, rather than using libchan. As long as the two components send and receive objects using the libchan Sender and Receiver objects, the code will work with a co-located deployment, or a remote one.

## Local example

To run the local example just run `go run cmd/main.go`. This will start up a local adapter, connect the repository to this adapter, and execute a set of actions against this repository.

## Remote example 

To run the remote example, first start up the server with `go run server/main.go`. This will start an adapter listening on port 8080. Next start the client with `go run client/main.go`, which will create a repository connecting to the remote adapter on port 8080.

## Conclusion

There is one primary different between the local and remote client side code, and that is how the `sender` and `receiver` variables are generated and passed to the repository and adapter. The send/receive code in both examples works exactly the same way. Much of the marhsalling complexity can be abstracted using something like PB, or a better json implementation.

Something to note is that one area where this does not act in a completely transparent manner relates to marhsalling across HTTP. When using a local Go channel transport, you can pass Request objects with interface{} type fields, and on the receiving side type-check those objects to handle them. When dispatching messages across a HTTP transport, those objects are extracted as maps, making type checking impossible. For this reason marshalling is used. 
