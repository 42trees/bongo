package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

const bongoVersion = "0.0.1"

func main() {

	/*
		  default usage: bongo - builds the site in _site/

			flags:
			-content
			-new
			-help
			-server
			-version
	*/
	var projectDir = flag.String("new", "", "Create a new bongo project in the specified directory")
	var contentPath = flag.String("content", "content", "Path to content")
	var helpFlag = flag.Bool("help", false, "Show usage")
	var versionFlag = flag.Bool("version", false, "Show version")
	var serverFlag = flag.Bool("server", false, "Build the site and start a webserver")
	var port = flag.String("port", "4242", "Port the webserver will listen on")

	flag.Parse()

	//	flag.PrintDefaults()

	if *projectDir != "" {
		newProject(*projectDir)
		return
	}

	if *helpFlag || *versionFlag {
		help()
		return
	}

	if *serverFlag {
		server(*port)
		return
	}

	startTime := time.Now()

	fmt.Println("content:", *contentPath)
	fmt.Println("content:", *contentPath)

	dirs, _ := ioutil.ReadDir(*contentPath)
	for _, n := range dirs {
		if n.IsDir() { //Run function that builds the posts or pages
			fmt.Println(n.Name())
			build(*contentPath + "/" + n.Name())
		}
	}

	index()

	fmt.Printf("Built in %v ms\n", int(1000*time.Since(startTime).Seconds()))

}

func help() {
	fmt.Printf("bongo %v\n", bongoVersion)
	flag.PrintDefaults()
	return
}

func server(p string) {
	fmt.Println("Listening on port", p)
	log.Fatal(http.ListenAndServe(":"+p, http.FileServer(http.Dir("./_site/"))))
}

func newProject(d string) error {

	var c = filepath.Clean(d)
	fmt.Printf("Creating %v/templates\n", c)
	var e = makeDir(c + "/templates")
	e = makeDir(c + "/content")
	e = makeDir(c + "/_site")
	e = makeDir(c + "/scss")
	e = makeDir(c + "/_site/js")
	e = makeDir(c + "/_site/css")
	return e

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
	md, _ := filepath.Glob(p + "/*.md")
	parseFiles(md)

	html, _ := filepath.Glob(p + "/*.html")
	parseFiles(html)

}

func parseFiles(f []string) {
	for i, n := range f {

		_, filename := filepath.Split(n)
		title := strings.TrimSuffix(filename, filepath.Ext(filename))
		fmt.Println(title)

		var td TemplateData

		//date, err := dateStr(title)
		//fmt.Println(date)
		//fmt.Println(err)

		fmt.Printf("%v:%v\n", i, n)

		td, _ = frontmatter(n)

		if td.Slug == "" {
			fmt.Println("autoslug is:", title)
			td.Slug = title
		}

		t, _ := template.ParseFiles("templates/layout.html")

		d := "_site/" + td.Slug
		makeDir(d)
		fmt.Println(d)
		index := d + "/index.html"
		fmt.Println(index)
		file, err := os.Create(index)

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
	Slug    string
	Title   string
	Layout  string
	Content template.HTML
}

func index() {

	t, _ := template.ParseFiles("templates/layout.html")
	d := "_site"
	index := d + "/index.html"
	fmt.Println(index)
	file, err := os.Create(index)

	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(file)

	var td TemplateData

	td, _ = frontmatter("content/index.html")

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
	sb := []string{} //stripped buffer. ie. No frontmatter

	for scanner.Scan() {

		t := scanner.Text()

		if start && end { //already done with FM. Add to stripped buffer
			sb = append(sb, t)
		}

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

		} else { //not dashes

			if start && !end { //are inside the FM
				b = append(b, t)
				fmt.Println(t)
			}

			if !start && !end { //No frontmatter. Just add it to the 'stripped' buffer
				sb = append(sb, t)
			}

		}
	}

	var td TemplateData

	//Parse the actual YAML inside the front matter
	s := strings.Join(b, "\n")
	e := yaml.Unmarshal([]byte(s), &td)
	if e != nil {
		log.Printf("YAML error: %v", err)
	}

	ext := filepath.Ext(f)
	fmt.Println(ext)

	//This works, but is a bit messy
	if ext == ".md" {
		s := strings.Join(sb, "\n")
		md := blackfriday.MarkdownCommon([]byte(s))
		s = string(md[:])
		td.Content = template.HTML(s)
	}

	if ext == ".html" {
		td.Content = template.HTML(strings.Join(sb, "\n")) //join the array together, convert to HTML
	}

	return td, e
}
