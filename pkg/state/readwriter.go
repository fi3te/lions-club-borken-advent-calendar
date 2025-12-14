package state

import (
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/common"
)

func FileExists() bool {
	return common.FileExists(FileName)
}

func ReadState() (*State, error) {
	return common.ReadYaml[State](FileName)
}

func WriteState(state *State) error {
	return common.WriteYaml[State](FileName, state)
}
