// Package nav provides a component gear for navigation.
package nav

import (
	"fmt"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"

	. "github.com/johnsiilver/webgear/html"
)

// menuList provides a method, Items() that implements Dynamic. It takes a list of MenuItems from the
// Pipeline and renders them as Li objects wrapping A representing a menu.
type menuList struct {
	conf *config.VideoFiles
}

// Items implements DynamicFunc.
func (m menuList) Items(pipe Pipeline) []Element {
	elements := []Element{}

	for _, item := range *m.conf {
		elements = append(
			elements,
			&Li{
				Elements: []Element{
					&A{
						Href:     URLParse(fmt.Sprintf("/video/%d", item.Index)),
						Elements: []Element{TextElement(item.Name)},
					},
				},
			},
		)
	}

	return elements
}

func scriptsToElements(scripts []*Script) []Element {
	n := make([]Element, len(scripts))
	for _, s := range scripts {
		n = append(n, s)
	}
	return n
}

// New constructs a new component that shows a nav bar.
func New(name string, conf *config.VideoFiles, scripts []*Script, options ...component.Option) (*component.Gear, error) {
	build := builder.NewHTML(&Head{}, &Body{})
	build.Into(
		&Div{
			GlobalAttrs: GlobalAttrs{ID: "container"},
		},
	)
	build.Add(&Link{Rel: "stylesheet", Href: URLParse("/static/components/nav/nav.css")})
	build.Add(scriptsToElements(scripts)...)
	build.Into(&Nav{GlobalAttrs: GlobalAttrs{ID: "nav"}})

	build.Into(&Ul{GlobalAttrs: GlobalAttrs{ID: "navList"}})

	build.Into(&Li{})

	// Top level button "Sections"
	build.Add(&A{GlobalAttrs: GlobalAttrs{Class: "title"}, Href:        URLParse("#"), Elements:    []Element{TextElement("Sections")}})
	build.Add(&Ul{GlobalAttrs: GlobalAttrs{Class: "items", Style: "visibility: hidden;"}, Elements: []Element{Dynamic(menuList{conf}.Items)}})
	build.Up() // At navList level

	// Top level button "About"
	build.Into(&Li{})

	build.Add(&A{GlobalAttrs: GlobalAttrs{Class: "title"}, Href:        URLParse("/about"), Elements:    []Element{TextElement("About")}})

	gear, err := component.New(name, build.Doc(), options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
