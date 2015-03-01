package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Good reading on unmarshalling into structures:
// https://eager.io/blog/go-and-json/

// http://golang.org/pkg/encoding/json/#Unmarshal

// https://github.com/ChimeraCoder/gojson
// Nicer formating (tab spaces, alphabetical, capitalizing ID and URL)
// That capitalization is consistent with the golang style: https://github.com/golang/go/wiki/CodeReviewComments#initialisms

type ChimeraCoderMoveRequest struct {
	Board [][]struct {
		Snake interface{} `json:"snake"`
		State string      `json:"state"`
	} `json:"board"`
	Food   [][]int `json:"food"`
	GameID string  `json:"game_id"`
	Snakes []struct {
		Color   string  `json:"color"`
		Coords  [][]int `json:"coords"`
		HeadURL string  `json:"head_url"`
		Name    string  `json:"name"`
		Taunt   string  `json:"taunt"`
		URL     string  `json:"url"`
	} `json:"snakes"`
	Turn int `json:"turn"`
}

// http://mholt.github.io/json-to-go/
type JSONToGoMoveRequest struct {
	Turn int `json:"turn"`
	GameId string `json:"game_id"`
	Food [][]int `json:"food"`  // Wish these were some sort of XY coords
	Snakes []struct {
		Url string `json:"url"`
		HeadUrl string `json:"head_url"`
		Coords [][]int `json:"coords"`
		Name string `json:"name"`
		Taunt string `json:"taunt"`
		Color string `json:"color"`
	} `json:"snakes"`
	Board [][]struct {
		State string `json:"state"`
		Snake interface{} `json:"snake"` // This actually should be: Snake string `json:"snake"` but strings can't be null in Go :(
	} `json:"board"`
}

// http://godoc.org/github.com/str1ngs/jflect
type JflectMoveRequest struct {
	Board  []interface{} `json:"board"`
	Food   []interface{} `json:"food"`
	GameId string        `json:"game_id"`
	Snakes []struct {
		Color   string        `json:"color"`
		Coords  []interface{} `json:"coords"`
		HeadUrl string        `json:"head_url"`
		Name    string        `json:"name"`
		Taunt   string        `json:"taunt"`
		Url     string        `json:"url"`
	} `json:"snakes"`
	Turn int `json:"turn"`
}

// https://github.com/tmc/json-to-struct
type JSONToStructMoveRequest struct {
	Board [][]struct {
		Snake interface{} `json:"snake"`
		State string      `json:"state"`
	} `json:"board"`
	Food   [][]float64 `json:"food"`		// Silly
	GameID string      `json:"game_id"`
	Snakes []struct {
		Color   string      `json:"color"`
		Coords  [][]float64 `json:"coords"`
		HeadURL string      `json:"head_url"`
		Name    string      `json:"name"`
		Taunt   string      `json:"taunt"`
		URL     string      `json:"url"`
	} `json:"snakes"`
	Turn float64 `json:"turn"`
}

// More details on null strings:
// http://invalidlogic.com/2012/10/16/golang-oddity-1/#comment-689752257
// http://play.golang.org/p/tXksf9e9rU
// tl;dr, if a string can be null, need to decode into a pointer :/

type Coord struct {
	X int
	Y int
}

func (c *Coord) UnmarshalJSON(b []byte) error {
    var tmp []int
    if err := json.Unmarshal(b, &tmp); err != nil {
        return err
    }

    if len(tmp) != 2 {
    	return errors.New("Coord only accepts a length two array")
    }

    c.X, c.Y = tmp[0], tmp[1]

    return nil
}

func (c Coord) String() string {
    return fmt.Sprintf("%d,%d", c.X, c.Y)
}


type CustomMoveRequest struct {
	Board [][]struct {
		Snake *string `json:"snake"`
		State string      `json:"state"`
	} `json:"board"`
	Food   []Coord `json:"food"`
	GameID string  `json:"game_id"`
	Snakes []struct {
		Color   string  `json:"color"`
		Coords  []Coord `json:"coords"`
		HeadURL string  `json:"head_url"`
		Name    string  `json:"name"`
		Taunt   string  `json:"taunt"`
		URL     string  `json:"url"`
	} `json:"snakes"`
	Turn int `json:"turn"`
}


type MoveResponse struct {
	Move string  `json:"move"`
	Taunt string `json:"taunt,omitempty"`
}

func StrictStructureMoveHandler(w http.ResponseWriter, r *http.Request) {
	request := CustomMoveRequest{}

	// use NewDecoder(bytes).Decode() instead of Unmarshall
	// https://www.datadoghq.com/2014/07/crossing-streams-love-letter-gos-io-reader/
	// http://stackoverflow.com/questions/21197239/decoding-json-in-golang-using-json-unmarshal-vs-json-newdecoder-decode

	// Doing this only defines err for the body of the if statement
	// https://groups.google.com/forum/#!msg/golang-nuts/__NG4nQJFMo/35ZPdeakCWAJ
	// http://stackoverflow.com/a/18698091
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil { 
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// when printing structs, the plus flag (%+v) adds field names
	// https://golang.org/pkg/fmt/
	fmt.Printf("%+v\n", request)
	fmt.Printf("%d\n", request.Turn)
	fmt.Printf("Snake position: %s\n", request.Snakes[0].Coords[0])

	response := MoveResponse{Move: "down", Taunt: request.Snakes[0].Taunt}

    // http://www.alexedwards.net/blog/golang-response-snippets#json
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
