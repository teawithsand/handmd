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
	return strings.ReplaceAll(text, "\n", "\\n")
}
