package content

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/html"
)

func TestContent(t *testing.T) {
	conf := &config.VideoFiles{
		&config.VideoFile{
			Notes: []byte("#markdown"),
		},
	}
	req := &http.Request{
		URL: html.URLParse("/video/0"),
	}

	buff := &strings.Builder{}
	pipe := html.NewPipeline(context.Background(), req, buff)

	c, err := New("my-content", conf)
	if err != nil {
		panic(err)
	}

	c.Execute(pipe)

	out, err := ioutil.ReadFile("testdata/want")
	if err != nil {
		panic(err)
	}

	space := regexp.MustCompile(`\s+`)
	got := strings.TrimSpace(space.ReplaceAllString(buff.String(), " "))
	want := strings.TrimSpace(space.ReplaceAllString(string(out), " "))

	if got != want {
		t.Errorf("TestContent: want:\n%s\ngot:\n%s", want, got)
	}
}
