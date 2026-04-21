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
		"teamInitials": func(team string) string {
			words := strings.Fields(team)
			if len(words) == 0 {
				return "CC"
			}
			if len(words) == 1 {
				return strings.ToUpper(firstRunes(words[0], 2))
			}
			var b strings.Builder
			for _, word := range words {
				if word == "" {
					continue
				}
				b.WriteString(strings.ToUpper(firstRunes(word, 1)))
			}
			value := b.String()
			if len(value) > 3 {
				return value[len(value)-3:]
			}
			return value
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

func firstRunes(value string, count int) string {
	var b strings.Builder
	for i, r := range value {
		if i >= count {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}
