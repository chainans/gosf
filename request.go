package gosf

import (
	"log"

	io "github.com/ambelovsky/gosf-socketio"
)

// Message - Standard message type for Socket communications
type Message struct {
	ID      int  `json:"id,omitempty"`
	Success bool `json:"success"`

	Text string      `json:"text,omitempty"`
	Body interface{} `json:"body,omitempty"`
}

// Request represents a single request over an active connection
type Request struct {
	Channel  *io.Channel
	Endpoint string
	Message  *Message
}

// Broadcast sends a message to connected clients joined to the same room
func Broadcast(room string, endpoint string, message *Message) {
	if room != "" {
		ioServer.BroadcastTo(room, endpoint, message)
	} else {
		ioServer.BroadcastToAll(endpoint, message)
	}
}

// Listen creates a listener on an endpoint
func Listen(endpoint string, callback func(request *Request) *Message) {
	ioServer.On(endpoint, func(channel *io.Channel, clientMessage *Message) *Message {
		client := new(Client)
		client.Channel = channel

		request := new(Request)
		request.Endpoint = endpoint
		request.Message = clientMessage

		emit("before-request", client, request)

		response := callback(request)
		if response == nil {
			response = new(Message)
		}

		emit("after-request", client, request, response)

		// Deferred until after return fires so the ack (acknowledgement) has a chance to go back to the client
		defer emit("after-response", client, request, response)

		return request.respond(response)
	})
}

// Respond sends a message back to the client
func (request Request) respond(response *Message) *Message {
	client := new(Client)
	client.Channel = request.Channel

	if &request.Message.ID != nil {
		response.ID = request.Message.ID
	}

	emit("before-response", client, &request, response)

	log.Println("t1")
	client.Channel.Emit(request.Endpoint, response)
	log.Println("t2")
	return response
}
