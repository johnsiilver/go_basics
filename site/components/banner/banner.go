package banner

import (
	"fmt"

	"github.com/johnsiilver/go_basics/site/components/nav"
	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"

	. "github.com/johnsiilver/webgear/html"
)

// New constructs a new component that shows a banner.  Execution requires passing a
// github.com/johnsiilver/go_basics/site/components/nav []MenuItem in the pipeline with key {{name}}-navbar.
func New(name string, conf *config.VideoFiles, options ...component.Option) (*component.Gear, error) {
	nav, err := nav.New(fmt.Sprintf("%s-navbar", name), conf, nil)
	if err != nil {
		return nil, err
	}

	build := builder.NewHTML(&Head{}, &Body{})
	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "banner"}})
	build.Add(
		&Link{Rel: "stylesheet", Href: URLParse("/static/components/banner/banner.css")},
		&A{
			Href: URLParse("/"),
			Elements: []Element{
				&Img{GlobalAttrs: GlobalAttrs{ID: "gopher"}, Src: URLParse("/static/components/banner/scientist.svg")},
			},
		},
		&Span{GlobalAttrs: GlobalAttrs{ID: "title"}, Elements: []Element{TextElement("Go Language Basics")}},
		&Div{
			GlobalAttrs: GlobalAttrs{ID: "navHolder"},
			Elements:    []Element{&Component{GlobalAttrs: GlobalAttrs{ID: "navHolder"}, Gear: nav}},
		},
	)

	options = append(options, component.AddGear(nav))

	gear, err := component.New(name, build.Doc(), options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
