package about

import (
	"html/template"

	"github.com/johnsiilver/go_basics/site/components/banner"
	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/html"
)

const (
	bannerGearName = "banner-component"
)

type aboutMeSection struct {
	Title string
	Value []string
}

type aboutMe struct {
	elements []html.Element
}

func newAboutMe(sections []aboutMeSection) aboutMe {
	elements := []html.Element{
		&html.H{
			Level: 1,
			Elements: []html.Element{
				html.TextElement("A Little About Me"),
			},
		},
		&html.Img{
			GlobalAttrs: html.GlobalAttrs{ID: "me"},
			Src:         html.URLParse("/static/pages/about/me.png"),
		},
	}

	for _, section := range sections {
		elements = append(
			elements,
			&html.H{
				Level: 2,
				Elements: []html.Element{
					html.TextElement(section.Title),
				},
			},
		)
		for i := 0; i < len(section.Value); i++ {
			elements = append(
				elements,
				&html.Span{
					GlobalAttrs: html.GlobalAttrs{Class: "textStyle"},
					Elements: []html.Element{
						html.TextElement(section.Value[i]),
					},
				},
			)
			if i+1 < len(section.Value) {
				elements = append(elements, &html.BR{})
			}
		}
	}
	return aboutMe{elements}
}

func (a aboutMe) dynamic(pipe html.Pipeline) []html.Element {
	return a.elements
}

// New constructs the page for "about" the site.
func New(conf *config.VideoFiles) (*html.Doc, error) {
	ban, err := banner.New(bannerGearName, conf)
	if err != nil {
		return nil, err
	}

	doc := &html.Doc{
		Head: &html.Head{
			Elements: []html.Element{
				&html.Title{TagValue: html.TextElement("Go Language Basics")},
				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/pages/about/about.css")},
				&html.Link{Rel: "stylesheet", Href: html.URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap")},
			},
		},
		Body: &html.Body{
			Elements: []html.Element{
				ban,
				&html.Component{GlobalAttrs: html.GlobalAttrs{ID: "banner"}, TagType: template.HTMLAttr(ban.Name())},
				&html.Div{
					GlobalAttrs: html.GlobalAttrs{ID: "mainPane"},
					Elements: []html.Element{
						&html.Div{
							GlobalAttrs: html.GlobalAttrs{
								ID: "aboutSite",
							},
							Elements: []html.Element{
								&html.H{
									Level:    1,
									Elements: []html.Element{html.TextElement("About The Site")},
								},
								&html.P{
									Elements: []html.Element{
										html.TextElement("Golang Basics is for developers wanting an introduction into programming in Go."),
									},
								},
								&html.P{
									Elements: []html.Element{
										html.TextElement("It is built off of several years experience teaching Go around the world for Google."),
									},
								},
								&html.P{
									Elements: []html.Element{
										html.TextElement(
											"This class will teach the basics of the language so that you can then begin exploring " +
												"Go tooling and advanced concepts that many great Gophers (what we call Go developers) have written.",
										),
									},
								},
							},
						},
						&html.Div{
							GlobalAttrs: html.GlobalAttrs{ID: "aboutMe"},
							Elements: []html.Element{
								html.Dynamic(newAboutMe(
									[]aboutMeSection{
										{
											"Name",
											[]string{"John Doak"},
										},
										{
											"Occupation",
											[]string{"Principal SWE Manager, Microsoft Azure"},
										},
										{
											"Formerly",
											[]string{
												"Staff SRE, Google",
												"Staff Network Systems Engineer, Google (I was the first network focused SRE)",
												"Network Engineer, Google",
												"Network Engineer, LucasFilm/LucasArts/ILM",
											},
										},
									},
								).dynamic),
								&html.H{
									Level: 2,
									Elements: []html.Element{
										html.TextElement("Websites"),
									},
								},
								&html.A{
									Href: "http://www.gophersre.com",
									Elements: []html.Element{
										html.TextElement("gophersre.com"),
									},
								},
								&html.A{
									Href: "http://www.obscuredworld.com",
									Elements: []html.Element{
										html.TextElement("obscuredworld"),
									},
								},
								&html.A{
									Href: "http://www.linkedin.com/in/johngdoak/",
									Elements: []html.Element{
										html.TextElement("LinkedIn"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return doc, nil
}
