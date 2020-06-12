package banner

import (
	"fmt"
	"html/template"

	"github.com/johnsiilver/go_basics/site/config"

	"github.com/johnsiilver/go_basics/site/components/nav"
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html"
)

// New constructs a new component that shows a banner.  Execution requires passing a
// github.com/johnsiilver/go_basics/site/components/nav []MenuItem in the pipeline with key {{name}}-navbar.
func New(name string, conf *config.VideoFiles, options ...component.Option) (*component.Gear, error) {
	nav, err := nav.New(fmt.Sprintf("%s-navbar", name), conf, nil)
	if err != nil {
		return nil, err
	}

	doc := &html.Doc{
		Body: &html.Body{
			GlobalAttrs: html.GlobalAttrs{ID: "banner"},
			Elements: []html.Element{
				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/components/banner/banner.css")},
				&html.A{
					Href: "/",
					Elements: []html.Element{
						&html.Img{
							GlobalAttrs: html.GlobalAttrs{ID: "gopher"},
							Src:         html.URLParse("/static/components/banner/scientist.svg"),
						},
					},
				},
				&html.A{
					Href: "/",
					Elements: []html.Element{
						&html.Span{
							GlobalAttrs: html.GlobalAttrs{ID: "title"},
							Elements: []html.Element{
								html.TextElement("Go Language Basics"),
							},
						},
					},
				},
				&html.Component{
					GlobalAttrs: html.GlobalAttrs{ID: "navHolder"},
					TagType:     template.HTMLAttr(nav.Name()),
				},
			},
		},
	}

	options = append(options, component.AddGear(nav))

	gear, err := component.New(name, doc, options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
