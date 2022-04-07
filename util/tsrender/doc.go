// Package tsrender contains utils for rendering typescript files.
package tsrender

import (
	"fmt"
	"strings"
)

type Import struct {
	From          string      `json:"from" toml:"from" yaml:"from"`
	Default       string      `json:"default" toml:"default" yaml:"default"`
	Others        []string    `json:"others" toml:"others" yaml:"others"`
	OthersAliased [][2]string `json:"othersAliased" toml:"others_aliased" yaml:"others_aliased"`
}

func (i *Import) Render() string {
	res := "import "

	if len(i.Default) > 0 {
		res += i.Default
	}
	if len(i.Others) > 0 || len(i.OthersAliased) > 0 {
		res += " { "
		res += strings.Join(i.Others, ", ")
		var aliased []string
		for _, a := range i.OthersAliased {
			aliased = append(aliased, fmt.Sprintf("%s as %s", a[0], a[1]))
		}
		res += strings.Join(aliased, ", ")
		res += " } "

	}

	res += " from "
	res += "\""
	res += i.From
	res += "\""

	return res
}
