// Package nav provides a component gear for navigation.
package nav

import (
	"fmt"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html"
)

// menuList provides a method, Items() that implements html.Dynamic. It takes a list of MenuItems from the
// Pipeline and renders them as html.Li objects wrapping html.A representing a menu.
type menuList struct {
	conf *config.VideoFiles
}

// Items implements html.DynamicFunc.
func (m menuList) Items(pipe html.Pipeline) []html.Element {
	elements := []html.Element{}

	for _, item := range *m.conf {
		elements = append(
			elements,
			&html.Li{
				Elements: []html.Element{
					&html.A{
						Href:     fmt.Sprintf("/video/%d", item.Index),
						Elements: []html.Element{html.TextElement(item.Name)},
					},
				},
			},
		)
	}

	return elements
}

func scriptsToElements(scripts []*html.Script) []html.Element {
	n := make([]html.Element, len(scripts))
	for _, s := range scripts {
		n = append(n, s)
	}
	return n
}

// New constructs a new component that shows a nav bar.
func New(name string, conf *config.VideoFiles, scripts []*html.Script, options ...component.Option) (*component.Gear, error) {
	var doc = &html.Doc{
		Body: &html.Body{
			Elements: append(
				append(
					[]html.Element{
						&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/components/nav/nav.css")},
					},
					scriptsToElements(scripts)...,
				),
				&html.Nav{
					GlobalAttrs: html.GlobalAttrs{
						ID: "nav",
					},
					Elements: []html.Element{
						&html.Ul{
							GlobalAttrs: html.GlobalAttrs{ID: "navList"},
							Elements: []html.Element{
								&html.Li{
									Elements: []html.Element{
										&html.A{ // Top level button.
											Href:        "#",
											Elements:    []html.Element{html.TextElement("Sections")},
											GlobalAttrs: html.GlobalAttrs{Class: "title"},
										},
										&html.Ul{
											Elements: []html.Element{html.Dynamic(menuList{conf}.Items)},
										},
									},
								},
								&html.Li{
									Elements: []html.Element{
										&html.A{ // Top level button.
											Href:        "/about",
											Elements:    []html.Element{html.TextElement("About")},
											GlobalAttrs: html.GlobalAttrs{Class: "title"},
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
