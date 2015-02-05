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

func Build(c *string) {

	startTime := time.Now()
	dirs, err := ioutil.ReadDir(*c)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, n := range dirs {
		if n.IsDir() { //Run function that builds the posts or pages
			fmt.Println(n.Name())
			p := *c + "/" + n.Name()

			//@TODO check errors
			md, _ := filepath.Glob(p + "/*.md")
			parseFiles(md)

			html, _ := filepath.Glob(p + "/*.html")
			parseFiles(html)

		}
	}

	Index()
	fmt.Printf("Built in %v ms\n", int(1000*time.Since(startTime).Seconds()))
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

		//@TODO don't ignore these
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

func Index() {

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
