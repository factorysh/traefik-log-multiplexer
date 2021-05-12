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
	fmt.Printf("  %v\n  %v\n%v", ts, line, meta)
	return nil
}

func (s *StdoutOutput) Close() error {
	return nil
}

func StdoutOutputFactory(cfg map[string]interface{}) (api.Output, error) {
	return &StdoutOutput{}, nil
}
