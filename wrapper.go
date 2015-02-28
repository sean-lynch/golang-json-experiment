package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/karlseguin/typed"
	"github.com/bitly/go-simplejson"
	"github.com/antonholmquist/jason"
)

// Considered a few different libraries
// https://github.com/karlseguin/typed
// Liked the syntax the most
// No way to get nested arrays with typed, had to write a Json Stream handler

func TypedMoveHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	request, err := typed.Json(bytes)
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

// https://github.com/bitly/go-simplejson
// Only one that supports creation

func SimpleJSONMoveHandler(w http.ResponseWriter, r *http.Request) {
	request, err := simplejson.NewFromReader(r.Body)
	if err != nil { 
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", request)
	turn, _ := request.Get("turn").Int()
	fmt.Printf("Turn: %d\n", turn)
	
	// Can't get nested arrays but not typed ones
	// so we're back to messing with json.Numbers
	c, _ := request.Get("snakes").GetIndex(0).Get("coords").GetIndex(0).Array()
	x, _ := c[0].(json.Number).Int64()
	y, _ := c[1].(json.Number).Int64()
	fmt.Printf("Snake position: %d, %d\n", x, y)


	taunt, _ := request.Get("snakes").GetIndex(0).Get("taunt").String()

	response := simplejson.New()
	response.Set("move", "down")
	response.Set("taunt", taunt)

	jsonResponse, err := response.Encode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// https://github.com/antonholmquist/jason
// Not good for chaining
// Can't get nested arrays either

func JasonMoveHandler(w http.ResponseWriter, r *http.Request) {
	request, err := jason.NewObjectFromReader(r.Body)
	if err != nil { 
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	turn, _ := request.GetInt64("turn")
	fmt.Printf("Turn: %d\n", turn)

	snakes, _ := request.GetObjectArray("snakes")
	snake := snakes[0]

	// Can't get nested arrays, so no coords


	taunt,_ := snake.GetString("taunt")

	jsonResponse, err := json.Marshal(map[string]interface{} {
        "move": "down",
        "taunt": taunt,
    })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}