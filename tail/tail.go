package tail

import (
	"github.com/influxdata/tail"
)

type Tail struct {
	tailer *tail.Tail
}

func New() *Tail {
	return &Tail{}
}
