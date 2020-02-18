package output

import (
	"github.com/influxdata/tail"
)

var Outputs map[string]NewOutput

type NewOutput func(map[string]interface{}) (Output, error)

// Output send logs somewhere
type Output interface {
	Read(project string, line *tail.Line)
	RemoveProject(project string)
}
