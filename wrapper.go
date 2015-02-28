package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karlseguin/typed"
)

// Considered a few different libraries
// https://github.com/karlseguin/typed <-- liked the syntax the most
// https://github.com/bitly/go-simplejson <-- only one that supports creation
// https://github.com/antonholmquist/jason

func WrapperMoveHandler(w http.ResponseWriter, r *http.Request) {
	//bytes, _ := ioutil.ReadAll(r.Body)
	request, err := typed.JsonStream(r.Body)
	if err != nil { 
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", request)
	fmt.Printf("Turn: %d\n", request.Int("turn"))
	
	// actually no way to get nested arrays with typed
	//c := request.Objects("snakes")[0] ...
	//fmt.Printf("Snake position: %v\n", c)

	jsonResponse, err := json.Marshal(map[string]interface{} {
        "move": "down",
        "taunt": request.Objects("snakes")[0].String("taunt"),
    })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}