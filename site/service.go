package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/johnsiilver/go_basics/site/config"
)

var (
	port  = flag.Int("port", 8081, "The port to run the server on")
	debug = flag.Bool("debug", false, "If the server is in developer debug mode")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	log.Println("CommandLine Flags:")
	flag.CommandLine.VisitAll(
		func(fl *flag.Flag) {
			log.Printf("\tFlag: %s: %v", fl.Name, fl.Value)
		},
	)

	mux := http.NewServeMux()
	h, err := newHandlers(mux, "etc/videos.config", *debug)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

// VideoDot is the pipeline object that is passed to the templates.
type VideosDot struct {
	// Files is a list of all video files for the site.
	Files *config.VideoFiles
	// Index is the current video file that the page should represent.
	Index int
}

// Next tells what the next video will be.  -1 if we are on the last video.
func (v VideosDot) Next() int {
	if v.Index+1 == len(*v.Files) {
		return -1
	}
	return v.Index + 1
}

// Prev tells what the previous video will be.  -1 if we are on the first video.
func (v VideosDot) Prev() int {
	if v.Index-1 < 0 {
		return -1
	}
	return v.Index - 1
}

// handler collects all our http handlers together in a single object with helpers.
type handlers struct {
	mux http.Handler

	configPath string
	files      atomic.Value // *config.VideoFiles
	tmpl       atomic.Value // *template.Template
	fs         dotFileHidingFileSystem

	gzPool sync.Pool

	debug bool
}

// newHandlers is the constructor for handlers.
func newHandlers(mux *http.ServeMux, configPath string, debug bool) (*handlers, error) {
	vConfig := &config.VideoFiles{}
	if err := vConfig.FromFile(configPath); err != nil {
		return nil, err
	}

	h := handlers{
		configPath: configPath,
		gzPool: sync.Pool{
			New: func() interface{} {
				w := gzip.NewWriter(ioutil.Discard)
				return w
			},
		},
		fs:    dotFileHidingFileSystem{http.Dir("assets")},
		debug: debug,
	}
	h.files.Store(vConfig)
	h.parseTemplates()

	// Go through config, create a handler for each video.
	for _, config := range *vConfig {
		mux.Handle(fmt.Sprintf("/video/%d", config.Index), http.HandlerFunc(h.video))
	}

	// Load all our assets into /assets/{{dir}}/{{file name}}
	mux.Handle("/assets/", http.FileServer(h.fs))

	mux.Handle("/about", http.HandlerFunc(h.about))

	// Handle everything else.
	mux.Handle("/", http.HandlerFunc(h.indexRedirect))

	h.mux = h.handlerWrap(mux)

	return &h, nil
}

func (h *handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// handlerWrap wraps all our page handlers with handlers that setup debugging, compression, etc... if needed.
func (h *handlers) handlerWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.doDebug(w); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.gzip(next).ServeHTTP(w, r)
	})
}

func (h *handlers) indexRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/video/0", 301)
}

func (h *handlers) about(w http.ResponseWriter, r *http.Request) {
	tmpl, dot := h.getData(-1)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := tmpl.ExecuteTemplate(w, "about.tmpl", dot); err != nil {
		panic(err)
	}
}

// video is the http handler for relative URL /video/%d .
func (h *handlers) video(w http.ResponseWriter, r *http.Request) {
	if err := h.doDebug(w); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	indexStr := path.Base(r.URL.Path)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, dot := h.getData(index)

	if err := tmpl.ExecuteTemplate(w, "index.tmpl", dot); err != nil {
		panic(err)
	}
}

// gzip wraps an http.Handler so that it's content is gzipped.
func (h *handlers) gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := h.gzPool.Get().(*gzip.Writer)
		defer h.gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// doDebug causes our templates and files to be reparsed and loaded on every call.  This is expensive
// but great for live reloading of data when developing the site.
func (h *handlers) doDebug(w http.ResponseWriter) error {
	if !h.debug {
		return nil
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := h.parseTemplates(); err != nil {
		return err
	}

	if err := h.files.Load().(*config.VideoFiles).FromFile(h.configPath); err != nil {
		return err
	}

	return nil
}

// parseTemplates parses all of our templates from the directory.
func (h *handlers) parseTemplates() error {
	files, err := ioutil.ReadDir("html")
	if err != nil {
		return err
	}

	fileList := []string{}
	for _, fi := range files {
		if strings.HasSuffix(fi.Name(), ".tmpl") {
			fileList = append(fileList, "html/"+fi.Name())
		}
	}

	tmpl, err := template.ParseFiles(fileList...)
	if err != nil {
		return err
	}

	h.tmpl.Store(tmpl)

	return nil
}

func (h *handlers) getData(index int) (tmpl *template.Template, dot VideosDot) {
	tmpl = h.tmpl.Load().(*template.Template)
	dot = VideosDot{Files: h.files.Load().(*config.VideoFiles), Index: index}

	return tmpl, dot
}

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fis = append(fis, file)
	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fs dotFileHidingFileSystem) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/assets")

	if containsDotFile(name) { // If dot file, return 403 response
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dotFileHidingFile{file}, err
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
