package tpl

import (
	"regexp"
)

type Template struct {
	vars      []string
	chunks    [][]byte
	chunkSize int
}

func Parse(txt []byte) (*Template, error) {
	r := regexp.MustCompile(`\$\{\s*([a-zA-Z0-9_.\-]+)\s*\}`)
	m := r.FindAllSubmatchIndex(txt, -1)
	if m == nil {
		return nil, nil
	}
	tpl := &Template{
		vars:   make([]string, len(m)),
		chunks: make([][]byte, len(m)+1),
	}
	start := 0
	for i, j := range m {
		tpl.vars[i] = string(txt[j[2]:j[3]])
		tpl.chunks[i] = txt[start:j[0]]
		tpl.chunkSize += j[0] - start
		start = j[1]
	}
	tpl.chunks[len(tpl.chunks)-1] = txt[start:]
	return tpl, nil
}

func (t *Template) Execute(data map[string]string) ([]byte, error) {
	valueSize := 0
	for _, k := range t.vars {
		valueSize += len(data[k])
	}
	response := make([]byte, valueSize+t.chunkSize)
	start := 0
	for i := 0; i < len(t.vars); i++ {
		copy(response[start:], t.chunks[i])
		start += len(t.chunks[i])
		copy(response[start:], data[t.vars[i]])
		start += len(data[t.vars[i]])
	}
	copy(response[start:], t.chunks[len(t.chunks)-1])
	return response, nil
}
