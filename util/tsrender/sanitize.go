package tsrender

import (
	"encoding/json"
	"strings"
)

type jsonSanitizeSafe interface {
	~string | ~int | ~int64 | ~int32 | ~uint64 | ~uint32 | ~float32 | ~float64 // TODO(teawithsand): rest of types here
}

func jsonSanitize[T jsonSanitizeSafe](data T) (res string) {
	imm, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(imm)
}

func mildSanitize(text string) (res string) {
	// note: this is kind of crappy, it should sanitize only those "`" chars, which are not already sanitized.
	// now we replace "\`" to "\\`", which is not valid thing to do(no escaping in this strings, chars are chars)
	return strings.ReplaceAll("`", "\\`", strings.ReplaceAll(text, "\n", "\\n"))
}
