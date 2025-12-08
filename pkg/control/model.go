package control

import "time"

type periodicTimer struct {
	timer                  *time.Timer
	nextRegularExecution   time.Time
	interval               time.Duration
	initialIntervalOnError time.Duration
	currentIntervalOnError time.Duration
}

func (pt *periodicTimer) setIntervalOnError() {
	var newInterval time.Duration
	if pt.currentIntervalOnError == 0 {
		newInterval = pt.initialIntervalOnError
	} else {
		newInterval = pt.currentIntervalOnError * 2
	}
	pt.currentIntervalOnError = newInterval
}

func (pt *periodicTimer) planNextExecution() {
	now := time.Now()
	for now.After(pt.nextRegularExecution) {
		pt.nextRegularExecution = pt.nextRegularExecution.Add(pt.interval)
	}
}
