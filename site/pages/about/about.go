package about

import (
	"github.com/johnsiilver/go_basics/site/components/banner"
	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/html/builder"

	. "github.com/johnsiilver/webgear/html"
)

const (
	bannerGearName = "banner-component"
)

type aboutMeSection struct {
	Title string
	Value []string
}

type aboutMe struct {
	elements []Element
}

func newAboutMe(sections []aboutMeSection) aboutMe {
	build := builder.NewHTML(&Head{}, &Body{})
	build.Add(&H{Level: 1, Elements: []Element{TextElement("A Little About Me")}})
	build.Add(&Img{GlobalAttrs: GlobalAttrs{ID: "me"}, Src:         URLParse("/static/pages/about/me.png")})

	for _, section := range sections {
		build.Add(&H{Level: 2, Elements: []Element{TextElement(section.Title)}})

		for i := 0; i < len(section.Value); i++ {
			build.Add(&Span{GlobalAttrs: GlobalAttrs{Class: "textStyle"}, Elements: []Element{TextElement(section.Value[i])}})
			if i+1 < len(section.Value) {
				build.Add(&BR{})
			}
		}
	}
	return aboutMe{build.Doc().Body.Elements}
}

func (a aboutMe) dynamic(pipe Pipeline) []Element {
	return a.elements
}

// New constructs the page for "about" the site.
func New(conf *config.VideoFiles) (*Doc, error) {
	ban, err := banner.New(bannerGearName, conf)
	if err != nil {
		return nil, err
	}

	head := Head{
		Elements: []Element{
			&Title{TagValue: TextElement("Go Language Basics")},
			&Link{Rel: "stylesheet", Href: URLParse("/static/pages/about/about.css")},
			&Link{Rel: "stylesheet", Href: URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap")},
			&Link{Rel:"icon", Type: "image/jpg", Href: URLParse("/static/pages/index/go.jpg")},
		},
	}


	build := builder.NewHTML(&head, &Body{})
	build.Add(ban)
	build.Add(&Component{GlobalAttrs: GlobalAttrs{ID: "banner"}, Gear: ban})
	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "mainPane"}})
	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "aboutSite"}})
	build.Add(&H{Level: 1, Elements: []Element{TextElement("About The Site")}})
	build.Add(&P{Elements: []Element{TextElement("Golang Basics is for developers wanting an introduction into programming in Go.")}})
	build.Add(&P{Elements: []Element{TextElement("It is built off of several years experience teaching Go around the world for Google.")}})
	build.Add(
		&P{
			Elements: []Element{
				TextElement("This class will teach the basics of the language so that you can then begin exploring " +
					"Go tooling and advanced concepts that many great Gophers (what we call Go developers) have written."),
			},
		},
	)
	build.Up()

	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "aboutMe"}})
	build.Add(
		Dynamic(
			newAboutMe(
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
			).dynamic,
		),
	)
	build.Add(
		&H{Level: 2, Elements: []Element{TextElement("Websites")}},
		&A{Href: URLParse("http://www.gophersre.com"),Elements: []Element{TextElement("gophersre.com")}},
		&A{Href: URLParse("http://www.obscuredworld.com"), Elements: []Element{TextElement("obscuredworld")}},
		&A{Href: URLParse("http://www.linkedin.com/in/johngdoak/"),Elements: []Element{TextElement("LinkedIn"),}},
	)

	return build.Doc(), nil
}
