package tsrender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type TSMarhshaler interface {
	MarshalTS() (res string, err error)
}

type TSMarshalerWrapper struct {
	TSMarhshaler
}

type Tag struct {
	Name    string
	Props   map[string]any // map of property name to anything, which can be marshaled for typescript.
	Content []LiteralTagContentEntry
}

type TagContentType uint8

const (
	TagTypeDefault       TagContentType = 0 // Acts like TagTypeString if content is not empty or like TagTypeOmit if content's empty.
	TagTypeString        TagContentType = 1 // json.Marshal sanitized string
	TagTypeOmit          TagContentType = 2 // Omits content and emits tag given as self-closed.
	TagTypeStringLiteral TagContentType = 3 // A "`" quotes js string, which can evalute JS expressions.
	TagTypeRaw           TagContentType = 4 // Raw tag. Copies and pastes content in between ts tag.
)

func (tct TagContentType) String() string {
	switch tct {
	case TagTypeDefault:
		return "default"
	case TagTypeOmit:
		return "omit"
	case TagTypeString:
		return "string"
	case TagTypeStringLiteral:
		return "quoteString"
	case TagTypeRaw:
		return "raw"
	default:
		return ""
	}
}
func (tct *TagContentType) fromString(text string) (err error) {
	switch text {
	case TagTypeDefault.String():
		*tct = TagTypeDefault
	case TagTypeOmit.String():
		*tct = TagTypeOmit
	case TagTypeString.String():
		*tct = TagTypeString
	case TagTypeStringLiteral.String():
		*tct = TagTypeStringLiteral
	case TagTypeRaw.String():
		*tct = TagTypeRaw
	default:
		err = fmt.Errorf("unknown tag content type: %s", text)
	}
	return
}

// MarshalJSON marshals the enum as a quoted json string
func (s TagContentType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *TagContentType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	s.fromString(j)
	return nil
}

type TagContent []TagContentEntry

func (tc TagContent) Render(ctx context.Context) (res string, err error) {
	isFirsIter := true
	for _, e := range tc {
		if !isFirsIter {
			res += "\n"
		}
		isFirsIter = false

		var immRes string
		immRes, err = e.Render(ctx)
		if err != nil {
			return
		}
		res += immRes
	}
	return
}

type TagContentEntry interface {
	Render(ctx context.Context) (res string, err error)
}

type LiteralTagContentEntry struct {
	Type    TagContentType
	Content string
}

func jsonSanitize(text string) (res string) {
	imm, err := json.Marshal(text)
	if err != nil {
		panic(err)
	}
	return string(imm)
}

func softSanitize(text string) (res string) {
	return strings.ReplaceAll(text, "\n", "\\n")
}

func (tce *LiteralTagContentEntry) Render(ctx context.Context) (res string, err error) {
	switch tce.Type {
	case TagTypeDefault:
		if len(tce.Content) == 0 {
			return
		}
		res = fmt.Sprintf("{%s}", jsonSanitize(tce.Content))
		return
	case TagTypeOmit:
		return
	case TagTypeString:
		res = fmt.Sprintf("{%s}", jsonSanitize(tce.Content))
		return
	case TagTypeStringLiteral:
		res = fmt.Sprintf("{`%s`}", softSanitize(tce.Content))
		return
	case TagTypeRaw:
		res = tce.Content
		return
	default:
		err = fmt.Errorf("handmd/util/tsrender: invalid tag type: %d", tce.Type)
		return
	}

}
