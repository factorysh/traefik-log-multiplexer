package tpl

import (
	"fmt"
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

func (t *Template) Execute(data map[string]interface{}) ([]byte, error) {
	valueSize := 0
	for _, k := range t.vars {
		v, ok := data[k].(string)
		if !ok {
			return nil, fmt.Errorf("only string is handled : %s=>%v", k, data[k])
		}
		valueSize += len(v)
	}
	response := make([]byte, valueSize+t.chunkSize)
	start := 0
	for i := 0; i < len(t.vars); i++ {
		copy(response[start:], t.chunks[i])
		start += len(t.chunks[i])
		v, _ := data[t.vars[i]].(string)
		copy(response[start:], v)
		start += len(v)
	}
	copy(response[start:], t.chunks[len(t.chunks)-1])
	return response, nil
}
