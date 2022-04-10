package httpext

import "net/http"

// See https://developer.mozilla.org/pl/docs/Web/HTTP/Headers/Referrer-Policy
// for documentation.
type ReferrerPolicy uint8

const ReferrerPolicyHeader = "Referrer-Policy"

const (
	NoReferrer                  ReferrerPolicy = 1
	NoReferrerWhenDowngrade     ReferrerPolicy = 2
	Origin                      ReferrerPolicy = 3
	OriginWhenCrossOrigin       ReferrerPolicy = 4
	SameOrigin                  ReferrerPolicy = 5
	StrictOrigin                ReferrerPolicy = 6
	StrictOriginWhenCrossOrigin ReferrerPolicy = 7
	UnsafeUrl                   ReferrerPolicy = 8
)

func (rp ReferrerPolicy) String() string {
	switch rp {
	case NoReferrer:
		return "no-referrer"
	case NoReferrerWhenDowngrade:
		return "no-referrer-when-downgrade"
	case Origin:
		return "origin"
	case OriginWhenCrossOrigin:
		return "origin-when-cross-origin"
	case SameOrigin:
		return "same-origin"
	case StrictOrigin:
		return "strict-origin"
	case StrictOriginWhenCrossOrigin:
		return "strict-origin-when-cross-origin"
	case UnsafeUrl:
		return "unsafe-url"
	default:
		return ""
	}
}

type ReferrerPolicyMW struct {
	Policy ReferrerPolicy
}

func (f ReferrerPolicyMW) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ReferrerPolicyHeader, f.Policy.String())
		h.ServeHTTP(w, r)
	})
}
