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
		Close: simplesite.SimpleTagClose,
	}

	res, err := tag.RenderSimple()
	if err != nil {
		t.Error(err)
		return
	}

	if res != "<script src=\"asdf.js\"</script>" {
		t.Error("expected different result")
		return
	}
}
