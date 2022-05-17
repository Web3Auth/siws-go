package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/Web3Auth/siws-go/pkg/types"
)
func hello(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		switch r.URL.Path {
			case "/verify" :
				Verify(w,r)
			case "/create":
				Create(w,r)
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
	log.Println(string(body))
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", string(body))

}

func Verify(w http.ResponseWriter,r *http.Request) {

}