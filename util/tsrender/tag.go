package tsrender

import (
	"context"
	"fmt"
	"io"
)

func marshalTagProps(ctx context.Context, props map[string]any, w io.Writer) (err error) {
	for k, v := range props {
		switch typedV := v.(type) {
		// TODO(teawithsand): add more types here
		case int:
			_, err = w.Write([]byte(fmt.Sprint("%s={%s}\n", k, jsonSanitize(typedV))))
			if err != nil {
				return
			}
		case int64:
			_, err = w.Write([]byte(fmt.Sprint("%s={%s}\n", k, jsonSanitize(typedV))))
			if err != nil {
				return
			}
		case string:
			_, err = w.Write([]byte(fmt.Sprint("%s={%s}\n", k, jsonSanitize(typedV))))
			if err != nil {
				return
			}
		case LiteralTagContent: // This is HACK. It should be string literal or sth like that
			var res string
			res, err = typedV.Render(ctx)
			if err != nil {
				return
			}

			_, err = w.Write([]byte(fmt.Sprint("%s=%s\n", k, res)))
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("unsupported type for tag props: %T", v)
			return
		}
	}
	return
}

// Simple TSX tag, which has either no or only string child.
type SimpleTag struct {
	Name    string
	Props   map[string]any // map of property name to anything, which can be marshaled for typescript.
	Content LiteralTagContent
}

func (t SimpleTag) Render(ctx context.Context, w io.Writer) (err error) {
	sc, err := t.Content.IsSelfClosing()
	if err != nil {
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("<%s ", t.Name)))
	if err != nil {
		return
	}

	err = marshalTagProps(ctx, t.Props, w)
	if err != nil {
		return
	}

	if sc {
		_, err = w.Write([]byte("/>"))
		if err != nil {
			return
		}
	} else {
		_, err = w.Write([]byte(">"))
		if err != nil {
			return
		}

		var content string
		content, err = t.Content.Render(ctx)
		if err != nil {
			return
		}

		_, err = w.Write([]byte(content))
		if err != nil {
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("</%s>", t.Name)))
		if err != nil {
			return
		}
	}

	return
}
