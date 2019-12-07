package processor

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

//
// Process contains information about a single proces.
//
type Process struct {
	Filesize  int    `json:"filesize"`
	Checksum  uint32 `json:"checksum"`
	Extension string `json:"extension"`
}

func handleProcess() error {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		"images",
		true,
		false,
		false,
		false,
		nil,
	)

	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	messages, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	for message := range messages {
		fmt.Println("==> Recivied a message: " + string(message.Body))

		process := &Process{}
		err = json.NewDecoder(bytes.NewReader(message.Body)).Decode(process)
		if err != nil {
			fmt.Println(`==> JSON decode error: ` + err.Error())
		}

		err = handleMessage(process)
		if err != nil {
			fmt.Println(`==> Handle message error: ` + err.Error())
		}

		message.Ack(false)
	}

	return nil

}
