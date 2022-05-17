package main

import (
	"encoding/json"
	"fmt"
	siwsMessage "github.com/Web3Auth/siws-go/pkg/message"
	"github.com/Web3Auth/siws-go/runner"
	"io/ioutil"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")


	switch r.Method {
	case "POST":
		switch r.URL.Path {
			case "/verify" :
				Verify(w,r)
			case "/create":
				Create(w,r)
			case "/prepareMessage":
				PrepareMessage(w,r)
			default:
				http.Error(w, "404 not found.", http.StatusNotFound)
				return
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", hello)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func Create(w http.ResponseWriter,r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	var payload Payload
	err = json.Unmarshal(body,&payload)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	SIWSMessage, err := runner.InitMessage(payload.Domain,payload.Address,payload.URI,payload.Version,payload.Options)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	siwsAsJSON, err := json.Marshal(&SIWSMessage)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	fmt.Fprintf(w, string(siwsAsJSON))
}

func PrepareMessage(w http.ResponseWriter,r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	var siwsMessage siwsMessage.Message
	err = json.Unmarshal(body,&siwsMessage)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	fmt.Fprintf(w, siwsMessage.PrepareMessage())

}

func Verify(w http.ResponseWriter,r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	var siwsMessage siwsMessage.Message
	err = json.Unmarshal(body,&siwsMessage)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	res,err := siwsMessage.Verify(siwsMessage.Signature.S,nil,nil)
	if err != nil {
		fmt.Fprintf(w, "Err: %v", err)
	}
	fmt.Fprintf(w,"%v",res)
}