package state

import "time"

const FileName = "state.yml"

type State struct {
	LastSuccessfulExecution time.Time `yaml:"last-successful-execution"`
}
