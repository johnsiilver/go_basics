package content

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html"

	"github.com/russross/blackfriday/v2"
)

type extractor struct {
	conf *config.VideoFiles
}

// DataFunc implements component.DataFunc().
func (e extractor) DataFunc(r *http.Request) (interface{}, error) {
	indexStr := path.Base(r.URL.Path)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, fmt.Errorf("video component could not local URL information on the video: %s", err)
	}

	if index < 0 || index > len(*e.conf) {
		return nil, fmt.Errorf("video component received video index %d that isn't valid", index)
	}

	return (*e.conf)[index], nil
}

func (e extractor) videoFile(pipe html.Pipeline) []html.Element {
	vf, ok := pipe.GearData.(*config.VideoFile)
	if !ok {
		pipe.Error(fmt.Errorf("content component found the content, expected *config.VideoFile, got %T", pipe.GearData))
		return nil
	}

	return []html.Element{html.TextElement(blackfriday.Run(vf.Notes))}
}

// New constructs a new component that shows a video based on the url.
func New(name string, conf *config.VideoFiles, options ...component.Option) (*component.Gear, error) {
	extract := extractor{conf}

	var doc = &html.Doc{
		Body: &html.Body{
			Elements: []html.Element{
				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/components/content/content.css")},
				&html.Div{
					GlobalAttrs: html.GlobalAttrs{ID: "content", Class: "markdown"},
					Elements: []html.Element{
						html.Dynamic(extract.videoFile),
					},
				},
			},
		},
	}

	options = append(options, component.ApplyDataFunc(extract.DataFunc))

	gear, err := component.New(name, doc, options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
