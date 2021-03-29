package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/go_basics/site/pages/about"
	"github.com/johnsiilver/go_basics/site/pages/index"

	"github.com/johnsiilver/webgear/handlers"
)

var (
	port  = flag.Int("port", 8081, "The port to run the server on")
	debug = flag.Bool("debug", false, "If the server is in developer debug mode")
)

func indexRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/video/0", 301)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	log.Println("CommandLine Flags:")
	flag.CommandLine.VisitAll(
		func(fl *flag.Flag) {
			log.Printf("\tFlag: %s: %v", fl.Name, fl.Value)
		},
	)

	conf := &config.VideoFiles{}
	if err := conf.FromFile("etc/videos.config"); err != nil {
		panic(err)
	}

	var opts []handlers.Option
	if *debug {
		opts = append(opts, handlers.DoNotCache())
	}

	h := handlers.New(opts...)
	h.ServeFilesWorkingDir([]string{".css", ".jpg", ".svg", ".png", ".ico"})

	index, err := index.New(conf)
	if err != nil {
		panic(err)
	}

	about, err := about.New(conf)
	if err != nil {
		panic(err)
	}

	h.HTTPHandler("/", http.HandlerFunc(indexRedirect))
	h.Handle("/video/", index)
	h.Handle("/about/", about)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("http server serving on :%d", *port)

	log.Fatal(server.ListenAndServe())
}
