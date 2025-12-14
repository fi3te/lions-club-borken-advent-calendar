package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/config"
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/control"
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/domain"
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/mail"
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/state"
)

var cfg *config.Config
var st *state.State

func main() {
	log.Println("Started application")

	var err error
	if state.FileExists() {
		st, err = state.ReadState()
		if err != nil {
			log.Printf("Failed to read state: %v\n", err)
		} else {
			log.Println("Read state from file")
		}
	} else {
		log.Println("State file does not exist and will be created.")
	}

	if st == nil {
		st = &state.State{}
	}

	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("An error occurred while reading config: %v\n", err)
	}
	log.Println("Read configuration file")
	log.Println("URL: " + cfg.Url)
	log.Println("Content CSS path: " + cfg.ContentCssPath)
	log.Println("Sender: " + cfg.Email.Sender.Name + " <" + cfg.Email.Sender.Address + ">")
	log.Println("Recipients: " + strings.Join(cfg.Email.Recipients, ", "))
	todaysPlannedExecutionTime := cfg.TodaysPlannedExecutionTime()
	log.Println("Today's planned execution time: " + todaysPlannedExecutionTime.String())
	log.Println("Location: " + cfg.Location)
	firstRetryAfterError := cfg.FirstRetryAfterError()
	log.Println("First retry after error: " + firstRetryAfterError.String())

	nextRegularExecutionTime, runImmediately := getNextRegularExecutionTime(todaysPlannedExecutionTime)
	if runImmediately {
		log.Println("Today's planned execution time has passed and no execution has taken place. This will be rectified immediately.")
		if err := executeTask(time.Now()); err != nil {
			log.Printf("The unscheduled execution failed and will not be repeated: %v\n", err)
		}
	}

	stopTimerChan := make(chan bool)

	go control.RunPeriodicTimer(nextRegularExecutionTime, 24*time.Hour, firstRetryAfterError, stopTimerChan, executeTask)

	control.WaitForInterrupt()

	log.Println("Stopped application")
	stopTimerChan <- true
}

func getNextRegularExecutionTime(todaysPlannedExecutionTime time.Time) (nextRegularExecutionTime time.Time, runImmediately bool) {
	now := time.Now()
	if todaysPlannedExecutionTime.After(now) {
		nextRegularExecutionTime = todaysPlannedExecutionTime
		log.Println("Today's planned execution time is targeted.")
	} else {
		if st.LastSuccessfulExecution.Before(todaysPlannedExecutionTime) {
			runImmediately = true
		}
		nextRegularExecutionTime = todaysPlannedExecutionTime.AddDate(0, 0, 1)
		log.Println("Planned next regular execution for tomorrow: " + nextRegularExecutionTime.String())
	}
	return
}

func executeTask(t time.Time) error {
	log.Println("Executing task...")
	door, err := domain.GetAdventCalendarDoor(cfg, t)
	if err != nil {
		return err
	}

	err = mail.SendMail(cfg.Email, "Lions Club Adventskalender - Tag "+strconv.Itoa(t.Day()), door.HtmlContent)
	if err != nil {
		return err
	}

	st.LastSuccessfulExecution = t
	err = state.WriteState(st)
	if err != nil {
		log.Printf("Failed to write state: %v\n", err)
	} else {
		log.Println("Wrote state to file")
	}

	// err for writing state is not relevant
	return nil
}
