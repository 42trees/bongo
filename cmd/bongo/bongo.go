package main

import "fmt"
import "path/filepath"

func main() {


  /*
  usage: bongo (no params)
  bongo builds the site in _site/

  design decision. Simpler. Opionionated. Migrate the site to bongo, worry less write more.



  flags: -s [or something]

  rebuild the site: serve it statically on port 4242

  NO admin interface or GUI. just simple text.

  build:

  pages/ 
  my-pagename.md - layout defined in config or YAML in header
  => builds as http://mysite.local:4242/my-pagename/index.html


  accessible in pages['category'][] somehow.

  Need to be able to create index or posts/pages. Pagination etc
  */


  fmt.Println("usage: bongo")
  content()

}

func content() {
  f,_ := filepath.Glob("/home/karlcordes/go/src/github.com/42trees/bongo/content/*.md")


  for i, n := range f {
    fmt.Printf("%v:%v\n", i, n)
  }

}
