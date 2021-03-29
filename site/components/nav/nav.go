// Package nav provides a component gear for navigation.
package nav

import (
	"fmt"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/component"

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
	var doc = &Doc{
		Body: &Body{
			Elements: append(
				append(
					[]Element{
						&Link{Rel: "stylesheet", Href: URLParse("/static/components/nav/nav.css")},
					},
					scriptsToElements(scripts)...,
				),
				&Nav{
					GlobalAttrs: GlobalAttrs{
						ID: "nav",
					},
					Elements: []Element{
						&Ul{
							GlobalAttrs: GlobalAttrs{ID: "navList"},
							Elements: []Element{
								&Li{
									Elements: []Element{
										&A{ // Top level button.
											Href:        URLParse("#"),
											Elements:    []Element{TextElement("Sections")},
											GlobalAttrs: GlobalAttrs{Class: "title"},
										},
										&Ul{
											Elements: []Element{Dynamic(menuList{conf}.Items)},
										},
									},
								},
								&Li{
									Elements: []Element{
										&A{ // Top level button.
											Href:        URLParse("/about"),
											Elements:    []Element{TextElement("About")},
											GlobalAttrs: GlobalAttrs{Class: "title"},
										},
									},
								},
							},
						},
					},
				},
			),
		},
	}

	gear, err := component.New(name, doc, options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
