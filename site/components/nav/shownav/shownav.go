package main

import (
	"flag"

	"github.com/johnsiilver/go_basics/site/components/nav"
	"github.com/johnsiilver/go_basics/site/config"

	"github.com/johnsiilver/webgear/component/viewer"
)

var (
	port = flag.Int("port", 8080, "The port to run the server on")
)

func main() {
	conf := &config.VideoFiles{
		&config.VideoFile{
			Index: 0,
			Name:  "Google",
			URL:   "http://google.com",
		},
	}

	nav, err := nav.New("nav-component", conf, nil)
	if err != nil {
		panic(err)
	}

	v := viewer.New(
		8080,
		nav,
		viewer.BackgroundColor("black"),
		viewer.ServeOtherFiles("../../../", []string{".css", ".jpg", ".svg", ".png"}),
	)

	v.Run()
}
