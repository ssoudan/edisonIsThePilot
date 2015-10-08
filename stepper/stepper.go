/*
Copyright 2015 Sebastien Soudan

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy
of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations
under the License.
*/

/*
* @Author: Sebastien Soudan
* @Date:   2015-09-29 10:43:34
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-03 23:01:03
 */

package stepper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	"github.com/ssoudan/edisonIsThePilot/steering"

	"github.com/adrianmo/go-nmea"
)

var log = logger.Log("stepper")

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%d\"", time.Time(t).Unix())
	return []byte(stamp), nil
}

type Input struct {
	Duration JSONDuration
	Heading  float64
	Step     float64
}

type JSONDuration time.Duration

func (d JSONDuration) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%f\"", time.Duration(d).Seconds())
	return []byte(stamp), nil
}

type Point struct {
	Timestamp      JSONTime
	Fix_time       string
	Fix_date       string
	Course         float64
	Speed          float64
	Delta_steering float64
	Latitude       nmea.LatLong
	Longitude      nmea.LatLong
	Validity       bool
}

const (
	UNDEFINED = iota

	ARMED   = iota
	GO      = iota
	RUNNING = iota
	DONE    = iota
	ABORTED = iota
)

type State int

func (d State) MarshalJSON() ([]byte, error) {
	field := "UNKNOWN"
	switch d {
	case ARMED:
		field = "\"ARMED\""
	case GO:
		field = "\"GO\""
	case RUNNING:
		field = "\"RUNNING\""
	case DONE:
		field = "\"DONE\""
	case ABORTED:
		field = "\"ABORTED\""
	}
	return []byte(field), nil
}

type plan struct {
	State        State
	Start        JSONTime
	Test_type    string
	Plot_command string
	Input        Input
	Points       []Point
	Description  string
}

type Stepper struct {
	mu   sync.RWMutex
	plan plan // protected by mu

	// channels
	inputChan    chan interface{}
	steeringChan chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

func New() *Stepper {
	return &Stepper{shutdownChan: make(chan interface{}), plan: plan{State: UNDEFINED}}
}

type message struct {
	step        float64
	duration    time.Duration
	description string
}

func NewStep(step float64, duration time.Duration, description string) interface{} {
	return message{step: step, duration: duration, description: description}
}

func (d *Stepper) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

func (d *Stepper) SetPanicChan(c chan interface{}) {
	d.panicChan = c
}

func (d *Stepper) SetSteeringChan(c chan interface{}) {
	d.steeringChan = c
}

type EnableAction struct {
}

type DisableAction struct {
}

// Enable the autopilot
func (d *Stepper) Enable() error {
	d.inputChan <- EnableAction{}
	return nil
}

// Disable the autopilot
func (d *Stepper) Disable() error {
	d.inputChan <- DisableAction{}
	return nil
}

func (d *Stepper) processNewStepMessage(m message) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.plan.State == UNDEFINED {
		d.plan.State = ARMED
		d.plan.Test_type = fmt.Sprintf("bump test of %f", m.step)
		d.plan.Input.Duration = JSONDuration(m.duration)
		d.plan.Input.Step = m.step
		d.plan.Description = m.description
	}

}

func (d *Stepper) processGPSMessage(m pilot.GPSFeedBackAction) {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()

	switch d.plan.State {

	case GO:
		d.plan.Input.Heading = m.Heading
		d.plan.State = RUNNING
		d.plan.Points = []Point{{
			Timestamp:      JSONTime(now),
			Fix_date:       m.Date,
			Fix_time:       m.Time,
			Course:         m.Heading,
			Speed:          m.Speed,
			Delta_steering: d.plan.Input.Step,
			Latitude:       m.Latitude,
			Longitude:      m.Longitude,
			Validity:       m.Validity,
		}}
		d.plan.Start = JSONTime(now)

		// send message to steering -- that's where we punch the system
		d.steeringChan <- steering.NewMessage(d.plan.Input.Step, true)

	case RUNNING:
		d.plan.Points = append(d.plan.Points, Point{
			Timestamp:      JSONTime(now),
			Fix_date:       m.Date,
			Fix_time:       m.Time,
			Course:         m.Heading,
			Speed:          m.Speed,
			Delta_steering: 0,
			Latitude:       m.Latitude,
			Longitude:      m.Longitude,
			Validity:       m.Validity,
		})

		if time.Time(d.plan.Start).Add(time.Duration(d.plan.Input.Duration)).Before(now) {
			// Time is up
			// write to file
			f, err := os.Create("/tmp/systemCalibration-" + time.Now().Format(time.RFC3339))
			if err != nil {
				log.Error("Failed to open experiment log file: %v", err)
				break
			}
			defer f.Close()
			json.NewEncoder(f).Encode(d.plan)

			// tell the pilot the calibration test is over and data can be collected
			log.Notice("The bump test is over. You can disable the autopilot, stop the program and start another test.")

			d.steeringChan <- steering.NewMessage(0, false)

			d.plan.State = DONE
		}
	case DONE:
		// Nothing
	}

}

func (d *Stepper) enable() {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()

	switch d.plan.State {
	case ARMED:
		d.plan.State = GO
	}
}

func (d *Stepper) disable() {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.plan.State != DONE {
		d.plan.State = ABORTED
	}

}

// Shutdown sets all the state to down and notify the handlers
func (d *Stepper) Shutdown() {

	d.shutdownChan <- 1
	<-d.shutdownChan
}

func (d *Stepper) shutdown() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.plan.State != DONE {
		d.plan.State = ABORTED
	}

	close(d.shutdownChan)
}

// Start the event loop of the Stepper component
func (d *Stepper) Start() {

	go func() {

		defer func() {
			if r := recover(); r != nil {
				d.panicChan <- r
			}
		}()

		for {
			select {
			case m := <-d.inputChan:
				switch m := m.(type) {
				case pilot.GPSFeedBackAction:
					d.processGPSMessage(m)
				case message:
					d.processNewStepMessage(m)
				case EnableAction:
					d.enable()
				case DisableAction:
					d.disable()
				case error:
					log.Error("Received an error: %v", m)
				}
			case <-d.shutdownChan:
				d.shutdown()

				return
			}

		}
	}()

}

func (d *Stepper) CalibrationEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d.mu.RLock()
	defer d.mu.RUnlock()
	err := json.NewEncoder(w).Encode(d.plan)
	if err != nil {
		log.Error("Failed to encode %v", err)
	}
}
