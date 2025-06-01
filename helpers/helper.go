package helpers

import (
	"io"
	"text/template"
)

func RenderTemplate(name string) func(interface{}, *template.Template, io.Writer) error {
	return func(ctx interface{}, tmpl *template.Template, wr io.Writer) error {
		return tmpl.ExecuteTemplate(wr, name, ctx)
	}
}
