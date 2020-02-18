package output

import (
	"github.com/influxdata/tail"
)

// Output send logs somewhere
type Output interface {
	Read(project string, line *tail.Line)
	RemoveProject(project string)
}
