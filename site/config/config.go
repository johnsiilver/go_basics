package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

type VideoFiles []*VideoFile

func (v *VideoFiles) FromFile(p string) error {
	if len(*v) > 0 {
		*v = (*v)[0:0]
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(b))

	_, err = dec.Token()
	if err != nil {
		return err
	}

	for dec.More() {
		vf := &VideoFile{}
		if err := dec.Decode(&vf); err != nil {
			return err
		}
		*v = append(*v, vf)
	}

	_, err = dec.Token()
	if err != nil {
		return err
	}

	v.importNotes(p)

	return nil
}

func (v *VideoFiles) importNotes(p string) {
	dir := filepath.Join(filepath.Dir(p), "notes")
	for _, vf := range *v {
		noteName := filepath.Join(dir, fmt.Sprintf("%d.md", vf.Index))
		b, err := ioutil.ReadFile(noteName)
		if err != nil {
			continue
		}
		vf.Notes = b
	}
}

type VideoFile struct {
	Index int
	Name  string
	URL   string
	Notes []byte
}

func (v *VideoFile) Render() template.HTML {
	return template.HTML(blackfriday.Run(v.Notes))
}
