package main

import (
    "net/http"
)

func main() {
    http.HandleFunc("/move", JasonMoveHandler)
    http.ListenAndServe(":8000", nil)
}
