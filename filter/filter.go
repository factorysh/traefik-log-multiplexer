package filter

import "github.com/factorysh/traefik-log-multiplexer/api"

var Filters map[string]FilterFactory

func init() {
	if Filters == nil {
		Filters = map[string]FilterFactory{}
	}
}

type FilterFactory func(map[string]interface{}) (api.Filter, error)
