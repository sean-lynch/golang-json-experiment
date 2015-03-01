package main

import (
    "net/http"
)

func main() {
    http.HandleFunc("/move", MoveHandler)
    http.ListenAndServe(":8000", nil)
}
