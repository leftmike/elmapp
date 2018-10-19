package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/leftmike/elmapp/api"
)

func main() {
	fmt.Println("Elm App")

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
