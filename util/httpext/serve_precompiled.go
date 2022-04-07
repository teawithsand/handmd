package httpext

import (
	"net/http"
	"path"
	"strings"
)

type fileSystemFunc func(name string) (http.File, error)

func (f fileSystemFunc) Open(name string) (http.File, error) {
	return f(name)
}

var validCompressedExtensionsWithDots = []string{".js", ".css", ".html", ".json"}
var validCompressedContentType = []string{"text/javascript", "text/css", "text/html", "application/json"}

type ComrpessionAlgoData struct {
	Extension string
	Priority  int
}

func (cad ComrpessionAlgoData) GetPath(name string) string {
	return name + "." + cad.Extension
}

func PrecompressedHandler(baseFs http.FileSystem, algos map[string]ComrpessionAlgoData) http.Handler {
	if algos == nil {
		algos = map[string]ComrpessionAlgoData{
			"gzip": {
				Extension: "gz",
				Priority:  0,
			},
			"br": {
				Extension: "br",
				Priority:  1,
			},
		}
	}

	fileSystems := make(map[string]http.FileSystem)
	handlers := make(map[string]http.Handler)

	for algo, data := range algos {
		data := data
		fileSystems[algo] = fileSystemFunc(func(name string) (http.File, error) {
			return baseFs.Open(data.GetPath(name))
		})

		handlers[algo] = http.FileServer(fileSystems[algo])
	}

	fileSystems[""] = baseFs
	handlers[""] = http.FileServer(baseFs)

	fallbackHandler := handlers[""]

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var encoding string
		maxPriority := -999999999 // min int here
		for _, curEnc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			curEnc = strings.TrimSpace(curEnc)
			if len(curEnc) == 0 {
				continue
			}

			algo, ok := algos[curEnc]
			if ok {
				// use priorities here,
				// since brotil offers better compression than gzip, prefer it
				if algo.Priority >= maxPriority {
					maxPriority = algo.Priority
					encoding = curEnc
				}
			}
		}

		handler := handlers[encoding]
		fs := fileSystems[encoding]

		validExt := false
		var contentType string

		for i, e := range validCompressedExtensionsWithDots {
			if strings.HasSuffix(r.URL.Path, e) {
				validExt = true
				contentType = validCompressedContentType[i]
				break
			}
		}

		if encoding == "" || r.Header.Get("Upgrade") != "" || !validExt {
			fallbackHandler.ServeHTTP(w, r)
			return
		}

		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		file, err := fs.Open(upath)
		if err != nil {
			fallbackHandler.ServeHTTP(w, r)
			return
		}
		file.Close()

		r.Header.Del("Accept-Encoding")

		w.Header().Set("Content-Encoding", encoding)
		if len(contentType) > 0 {
			w.Header().Set("Content-Type", contentType)
		}

		handler.ServeHTTP(w, r)
	})
}
