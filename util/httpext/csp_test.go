package httpext_test

import (
	"testing"

	"github.com/teawithsand/handmd/util/httpext"
)

func TestCSP(t *testing.T) {
	b := httpext.CSPBuilder{}

	b.AddEntry(httpext.CSPEntry{
		Name: "script-src",
		Values: []string{
			"nonce-asdfasdfasdf",
		},
	})

	b.AddEntry(httpext.CSPEntry{
		Name: "style-src",
		Values: []string{
			"nonce-1234567890",
		},
	})

	res, err := b.RenderSimple()
	if err != nil {
		t.Error(err)
		return
	}
	if res != "script-src 'nonce-asdfasdfasdf'; style-src 'nonce-1234567890'" && res != "style-src 'nonce-1234567890'; script-src 'nonce-asdfasdfasdf'" {
		t.Error("expected different result csp\n", res)
		return
	}
	return
}
