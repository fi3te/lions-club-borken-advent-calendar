package control

import (
	"log"
	"time"
)

func RunPeriodicTimer(
	firstExecution time.Time,
	interval time.Duration,
	initialIntervalOnError time.Duration,
	stopTimerChan chan bool,
	callback func(time.Time) error) {

	duration := time.Until(firstExecution)
	timer := time.NewTimer(duration)

	periodicTimer := periodicTimer{
		timer,
		firstExecution,
		interval,
		initialIntervalOnError,
		0,
	}

	for {
		select {
		case t := <-timer.C:
			err := callback(t)
			if err != nil {
				log.Println("Error occurred: " + err.Error())
				periodicTimer.setIntervalOnError()
				interval := periodicTimer.currentIntervalOnError
				timer.Reset(interval)
				log.Printf("Waiting for %s after error...\n", interval)
			} else {
				periodicTimer.currentIntervalOnError = 0
				periodicTimer.planNextExecution()
				nextRegularExecution := periodicTimer.nextRegularExecution
				duration := time.Until(nextRegularExecution)
				timer.Reset(duration)
				log.Printf("Waiting for %s for the next regular execution at %s...\n", duration, nextRegularExecution)
			}
		case <-stopTimerChan:
			timer.Stop()
			return
		}
	}
}
