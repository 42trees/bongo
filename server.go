package bongo

import (
	"fmt"
	"log"
	"net/http"
)

func Server(p string) {
	fmt.Println("Listening on port", p)
	log.Fatal(http.ListenAndServe(":"+p, http.FileServer(http.Dir("./_site/"))))
}
