package video

import (
	"context"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/johnsiilver/go_basics/site/config"

	"github.com/johnsiilver/webgear/html"
)

func TestVideo(t *testing.T) {
	conf := &config.VideoFiles{
		{
			Index: 0,
			Name:  "Intro",
			URL:   "https://player.vimeo.com/video/422489803",
		},
		{
			Index: 1,
			Name:  "Packages",
			URL:   "https://player.vimeo.com/video/422491427",
		},
		{
			Index: 2,
			Name:  "types",
			URL:   "https://player.vimeo.com/video/422276275",
		},
	}

	tests := []struct {
		desc     string
		req      *http.Request
		wantFile string
	}{
		{
			desc: "On page 0, so only next button",
			req: &http.Request{
				URL: html.URLParse("blah.com/video/0"),
			},
			wantFile: "page0",
		},
		{
			desc: "On page 1, so prev and next button",
			req: &http.Request{
				URL: html.URLParse("blah.com/video/1"),
			},
			wantFile: "page1",
		},
		{
			desc: "On page 2, so only prev button",
			req: &http.Request{
				URL: html.URLParse("blah.com/video/2"),
			},
			wantFile: "page2",
		},
	}

	b, err := New("my-video", conf)
	if err != nil {
		panic(err)
	}

	for _, test := range tests {
		buff := &strings.Builder{}
		pipe := html.NewPipeline(context.Background(), test.req, buff)

		out, err := ioutil.ReadFile(filepath.Join("testdata", test.wantFile))
		if err != nil {
			panic(err)
		}

		space := regexp.MustCompile(`\s+`)

		b.Execute(pipe)

		got := strings.TrimSpace(space.ReplaceAllString(string(buff.String()), " "))
		want := strings.TrimSpace(space.ReplaceAllString(string(out), " "))

		if strings.TrimSpace(got) != want {
			t.Errorf("TestVideo(%s): want:\n%s\ngot:\n%s", test.desc, want, got)
		}
	}
}
