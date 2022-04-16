package httpext

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/teawithsand/handmd/util/iter"
)

// Note: values stored in these are assumed to be valid
type CSPEntry struct {
	Name   string
	Values []string
}

func (e *CSPEntry) mergeWith(other CSPEntry) {
	if e.Name != other.Name {
		panic("handmd/util/httpext: CSPEntry name mismatch when mergin")
	}
	e.Values = append(e.Values, other.Values...)
}

func (e *CSPEntry) Render() (res string) {
	res, err := iter.JoinString(context.Background(),
		iter.Map(iter.Slice(e.Values), func(ctx context.Context, data string) (string, error) {
			return fmt.Sprintf("'%s'", data), nil
		}),
		", ",
	)
	if err != nil {
		panic(err)
	}

	res = e.Name + " " + res

	return
}

type CSPBuilder struct {
	entries map[string]CSPEntry
}

const CSPHeaderName = "Content-Security-Policy"

func (b *CSPBuilder) init() {
	if b.entries == nil {
		b.entries = map[string]CSPEntry{}
	}
}

func (b *CSPBuilder) AddEntry(e CSPEntry) {
	b.init()

	other := b.entries[e.Name]
	if other.Name == "" {
		b.entries[e.Name] = e
	} else {
		other.mergeWith(e)
		b.entries[e.Name] = other
	}
}

func (b *CSPBuilder) Render(ctx context.Context, w io.Writer) (err error) {
	isFirst := true
	for _, e := range b.entries {
		if !isFirst {
			_, err = w.Write([]byte("; "))
			if err != nil {
				return
			}
		}
		isFirst = false

		renderedEntry := e.Render()

		_, err = w.Write([]byte(renderedEntry))
		if err != nil {
			return
		}
	}
	return
}

func (b *CSPBuilder) RenderSimple() (res string, err error) {
	buf := bytes.NewBuffer(nil)
	err = b.Render(context.Background(), buf)
	if err != nil {
		return
	}

	res = buf.String()

	return
}
