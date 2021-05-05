package tpl

import (
	"io"
	"regexp"
)

type Template struct {
	vars   []string
	chunks []string
}

func Parse(txt string) (*Template, error) {
	r := regexp.MustCompile(`\$\{\s*([a-zA-Z0-9_.\-]+)\s*\}`)
	m := r.FindAllSubmatchIndex([]byte(txt), -1)
	if m == nil {
		return nil, nil
	}
	tpl := &Template{
		vars:   make([]string, len(m)),
		chunks: make([]string, len(m)+1),
	}
	start := 0
	for i, j := range m {
		tpl.vars[i] = txt[j[2]:j[3]]
		tpl.chunks[i] = txt[start:j[0]]
		start = j[1]
	}
	tpl.chunks[len(tpl.chunks)-1] = txt[start:]
	return tpl, nil
}

func (t *Template) Execute(w io.Writer, data map[string]string) error {
	for i := 0; i < len(t.vars); i++ {
		w.Write([]byte(t.chunks[i]))
		w.Write([]byte(data[t.vars[i]]))
	}
	w.Write([]byte(t.chunks[len(t.chunks)-1]))
	return nil
}
