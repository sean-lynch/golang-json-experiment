package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Good reading here too:
// http://attilaolah.eu/2013/11/29/json-decoding-in-go/


func LooseStructureMoveHandler(w http.ResponseWriter, r *http.Request) {
    // Can't use := because dunno
    var request map[string]interface{}
    d := json.NewDecoder(r.Body)
    // http://golang.org/pkg/encoding/json/#Decoder.UseNumber
    // http://attilaolah.eu/2013/11/29/json-decoding-in-go/
    d.UseNumber()
	if err := d.Decode(&request); err != nil { 
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

    fmt.Printf("%+v\n", request)

    // Decode unmarshals to float64 by default for some reason 
    // So if you didn't use .UseNumber(), you'd need to do something insane like this:
    // https://code.google.com/p/go/issues/detail?id=5562
    // https://github.com/golang/go/issues/5562
    // turn := int(request["turn"].(float64))
    // And Int64() is just using strconv.ParseInt() underneath
    // http://golang.org/src/encoding/json/decode.go?s=5163:5201#L155
    // Oh shit, this is just as ugly.
    turn, _ := request["turn"].(json.Number).Int64()
	fmt.Printf("Turn: %d\n", turn)

    // http://www.alexedwards.net/blog/golang-response-snippets#json
	jsonResponse, err := json.Marshal(map[string]interface{} {
        "move": "down",
        "taunt": request["snakes"].([]interface{})[0].(map[string]interface{})["taunt"],  // Thanks: https://eager.io/blog/go-and-json/
    })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
