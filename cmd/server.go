package main

import "net/http"
import "log"

func main() {
	log.Fatal(http.ListenAndServe(":8081", http.FileServer(http.Dir("../_site"))))
}
