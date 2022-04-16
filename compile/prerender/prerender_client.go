package prerender

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client for https://www.npmjs.com/package/prerender
// It's very basic though.
type Client struct {
	Client          *http.Client
	PrerendererHost string
}

// Returns reader, which contains prerendered URL target given.
func (c *Client) PrerenderURL(ctx context.Context, targetURL string) (result io.ReadCloser, err error) {
	u := url.URL{
		Scheme:   "http",
		Host:     c.PrerendererHost,
		Path:     "/render",
		RawQuery: "url=" + url.PathEscape(targetURL),
	}

	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return
	}

	if r.Response.StatusCode != http.StatusOK {
		err = fmt.Errorf("handmd/compile/prerender: invalid response status: %d", r.Response.StatusCode)
		return
	}

	res, err := c.Client.Do(r)
	if err != nil {
		return
	}

	result = res.Body
	return
}
