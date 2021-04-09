package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/go_basics/site/pages/about"
	"github.com/johnsiilver/go_basics/site/pages/index"
	"github.com/johnsiilver/webgear/handlers"

	"github.com/caddyserver/certmagic"
)

var (
	port  = flag.Int("port", 8081, "The port to run the server on. Only works if debug==true")
	debug = flag.Bool("debug", false, "If the server is in developer debug mode")
)

func indexRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/video/0", 301)
}

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	host, p, _ := net.SplitHostPort(r.Host)
	if host == "" {
		host = r.Host
	}
	log.Printf("redirecting from host %s:%s to %s:443", host, p, host)
	u.Scheme="https"
	http.Redirect(w,r,"https://"+r.Host+r.URL.Path, http.StatusMovedPermanently)
}

func httpRedirectServ() {
	s := &http.Server{
		Addr:           ":80",
		Handler:        http.HandlerFunc(httpRedirect),
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	flag.Parse()

	var tlsListen net.Listener
	if !*debug {
		certmagic.DefaultACME.Agreed = true
		certmagic.DefaultACME.Email = "johnsiilver@gmail.com"
		//certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA

		// If we want to back to binding to a specific port and do port forwarding.
		/*
		certmagic.HTTPChallengePort = 8081
		certmagic.TLSALPNChallengePort = 8081
		tlsConfig, err := certmagic.TLS([]string{"gophersre.com", "golangsre.com"})
		if err != nil {
			log.Fatal(err)
		}

		tlsListen, err = tls.Listen("tcp", fmt.Sprintf(":%d", *port), tlsConfig)
		if err != nil {
			log.Fatal(err)
		}
		*/
		var err error
		tlsListen, err = certmagic.Listen([]string{"www.golangbasics.com", "golangbasics.com"})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("port 80 redirector to 443 started")
		go httpRedirectServ()
	}

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
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if !*debug {
		server.Addr = fmt.Sprintf(":%d", *port)
	}

	log.Printf("http server serving on :%d", *port)

	if *debug {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.Serve(tlsListen))
	}
}
