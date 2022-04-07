package defines

import (
	"context"
	"io"

	"github.com/teawithsand/handmd/util/iter"
	"github.com/teawithsand/handmd/util/tsrender"
)

type PostContentEntry struct {
	Type string `json:"type" toml:"type" yaml:"type"` // purely informative for renderer, ignored in output(for now ignored everywhere)

	// Simplified version of typescript tag here
	// tag name
	// It's content
	Tag         string                  `json:"tag" toml:"tag" yaml:"tag"`
	Content     string                  `json:"content" toml:"content" yaml:"content"`
	ContentType tsrender.TagLiteralType `json:"contentType" toml:"content_type" yaml:"content_type"`
	Props       map[string]string       `json:"props,omitempty" toml:"props,omitempty" yaml:"props,omitempty"`
}

type PostContent struct {
	Imports []tsrender.Import  `json:"imports,omitempty" toml:"imports,omitempty" yaml:"imports,omitempty"`
	Entries []PostContentEntry `json:"entries" toml:"entries" yaml:"entries"`
}

// ExtractPostContentText extracts text from iterable of PostContentEntry.
// It's useful for making index of any kind.
//
// Passed iterator may not be parallel. It must use only one goroutine.
func ExtractPostContentText(ctx context.Context, entries iter.Iterable[PostContentEntry], w io.Writer) (err error) {
	isFirst := false

	return entries.Iterate(ctx, iter.Receiver[PostContentEntry](func(ctx context.Context, e PostContentEntry) (err error) {
		if len(e.Content) == 0 {
			return
		}

		// In general, write all content that may be here
		// Images or stuff are passed by props, rather than content, so it's ok.
		_, err = w.Write([]byte(e.Content))
		if err != nil {
			return
		}

		if !isFirst {
			_, err = w.Write([]byte("\n"))
			if err != nil {
				return
			}
		}

		isFirst = false

		return
	}))
}
