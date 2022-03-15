package workerPool

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func StartServer(HTTPAddr string, pool *Pool) {
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", func(writer http.ResponseWriter, request *http.Request) {
		err := parseReq(writer, request, pool)
		if err != nil {
			fmt.Println("Unable to start the workers: ", err)
		}
	})

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", HTTPAddr)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err.Error())
	}
}

func parseReq(w http.ResponseWriter, r *http.Request, pool *Pool) error {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errors.New("only POST method allowed")
	}

	var workReq = WorkRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&workReq); err != nil {
		http.Error(w, "Invalid JSON payload "+err.Error(), http.StatusBadRequest)
		return err
	}
	defer r.Body.Close()

	// Get the delay.
	delay := time.Second * workReq.Delay

	// Check to make sure the delay is anywhere from 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return errors.New("invalid Delay")
	}

	// Now, we retrieve the person's name from the request.
	name := workReq.Name

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return errors.New("invalid name, it should not be an empty string")
	}

	expiryTime := time.Now().Add(time.Second * pool.Duration)

	// Now, we take the delay, the person's name, and expiry time to make a WorkRequest out of them.
	work := WorkRequest{Name: name, Delay: delay, ExpiryTime: expiryTime}

	fmt.Println("Work request queued")
	pool.collector(work)
	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return nil
}
