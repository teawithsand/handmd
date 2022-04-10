package tsrender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type TagLiteralType uint8

const (
	TagTypeDefault       TagLiteralType = 0 // Acts like TagTypeString if content is not empty or like TagTypeOmit if content's empty.
	TagTypeString        TagLiteralType = 1 // json.Marshal sanitized string
	TagTypeOmit          TagLiteralType = 2 // Omits content and emits tag given as self-closed.
	TagTypeStringLiteral TagLiteralType = 3 // A "`" quotes js string, which can evalute JS expressions.
	TagTypeRaw           TagLiteralType = 4 // Raw tag. Copies and pastes content in between ts tag.
)

func (tct TagLiteralType) String() string {
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
func (tct *TagLiteralType) fromString(text string) (err error) {
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
func (s TagLiteralType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *TagLiteralType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	s.fromString(j)
	return nil
}

type LiteralTagContent struct {
	Type    TagLiteralType
	Content string
}

func (tce *LiteralTagContent) IsSelfClosing() (selfClosing bool, err error) {
	switch tce.Type {
	case TagTypeDefault:
		selfClosing = len(tce.Content) == 0
		return
	case TagTypeOmit:
		selfClosing = true
		return
	default:
		return
	}
}

func (tce *LiteralTagContent) Render(ctx context.Context) (res string, err error) {
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
		res = fmt.Sprintf("{`%s`}", mildSanitize(tce.Content))
		return
	case TagTypeRaw:
		res = tce.Content
		return
	default:
		err = fmt.Errorf("handmd/util/tsrender: invalid tag type: %d", tce.Type)
		return
	}

}
