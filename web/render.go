package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

type Renderer struct {
	fsys fs.FS
}

func NewRenderer(fsys fs.FS) (*Renderer, error) {
	return &Renderer{fsys: fsys}, nil
}

func (r *Renderer) HTML(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	baseSrc, err := fs.ReadFile(r.fsys, "base.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageSrc, err := fs.ReadFile(r.fsys, name+".gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	funcs := template.FuncMap{
		"currency": func(cents int) string {
			return fmt.Sprintf("$%.2f", float64(cents)/100)
		},
		"slugjoin": func(parts ...string) string {
			return strings.Join(parts, " ")
		},
	}

	tmpl, err := template.New(name).Funcs(funcs).Parse(string(baseSrc) + "\n" + string(pageSrc))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(buf.Bytes())
}
