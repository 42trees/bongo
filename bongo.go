// Copyright 2015 Karl Cordes
// See LICENSE files for details

package bongo

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const bongoVersion = "0.0.2"

func Help() {
	fmt.Printf("bongo %v\n", bongoVersion)
	flag.PrintDefaults()
	return
}

//@TODO handle dates in filenames
func dateStr(f string) (time.Time, error) {
	s := strings.Split(f, "-")
	fmt.Println(len(s))
	d := strings.Join(s[:3], "-") //seems a bit silly. But it works for now
	return time.Parse("2006-01-02", d)
}

// Build loops through the content types and proceses them
func Build(c *string) {

	startTime := time.Now()
	dirs, err := ioutil.ReadDir(*c)

	if err != nil {
		log.Println("Content directory does not exist")
		log.Printf("Error: %v\n", err)
		return
	}
	var posts map[string]Page
	for _, n := range dirs {
		if n.IsDir() { //Run function that builds the posts or pages
			fmt.Println(n.Name())
			p := *c + "/" + n.Name()

			//@TODO check errors
			md, _ := filepath.Glob(p + "/*.md")
			posts = parseFiles(md)

			html, _ := filepath.Glob(p + "/*.html")
			parseFiles(html)

		}
	}

	Index(posts)
	fmt.Printf("Built in %v ms\n", int(1000*time.Since(startTime).Seconds()))
}

// parseFiles takes an array of filenames
// makes the directory for the page and creates index.html
func parseFiles(f []string) map[string]Page {

	var pages map[string]Page

	pages = make(map[string]Page, len(f))

	for i, n := range f {

		_, filename := filepath.Split(n)
		title := strings.TrimSuffix(filename, filepath.Ext(filename))
		fmt.Println(title)

		var p Page

		//date, err := dateStr(title)
		//fmt.Println(date)
		//fmt.Println(err)

		fmt.Printf("%v:%v\n", i, n)

		p, _ = frontmatter(n)

		if p.Slug == "" {
			fmt.Println("autoslug is:", title)
			p.Slug = title
		}
		pages[p.Slug] = p

		//@TODO don't ignore these
		t, _ := template.ParseFiles("templates/layout.html")

		d := "_site/" + p.Slug
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
		t.Execute(w, p)
		w.Flush()
	}
	return pages
}

type Page struct {
	Slug    string // eg. my-page
	Title   string // My page!
	Layout  string // @TODO
	Content template.HTML
	Excerpt template.HTML //the excerpt of the content to show on the index page
}

// Index builds the site index/home page
func Index(pages map[string]Page) {
	var p Page

	for key, value := range pages {
		log.Printf("%v => %v", key, value)
		//cb = append(cb, value.Excerpt)
	}

	t, _ := template.ParseFiles("templates/index.html")
	d := "_site"
	index := d + "/index.html"
	fmt.Println(index)
	file, err := os.Create(index)

	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(file)

	p.Title = "Home"
	t.Execute(w, p)
	w.Flush()

}

// Read the file line by line and find the frontmatter
// Return a Page struct if successful
func frontmatter(f string) (Page, error) {

	file, err := os.Open(f)
	if err != nil {
		fmt.Printf("Unable to open file: %v, %v", file, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	const fm = "---"

	start, end := false, false

	//Definitely should use a better method than this @TODO
	b := []string{}
	sb := []string{} //stripped buffer. ie. No frontmatter
	excerpt := sb
	for scanner.Scan() {

		t := scanner.Text()

		if start && end { //already done with FM. Add to stripped buffer

			if t == "<!--more-->" {
				log.Printf("Found more tag in %v\n", f)
				excerpt = sb
				log.Printf("excerpt length: %v\n", len(excerpt))
			}

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

	var p Page

	//Parse the actual YAML inside the front matter
	s := strings.Join(b, "\n")
	e := yaml.Unmarshal([]byte(s), &p)
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
		p.Content = template.HTML(s)

		if len(excerpt) > 0 {
			s = strings.Join(excerpt, "\n")
			md = blackfriday.MarkdownCommon([]byte(s))
			s = string(md[:])
			p.Excerpt = template.HTML(s)

			log.Printf("EX PAGE: %+v\n", p)
		}
	}

	if ext == ".html" {
		p.Content = template.HTML(strings.Join(sb, "\n")) //join the array together, convert to HTML
	}

	return p, e
}

func Server(p string) {
	fmt.Println("Listening on port", p)
	log.Fatal(http.ListenAndServe(":"+p, http.FileServer(http.Dir("./_site/"))))
}
