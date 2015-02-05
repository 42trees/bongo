package main

import (
	"flag"
	"fmt"
	"github.com/42trees/bongo"
)

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

	if *projectDir != "" {
		bongo.NewProject(*projectDir)
		return
	}

	if *helpFlag || *versionFlag {
		bongo.Help()
		return
	}

	if *serverFlag {
		bongo.Server(*port)
		return
	}

	bongo.Build(contentPath)

}
