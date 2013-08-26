package main

import (
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func queueRequest(request *http.Request, channel *amqp.Channel) {

}

func sendErrorResponse(err error, w *http.ResponseWriter) {

}

func main() {

	connection, err := amqp.Dial("amqp://localhost:5672")

	if err != nil {
		panic("Could not connect to amqp server")
	}

	channel, err := connection.Channel()

	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		uri_split := strings.Split("/", request.RequestURI)

		if len(uri_split) < 1 {
			//TODO error on this condition
		}

		exchange := uri_split[0]
		routing_key := ""
		content_type := request.Header.Get("ContentType")

		if len(uri_split) > 1 {
			routing_key = uri_split[1]
		}

		if content_type == "" {
			content_type = "text/json"
		}

		body, err := ioutil.ReadAll(request.Body)

		if err != nil {
			sendErrorResponse(err, &w)
		}

		err = channel.Publish(
			exchange,
			routing_key,
			false,
			false,
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     content_type,
				ContentEncoding: "",
				Body:            body,
				DeliveryMode:    amqp.Persistent,
				Priority:        0,
			},
		)

		if err != nil {
			sendErrorResponse(err, &w)
		}

		io.WriteString(w, "Wrote hello world to queue\n")
	})

	http.ListenAndServe(":12345", nil)
}
