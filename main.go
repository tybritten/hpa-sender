package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type eventdata struct {
	//
	FirstTimestamp string `json:"firstTimestamp"`
	InvolvedObject struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Name       string `json:"name"`
		Namespace  string `json:"namespace"`
	} `json:"involvedObject"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Source  struct {
		Component string `json:"component"`
	} `json:"source"`
}

func main() {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("will listen on :8080\n")
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, receive))
}

func receive(ctx context.Context, event cloudevents.Event) {
	data := &eventdata{}
	if err := event.DataAs(data); err != nil {
		fmt.Printf("Got Data Error: %s\n", err.Error())
	}
	if data.Source.Component == "horizontal-pod-autoscaler" {
		fmt.Printf("HPA Event Received \n")
		fmt.Printf("----------------------------\n")
		secret, _ := checkhpa(data.InvolvedObject.Name, data.InvolvedObject.Namespace)
		if len(secret) > 0 {
			fmt.Printf("Sending Message \n")
			sendmessage(data.FirstTimestamp, data.Message, secret)
		}

	}
}

func sendmessage(timestamp string, message string, secretloc string) {
	secret, err := getsecret(secretloc)
	data := prepMessage(message, timestamp, secret)
	b := bytes.NewBuffer(data.Body)
	req, err := http.NewRequest("POST", data.URL, b)
	if err != nil {
		print(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range data.Headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		print(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %+v\n", string(body))
	if err != nil {
		print(err.Error())
	}

}

func prepMessage(message string, timestamp string, secret map[string][]byte) SecretData {
	var secretstruct SecretData
	err := json.Unmarshal(secret["headers"], &secretstruct.Headers)
	if err != nil {
		print(err.Error())
	}
	b1 := bytes.ReplaceAll(secret["body"], []byte("_message_"), []byte(message))
	secretstruct.Body = bytes.ReplaceAll(b1, []byte("_time_"), []byte(timestamp))
	secretstruct.URL = string(secret["url"])
	return secretstruct
}
