package output

import (
	"fmt"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
)

func init() {
	Outputs["stdout"] = StdoutOutputFactory
}

type StdoutOutput struct {
}

func (s *StdoutOutput) Write(ts time.Time, line string, meta map[string]interface{}) error {
	fmt.Println(ts, line, meta)
	return nil
}

func StdoutOutputFactory(cfg map[string]interface{}) (api.Output, error) {
	return &StdoutOutput{}, nil
}
