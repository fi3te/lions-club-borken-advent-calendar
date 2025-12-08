package config

import (
	"fmt"
	"time"
)

type SmtpConfig struct {
	Host     string
	Username string
	Password string
}

func (s SmtpConfig) validate() error {
	if s.Host == "" {
		return errRequired("email smtp host")
	}
	if s.Username == "" {
		return errRequired("email smtp username")
	}
	if s.Password == "" {
		return errRequired("email smtp password")
	}
	return nil
}

type EmailConfig struct {
	Sender struct {
		Address string `yaml:"address"`
		Name    string `yaml:"name"`
	} `yaml:"sender"`
	Recipients []string `yaml:"recipients"`
	Smtp       SmtpConfig
}

func (e EmailConfig) validate() error {
	if e.Sender.Address == "" {
		return errRequired("email sender address")
	}
	if e.Sender.Name == "" {
		return errRequired("email sender name")
	}
	if e.Recipients == nil || len(e.Recipients) == 0 {
		return errRequired("email recipients")
	}
	return e.Smtp.validate()
}

type Config struct {
	Url                           string      `yaml:"url"`
	ContentCssPath                string      `yaml:"content-css-path"`
	Location                      string      `yaml:"location"`
	Time                          string      `yaml:"time"`
	FirstRetryAfterErrorInSeconds int64       `yaml:"first-retry-after-error-in-seconds"`
	Email                         EmailConfig `yaml:"email"`
}

func (cfg *Config) validate() error {
	if cfg.Url == "" {
		return errRequired("url")
	}
	if cfg.ContentCssPath == "" {
		return errRequired("content-css-path")
	}
	if cfg.Time == "" {
		return errRequired("time")
	}
	if _, err := parseLocation(cfg.Location); err != nil {
		return errGeneral("location", err)
	}
	if _, err := parseTime(cfg.Time); err != nil {
		return errGeneral("time", err)
	}
	if cfg.FirstRetryAfterErrorInSeconds < 1 {
		return errGtZero("first retry after error in seconds")
	}
	return cfg.Email.validate()
}

func (cfg *Config) NextExecution() time.Time {
	t, _ := parseTime(cfg.Time)
	now := time.Now()
	execution := time.Date(
		now.Year(), now.Month(), now.Day(),
		t.Hour(), t.Minute(), 0, 0,
		cfg.getLocation(),
	)
	if execution.After(now) {
		return execution
	} else {
		return execution.AddDate(0, 0, 1)
	}
}

func (cfg *Config) getLocation() *time.Location {
	l, _ := parseLocation(cfg.Location)
	return l
}

func parseLocation(s string) (*time.Location, error) {
	return time.LoadLocation(s)
}

func parseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}

func (cfg *Config) FirstRetryAfterError() time.Duration {
	return intToDuration(cfg.FirstRetryAfterErrorInSeconds, time.Second)
}

func intToDuration(value int64, unit time.Duration) time.Duration {
	return time.Duration(value * int64(unit))
}

func errRequired(description string) error {
	return fmt.Errorf("attribute '%s' must be specified", description)
}

func errGtZero(description string) error {
	return fmt.Errorf("attribute '%s' must be greater than zero", description)
}

func errGeneral(description string, err error) error {
	return fmt.Errorf("attribute '%s' is invalid: %v", description, err)
}
