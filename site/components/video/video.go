// Package video provides the display of vimeo videos from a list of videos that is passed to the component on init.
package video

import (
	"fmt"
	"path"
	"strconv"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html"
)

// video has controls for dealing with the videos.
type video struct {
	conf *config.VideoFiles
}

// prev looks at the URL, finds if we have previous videos in the list and if so, displays a link that loads the previous video.
func (v video) prev(pipe html.Pipeline) []html.Element {
	indexStr := path.Base(pipe.Req.URL.Path)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		pipe.Error(fmt.Errorf("video component could not locate URL information on the video: %s", err))
		return nil
	}

	if index < 1 {
		return nil
	}

	return []html.Element{
		&html.Span{
			GlobalAttrs: html.GlobalAttrs{
				ID:    "prevVideo",
				Class: "videoControls",
			},
			Elements: []html.Element{
				&html.A{
					Elements: []html.Element{html.TextElement("<")},
					Href:     fmt.Sprintf("/video/%d", index-1),
				},
			},
		},
	}
}

// next looks at the URL, finds if we have more videos in the list and if so, displays a link that loads the next video.
func (v video) next(pipe html.Pipeline) []html.Element {
	indexStr := path.Base(pipe.Req.URL.Path)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		pipe.Error(fmt.Errorf("video component could not locate URL information on the video: %s", err))
		return nil
	}

	if index+1 >= len(*v.conf) {
		return nil
	}

	return []html.Element{
		&html.Span{
			GlobalAttrs: html.GlobalAttrs{
				ID:    "nextVideo",
				Class: "videoControls",
			},
			Elements: []html.Element{
				&html.A{
					Elements: []html.Element{html.TextElement(">")},
					Href:     fmt.Sprintf("/video/%d", index+1),
				},
			},
		},
	}
}

// display looks at the URL, finds the video in our config and creates an iframe according to Vimeo specs.
func (v video) display(pipe html.Pipeline) []html.Element {
	indexStr := path.Base(pipe.Req.URL.Path)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		pipe.Error(fmt.Errorf("video component could not local URL information on the video: %s", err))
		return nil
	}

	if index < 0 || index > len(*v.conf) {
		pipe.Error(fmt.Errorf("video component received video index %d that isn't valid", index))
		return nil
	}

	videoConf := (*v.conf)[index]

	return []html.Element{
		&html.IFrame{
			GlobalAttrs: html.GlobalAttrs{
				ID:    "videoSrc",
				Style: "border:none;",
			},
			Src:             html.URLParse(videoConf.URL),
			Allow:           "autoplay; fullscreen",
			AllowFullscreen: true,
		},
	}
}

// New constructs a new component that shows a video based on the url.
func New(name string, conf *config.VideoFiles, options ...component.Option) (*component.Gear, error) {
	vc := video{conf}

	var doc = &html.Doc{
		Body: &html.Body{
			Elements: []html.Element{
				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/components/video/video.css")},
				html.Dynamic(vc.prev),
				html.Dynamic(vc.display),
				html.Dynamic(vc.next),
			},
		},
	}

	gear, err := component.New(name, doc, options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
