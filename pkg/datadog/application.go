package datadog

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	DatadogApp interface {
		Get() *DatadogAppBody
		ActivateDowntime() error
		RemoveDownTime() error
	}

	DatadogAppBody struct {
		Name     string
		Downtime time.Time
	}

	datadogAppHandler struct {
		Name      string
		StartDate time.Time
		EndDate   time.Time
	}

	Option func(*datadogAppHandler)
)

func NewDatadogApp(opts ...Option) DatadogApp {
	app := &datadogAppHandler{}
	for _, opt := range opts {
		opt(app)
	}

	return app
}

func WithName(name string) Option {
	return func(d *datadogAppHandler) {
		d.Name = name
	}
}

func WithStartDate(start time.Time) Option {
	return func(d *datadogAppHandler) {
		d.StartDate = start
	}
}

func WithEndDate(end time.Time) Option {
	return func(d *datadogAppHandler) {
		d.EndDate = end
	}
}

func (d *datadogAppHandler) Get() *DatadogAppBody {
	return &DatadogAppBody{Name: d.Name, Downtime: d.StartDate}
}

func (d *datadogAppHandler) ActivateDowntime() error {
	remoteApp := d.Get()
	if remoteApp.Downtime.Unix() != 0 {
		fmt.Printf("%s app already got a downtime", d.Name)
		return nil
	}

	if d.Name == "" {
		return fmt.Errorf("no app name specified")
	}

	if d.EndDate.Before(d.StartDate) {
		return fmt.Errorf("end date is before start date")
	}

	_, err := d.Marshal()
	if err != nil {
		return err
	}

	// Do datadog call

	return nil
}

func (d *datadogAppHandler) RemoveDownTime() error {
	remoteApp := d.Get()
	if remoteApp.Downtime.Unix() == 0 {
		fmt.Printf("no downtime specified for %s, skipping", d.Name)
		return nil
	}

	return nil
}

func (d *datadogAppHandler) Marshal() ([]byte, error) {
	content, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("unable to parse datadog downtime request, %w", err)
	}

	return content, nil
}
