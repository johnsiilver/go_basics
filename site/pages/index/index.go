// Package index implements the index page for the site.
package index

import (
	"github.com/johnsiilver/go_basics/site/components/banner"
	"github.com/johnsiilver/go_basics/site/components/content"
	"github.com/johnsiilver/go_basics/site/components/video"
	"github.com/johnsiilver/go_basics/site/config"

	. "github.com/johnsiilver/webgear/html"
)

const (
	bannerGearName  = "banner-component"
	videoGearName   = "video-component"
	contentGearName = "content-component"
)

// New creates a new Page object that can have the .Doc called to render the index page.
func New(conf *config.VideoFiles) (*Doc, error) {
	bannerGear, err := banner.New(bannerGearName, conf)
	if err != nil {
		return nil, err
	}
	videoGear, err := video.New(videoGearName, conf)
	if err != nil {
		return nil, err
	}
	contentGear, err := content.New(contentGearName, conf)
	if err != nil {
		return nil, err
	}

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "UTF-8"},
				&Title{TagValue: TextElement("Go Language Basics")},
				&Link{Rel: "stylesheet", Href: URLParse("/static/pages/index/index.css")},
				&Link{Href: URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
				&Link{Href: URLParse("https://fonts.googleapis.com/css2?family=Nanum+Gothic&display=swap"), Rel: "stylesheet"},
			},
		},
		Body: &Body{
			Elements: []Element{
				bannerGear, // This causes the code to render.
				videoGear,
				contentGear,
				&Component{GlobalAttrs: GlobalAttrs{ID: "banner"}, Gear: bannerGear},
				&Component{GlobalAttrs: GlobalAttrs{ID: "mainPane"}, Gear: videoGear},
				&Component{GlobalAttrs: GlobalAttrs{ID: "rightPane"}, Gear: contentGear},
			},
		},
	}

	return doc, nil
}
