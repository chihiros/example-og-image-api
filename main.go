package main

import (
	"net/http"
	"server/infra"
)

func main() {
	r := infra.NewRouter()
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
