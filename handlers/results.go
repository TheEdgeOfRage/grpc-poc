package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"grpc-test/client"
	"grpc-test/models"
)

func decodeJSON(r io.Reader, ch chan models.Row) {
	dec := json.NewDecoder(r)

	_, err := dec.Token()
	if err != nil {
		log.Fatalf("did not find [ token: %v", err)
	}

	for dec.More() {
		var row models.Row
		err := dec.Decode(&row)
		if err != nil {
			log.Fatalf("failed to decode row: %v", err)
		}
		ch <- row
	}

	_, err = dec.Token()
	if err != nil {
		log.Fatalf("did not find ] token: %v", err)
	}
	close(ch)
}

func getRows(r io.Reader) <-chan models.Row {
	ch := make(chan models.Row, 10)
	go decodeJSON(r, ch)
	return ch
}

func GetResults(c *gin.Context) {
	grpcClient := client.NewClient(c)
	r := grpcClient.GetResults()

	c.Header("content-type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)

	c.Writer.Write([]byte("["))
	enc := json.NewEncoder(c.Writer)
	notFirst := false
	for row := range getRows(r) {
		if notFirst {
			c.Writer.Write([]byte(","))
		}
		notFirst = true
		if err := enc.Encode(row); err != nil {
			log.Fatalf("failed to encode row: %v\n", err)
		}
	}
	c.Writer.Write([]byte("]"))
	c.Writer.Flush()
}
