package main

import "fmt"
import "path/filepath"
import "flag"
import "html/template"
import "bufio"
import "io/ioutil"
import "strings"
import "os"
import "time"

import "github.com/russross/blackfriday"

func main() {


  /*
  usage: bongo (no params)
  bongo builds the site in _site/

  design decision. Simpler. Migrate the site to bongo, worry less write more.


  Make it work. Make it right. Make it fast.

  flags: -content
  */
  var contentPath = flag.String("content", "content", "Path to content")

  flag.Parse()

//  flag.PrintDefaults()


  fmt.Println("content:", *contentPath)

  f,_ := filepath.Glob(*contentPath+"/*.md")


  for i, n := range f {

    _, filename := filepath.Split(n)
    title := strings.TrimSuffix(filename, filepath.Ext(filename))
    fmt.Println(title)
    
    //date, err := dateStr(title)
    //fmt.Println(date)
    //fmt.Println(err)

    fmt.Printf("%v:%v\n", i, n)

    s := parseMD(n)

    t, _ := template.ParseFiles("../templates/layout.html")
    t.Parse(s)


    d := "../_site/"+title
    makeDir(d)
    fmt.Println(d) 
    index := d+"/index.html"
    fmt.Println(index)
    file,err := os.Create(index)
    
    if err != nil {
      panic(err)
    }
    w := bufio.NewWriter(file)
    t.Execute(w, nil)
    w.Flush()
  }

}


func parseMD(filename string) string {

  start := "{{define \"content\"}}\n"
  end := "{{end}}\n"

  abs,_ := filepath.Abs(filename)

  f,_ := ioutil.ReadFile(abs)


  md := blackfriday.MarkdownCommon(f)
  s := string(md[:])

  strs := []string{start, s, end}

  mdstr := strings.Join(strs, "")

  return mdstr
}


func makeDir(p string) error {
  return os.MkdirAll(p, 0755)
}

//Might skip dates in filenames. What's the point?
func dateStr(f string) (time.Time, error) {
  s := strings.Split(f, "-")
  fmt.Println(len(s))
  d := strings.Join(s[:3], "-") //seems a bit silly. But it works for now
  return time.Parse("2006-01-02", d)
}



