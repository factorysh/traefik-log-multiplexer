package output

import (
	"github.com/factorysh/traefik-log-multiplexer/api"
)

var Outputs map[string]OutputFactory

func init() {
	if Outputs == nil {
		Outputs = make(map[string]OutputFactory)
	}
}

type OutputFactory func(map[string]interface{}) (api.Output, error)
