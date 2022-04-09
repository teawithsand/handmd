package simplesite

import (
	"fmt"
)

// TODO(teawithsand): make render methods there less crappy; sanitize input

func wrapTag(tag, value string) string {
	if len(value) == 0 {
		return ""
	}
	return fmt.Sprintf(` %s="%s"`, tag, value)
}

type Meta struct {
	Name    string
	Content string
}

func (m *Meta) Render() string {
	return fmt.Sprintf("<meta%s%s>", wrapTag("name", m.Name), wrapTag("content", m.Content))
}

type Link struct {
	Rel   string
	Href  string
	Nonce string
}

func (m *Link) Render() string {
	return fmt.Sprintf("<link%s%s%s >", wrapTag("rel", m.Rel), wrapTag("href", m.Href), wrapTag("nonce", m.Nonce))
}

type Script struct {
	Src         string
	Nonce       string
	Integrity   string
	Crossorigin string

	Async bool
	Defer bool
}

func (s *Script) Render() string {
	var t string
	if s.Async && s.Defer {
		t = "async defer"
	} else if s.Async {
		t = "async"
	} else if s.Defer {
		t = "defer"
	}

	return fmt.Sprintf(
		`<script %s%s%s%s%s ></script>`,
		t,
		wrapTag("src", s.Src),
		wrapTag("nonce", s.Nonce),
		wrapTag("crossorigin", s.Crossorigin),
		wrapTag("integrity", s.Integrity),
	)
}

type Input struct {
	ID    string
	Value string
}

func (i *Input) Render() string {
	return fmt.Sprintf(
		`<input %s%s%s ></script>`,
		wrapTag("id", i.ID),
		wrapTag("value", i.Value),
		"type=hidden",
	)
}

type SiteHead struct {
	Lang    string
	Title   string
	Metas   []Meta
	Links   []Link
	Scripts []Script
	Inputs  []Input
}
