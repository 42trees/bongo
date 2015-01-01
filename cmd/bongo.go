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
import "log"

import "github.com/russross/blackfriday"
import "gopkg.in/yaml.v2"

func main() {


  /*
  usage: bongo (no params)
  bongo builds the site in _site/

  Make it work. Make it right. Make it fast.

  flags: -content
  */
  var contentPath = flag.String("content", "content", "Path to content")

  flag.Parse()

  //  flag.PrintDefaults()


  fmt.Println("content:", *contentPath)

  dirs,_ := ioutil.ReadDir(*contentPath)
  for _,n := range dirs {
    if n.IsDir() { //Run function that builds the posts or pages
      fmt.Println(n.Name())
      build(*contentPath+"/"+n.Name())
    }
  }

  index()

  frontmatter("../content/pages/cv.html")


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

func parseHTML(filename string) string {
  start := "{{define \"content\"}}\n"
  end := "{{end}}\n"

  abs,_ := filepath.Abs(filename)
  f,_ := ioutil.ReadFile(abs)
  s := string(f[:])

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

func build(p string) {
  //@TODO check errors
  md,_ := filepath.Glob(p+"/*.md")
  parseFiles(md)

  html,_ := filepath.Glob(p+"/*.html")
  parseFiles(html)

}

func parseFiles(f []string) {
  for i, n := range f {

    _, filename := filepath.Split(n)
    title := strings.TrimSuffix(filename, filepath.Ext(filename))
    fmt.Println(title)

    var td TemplateData;

    td.Slug = title;

    //date, err := dateStr(title)
    //fmt.Println(date)
    //fmt.Println(err)

    fmt.Printf("%v:%v\n", i, n)

    ext := filepath.Ext(n)
    s := parseMD(n)
    if ext == ".md" {
      //@TODO scan the file. look for frontmatter opening and closing ---
      // if found, unmarshal the contents
      /* err := yaml.Unmarshal([]byte(s), &td)
      if err != nil {
        log.Printf("error: %v", err)
      }
      fmt.Printf("--- td:%+v\n\n", td)*/

    }

    t, _ := template.ParseFiles("../templates/layout.html")
    t.Parse(s)


    d := "../_site/"+title
    makeDir(d)
    fmt.Println(d) 
    index := d+"/index.html"
    fmt.Println(index)
    file,err := os.Create(index)

    if err != nil {
      fmt.Printf("ERROR: %v\n", err)
      panic(err)
    }

    w := bufio.NewWriter(file)
    t.Execute(w, td)
    w.Flush()
  }
}

type TemplateData struct {
  Slug string
  Title string
  Layout string
  Permalink string
}

func index() {

  t, _ := template.ParseFiles("../templates/layout.html", "../content/index.html")
  d := "../_site"
  index := d+"/index.html"
  fmt.Println(index)
  file,err := os.Create(index)
  if err != nil {
    panic(err)
  }
  w := bufio.NewWriter(file)

  var td TemplateData;
  td.Title = "Home"
  t.Execute(w, td)
  w.Flush()


}


// Read the file line by line and find the frontmatter
// Return a TemplateData struct if successful
func frontmatter(f string) (TemplateData, error) {

  file, err := os.Open(f)
  if err != nil {
    fmt.Printf("Unable to open file: %v, %v", file, err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  const fm = "---"

  start, end := false, false

  b := []string{}

  for scanner.Scan() {

    if start && end { //already done with FM. No need to continue
      fmt.Printf("b: %+v\n", b)
      break
    }

    t := scanner.Text()

    if t == fm {

      if start && end { //error state. already had frontmatter. Return an error on this file
        fmt.Println("Invalid Frontmatter")
      }


      if start { //already found open start --- tag
        end = true
      }

      start = true
      fmt.Printf("start: %v, end: %v\n", start, end)
      fmt.Println(scanner.Text())

    } else {

      if start {
        b = append(b, t)
      }
    }
  }
  fmt.Println(b)

  var td TemplateData
  s := strings.Join(b, "\n")
  e := yaml.Unmarshal([]byte(s), &td)
  if e != nil {
    log.Printf("error: %v", err)
  }
  fmt.Printf("--- td:%+v\n\n", td)

  return td, e 
}

