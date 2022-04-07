package httpext

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/CacheMW-Control
// for caching options description
type CacheMW struct {
	ForceDisable bool // Forces disable all caching for resource, overrides other options

	Infinite             bool          // Do whatever to cache resource for eternity. Useful for static resources. Overrides other options.
	MaxAge               time.Duration // Aligned to second, ignored when zero
	SharedMaxAge         time.Duration // Aligned to second, ignored when zero
	StaleWhileRevalidate time.Duration // Aligned to second, ignored when zero
	StaleIfError         time.Duration // Aligned to second, ignored when zero

	MustRevalidate   bool
	IsPublic         bool
	IsMustUnderstand bool
	IsNoTransform    bool

	VaryAll bool // sets vary to "*", overrides vary option
	Vary    []string

	Immutable bool
}

func (c *CacheMW) getCacheControl() (cache string) {
	if c.ForceDisable {
		return "max-age=0, s-maxage=0, no-store"
	}

	var elems []string

	if c.Infinite {
		elems = append(elems, "max-age=31536000")
		elems = append(elems, "s-maxage=31536000")
	} else {
		if c.MaxAge != 0 {
			elems = append(elems, fmt.Sprintf("max-age=%d", int64(c.MaxAge.Seconds())))
		}

		if c.SharedMaxAge != 0 {
			elems = append(elems, fmt.Sprintf("s-maxage=%d", int64(c.SharedMaxAge.Seconds())))
		}

		if c.StaleWhileRevalidate != 0 {
			elems = append(elems, fmt.Sprintf("stale-while-revalidate=%d", int64(c.StaleWhileRevalidate.Seconds())))
		}

		if c.StaleIfError != 0 {
			elems = append(elems, fmt.Sprintf("stale-if-error=%d", int64(c.StaleIfError.Seconds())))
		}

		if c.MustRevalidate {
			elems = append(elems, "must-revalidate")
		}
	}

	if c.IsPublic {
		elems = append(elems, "public")
	} else {
		elems = append(elems, "private")
	}

	if c.IsMustUnderstand {
		elems = append(elems, "must-understand")
	}

	if c.Immutable {
		elems = append(elems, "immutable")
	}

	if c.IsNoTransform {
		elems = append(elems, "no-transform")
	}

	cache = strings.Join(elems, ", ")
	return
}

func (c *CacheMW) varyString() (vary string) {
	if c.VaryAll {
		vary = "*"
		return
	}

	vary = strings.Join(c.Vary, ", ")
	return
}

func (c *CacheMW) Apply(h http.Handler) http.Handler {
	vary := c.varyString()
	cacheControl := c.getCacheControl()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(cacheControl) > 0 {
			w.Header().Set("Cache-Control", cacheControl)
		}
		if len(vary) > 0 {
			w.Header().Set("Vary", vary)
		}
		h.ServeHTTP(w, r)
	})
}
