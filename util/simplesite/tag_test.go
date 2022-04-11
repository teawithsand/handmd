package simplesite_test

import (
	"testing"

	"github.com/teawithsand/handmd/util/simplesite"
)

func TestTag(t *testing.T) {
	tag := simplesite.HTMLTag{
		Name: "script",
		Attributes: []simplesite.HTMLAttribute{
			{
				Name:  "src",
				Value: "asdf.js",
			},
		},
		RawAttributes: []string{"defer"},
		Close:         simplesite.SimpleTagClose,
	}

	res, err := tag.RenderSimple()
	if err != nil {
		t.Error(err)
		return
	}

	if res != "<script src=\"asdf.js\" defer></script>" {
		t.Error("expected different result", res)
		return
	}
}
