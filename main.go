package main

import (
    "net/http"
)

func main() {
    http.HandleFunc("/move", WrapperMoveHandler)
    http.ListenAndServe(":8000", nil)
}
