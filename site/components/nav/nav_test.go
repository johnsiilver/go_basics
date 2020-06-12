package nav

import (
	"context"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/html"
)

func TestNav(t *testing.T) {
	conf := &config.VideoFiles{
		{
			Index: 0,
			Name:  "Video 0",
			URL:   "/video/0",
		},
		{
			Index: 1,
			Name:  "Video 1",
			URL:   "/video/1",
		},
	}

	buff := &strings.Builder{}
	pipe := html.NewPipeline(context.Background(), nil, buff)

	n, err := New("nav-bar", conf, nil)
	if err != nil {
		panic(err)
	}

	out, err := ioutil.ReadFile("testdata/want")
	if err != nil {
		panic(err)
	}

	n.Execute(pipe)

	got := onlyOneSpaceChar(buff.String())
	want := onlyOneSpaceChar(string(out))

	if got != want {
		t.Errorf("TestNav: want:\n%s\ngot:\n%s", want, got)
	}
}

var spaceRE = regexp.MustCompile(`\s+`)

func onlyOneSpaceChar(s string) string {
	compare := s
	s = spaceRE.ReplaceAllString(s, " ")
	if s == compare {
		return strings.TrimSpace(s)
	}
	return onlyOneSpaceChar(s)
}
