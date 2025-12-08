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
)

var cfg *config.Config

func main() {
	log.Println("Started application")

	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("An error occurred while reading config: %v\n", err)
	}
	log.Println("URL: " + cfg.Url)
	log.Println("Content CSS path: " + cfg.ContentCssPath)
	log.Println("Sender: " + cfg.Email.Sender.Name + " <" + cfg.Email.Sender.Address + ">")
	log.Println("Recipients: " + strings.Join(cfg.Email.Recipients, ", "))
	firstExecution := cfg.NextExecution()
	log.Println("First execution: " + firstExecution.String())
	log.Println("Location: " + cfg.Location)
	firstRetryAfterError := cfg.FirstRetryAfterError()
	log.Println("First retry after error: " + firstRetryAfterError.String())

	stopTimerChan := make(chan bool)

	go control.RunPeriodicTimer(firstExecution, 24*time.Hour, firstRetryAfterError, stopTimerChan, executeTask)

	control.WaitForInterrupt()

	log.Println("Stopped application")
	stopTimerChan <- true
}

func executeTask(t time.Time) error {
	log.Println("Executing task...")
	door, err := domain.GetAdventCalendarDoor(cfg, t)
	if err != nil {
		return err
	}
	return mail.SendMail(cfg.Email, "Lions Club Adventskalender - Tag "+strconv.Itoa(t.Day()), door.HtmlContent)
}
