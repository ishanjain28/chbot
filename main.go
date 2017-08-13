package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// PORT on which HTTP Server is started
var PORT = os.Getenv("PORT")

func main() {

	http.HandleFunc("/", handler)

	log.Fatalln(http.ListenAndServe(":"+PORT, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)

	defer r.Body.Close()

	io.Copy(os.Stdout, r.Body)
}
