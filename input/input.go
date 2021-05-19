package input

import "github.com/factorysh/traefik-log-multiplexer/api"

var Inputs map[string]InputFactory

type InputFactory func(map[string]interface{}) (api.Input, error)
