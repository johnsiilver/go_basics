// Package index implements the index page for the site.
package index

import (
	"html/template"

	"github.com/johnsiilver/go_basics/site/components/banner"
	"github.com/johnsiilver/go_basics/site/components/content"
	"github.com/johnsiilver/go_basics/site/components/video"
	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/html"
)

const (
	bannerGearName  = "banner-component"
	videoGearName   = "video-component"
	contentGearName = "content-component"
)

// New creates a new Page object that can have the .Doc called to render the index page.
func New(conf *config.VideoFiles) (*html.Doc, error) {
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

	doc := &html.Doc{
		Head: &html.Head{
			Elements: []html.Element{
				&html.Meta{Charset: "UTF-8"},
				&html.Title{TagValue: html.TextElement("Go Language Basics")},
				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/pages/index/index.css")},
				&html.Link{Href: html.URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
				&html.Link{Href: html.URLParse("https://fonts.googleapis.com/css2?family=Nanum+Gothic&display=swap"), Rel: "stylesheet"},
			},
		},
		Body: &html.Body{
			Elements: []html.Element{
				bannerGear, // This causes the code to render.
				videoGear,
				contentGear,
				&html.Component{GlobalAttrs: html.GlobalAttrs{ID: "banner"}, TagType: template.HTMLAttr(bannerGear.Name())},
				&html.Component{GlobalAttrs: html.GlobalAttrs{ID: "mainPane"}, TagType: template.HTMLAttr(videoGear.Name())},
				&html.Component{GlobalAttrs: html.GlobalAttrs{ID: "rightPane"}, TagType: template.HTMLAttr(contentGear.Name())},
			},
		},
	}

	return doc, nil
}
