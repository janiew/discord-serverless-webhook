package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PratikMahajan/Twitter-Knative-Serverless-Video-Download/config"
	ce "github.com/cloudevents/sdk-go"
)



func main() {

	port, err := strconv.Atoi(config.MustGetEnvVar("PORT", "8080"))
	if err != nil {
		log.Fatalf("failed to parse port, %s", err.Error())
	}

	// Handler Mux
	mux := http.NewServeMux()

	// Ingres API Handler
	t, err := ce.NewHTTPTransport(
		ce.WithMethod("POST"),
		ce.WithPath("/"),
		ce.WithPort(port),
	)
	if err != nil {
		log.Fatalf("failed to create CloudEvents transport, %s", err.Error())
	}

	// wire handler for CE
	t.SetReceiver(&eventReceiver{})

	// Health Handler
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, err = fmt.Fprint(w, "ok")
		if err != nil{
			log.Printf("failed to write /health : %s", err.Error())
		}
	})

	// Events or UI Handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method, %s", r.Method)
		if r.Method == "POST" {
			t.ServeHTTP(w, r)
			return
		}
		_, err := fmt.Fprint(w, "Nothing to see here. Use POST to send CloudEvents")
		if err != nil{
			log.Panicf("failed to write / : %s", err.Error())
		}


	})

	a := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(a, mux))

}