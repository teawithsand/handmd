package simplesite

import (
	"bytes"
	"context"
	"errors"
	"html"
	"io"
)

type TagClose uint8

const (
	SimpleTagClose TagClose = 0
	NoTagClose     TagClose = 1
	SelfTagClose   TagClose = 2
)

type HTMLAttribute struct {
	Name  string
	Value string
}

// Simple way to render single HTML tag.
type HTMLTag struct {
	Name          string
	Attributes    []HTMLAttribute
	Content       string
	Close         TagClose
	RawAttributes []string
}

func (t *HTMLTag) Render(ctx context.Context, w io.Writer) (err error) {
	_, err = w.Write([]byte("<"))
	if err != nil {
		return
	}
	_, err = w.Write([]byte(t.Name))
	if err != nil {
		return
	}

	for _, attr := range t.Attributes {
		_, err = w.Write([]byte(" "))
		if err != nil {
			return
		}

		_, err = w.Write(
			[]byte(html.EscapeString(attr.Name)),
		)
		if err != nil {
			return
		}

		_, err = w.Write([]byte("=\""))
		if err != nil {
			return
		}

		_, err = w.Write(
			[]byte(html.EscapeString(attr.Value)),
		)
		if err != nil {
			return
		}

		_, err = w.Write([]byte("\""))
		if err != nil {
			return
		}
	}

	for _, attr := range t.RawAttributes {
		_, err = w.Write([]byte(" "))
		if err != nil {
			return
		}
		_, err = w.Write([]byte(attr))
		if err != nil {
			return
		}
	}

	switch t.Close {
	case SelfTagClose:
		if len(t.Content) != 0 {
			err = errors.New("handmd/util/simplesite: expected no content for SelfTagClose")
			return
		}
		_, err = w.Write([]byte("/>"))
		if err != nil {
			return
		}
	case NoTagClose:
		if len(t.Content) != 0 {
			err = errors.New("handmd/util/simplesite: expected no content for NoTagClose")
			return
		}
		_, err = w.Write([]byte(">"))
		if err != nil {
			return
		}
	case SimpleTagClose:
		fallthrough
	default:
		_, err = w.Write([]byte(">"))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(html.EscapeString(t.Content)))
		if err != nil {
			return
		}
		_, err = w.Write([]byte("</"))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(t.Name))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(">"))
		if err != nil {
			return
		}
	}
	return
}

func (t *HTMLTag) RenderSimple() (res string, err error) {
	b := bytes.NewBuffer(nil)
	err = t.Render(context.Background(), b)
	if err != nil {
		return
	}
	res = b.String()
	return
}
