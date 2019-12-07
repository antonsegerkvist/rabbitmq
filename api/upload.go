package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/streadway/amqp"
)

func handleUpload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(`Content-Type`, `application/json`)

	filenames := r.URL.Query()["filename"]
	if len(filenames) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"status":400}`)
		return
	}
	filename := filenames[0]

	buffer := bytes.Buffer{}
	_, err := io.Copy(&buffer, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"IO_COPY"}`)
		return
	}

	adler32checksum := adler32.Checksum(buffer.Bytes())
	filesize := buffer.Len()
	extension := path.Ext(filename)

	file, err := os.OpenFile(
		"storage/original/"+strconv.FormatUint(uint64(filesize), 10)+"_"+strconv.FormatUint(uint64(adler32checksum), 10)+extension,
		os.O_WRONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"OS_FILE_OPEN"}`)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, &buffer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"IO_COPY"}`)
		return
	}

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"AMQP_DIAL"}`)
		return
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"CONNECTION_CHANNEL"}`)
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"images",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"CHANNEL_QUEUE_DECLARE"}`)
		return
	}

	message := map[string]interface{}{
		"checksum":  adler32checksum,
		"filesize":  filesize,
		"extension": extension,
	}

	messageString, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"JSON_MARSHAL"}`)
		return
	}

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(messageString),
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"CHANNEL_PUBLISH"}`)
		return
	}

	response := map[string]interface{}{
		"status":    200,
		"checksum":  adler32checksum,
		"filesize":  filesize,
		"extension": extension,
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":500,"reason":"JSON_MARSHAL"}`)
		return
	}

	io.WriteString(w, string(jsonBody))

}
