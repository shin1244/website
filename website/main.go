package main

import (
	"mongo/website/app"
	"net/http"
)

func main() {
	mux := app.NewHandler()
	http.ListenAndServe(":3000", mux)
}
