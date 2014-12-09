package main

import(
  "net/http"
  "html/template"
  "fmt"
)


func handler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("templates/layout.html", "templates/moo.html")
  t.Execute(w, nil)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("templates/layout.html", "templates/new.html")
  t.Execute(w, nil)

}

func main() {

  http.HandleFunc("/new", newHandler)
  http.HandleFunc("/", handler)
  fmt.Println("Listening on port 4242")
  http.ListenAndServe(":4242", nil)

  /* 
  usage:
   bongod (no params)
    Starts up an admin web interface:
    
    post/page creation and modification.

    images can be uploaded and processed  
    configs edited [sqlite? json?]

    */
}
