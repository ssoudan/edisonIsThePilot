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
* @Last Modified time: 2015-09-29 12:32:24
 */

package stepper

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type input struct {
	duration JSONDuration
	heading  float64
	step     float64
}

type JSONDuration time.Duration

func (d JSONDuration) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%f\"", time.Duration(d).Seconds())
	return []byte(stamp), nil
}

type point struct {
	timestamp      JSONTime
	fix_time       string
	fix_date       string
	course         float64
	speed          float64
	delta_steering float64
	latitude       nmea.LatLong
	longitude      nmea.LatLong
	validity       bool
}

const (
	UNDEFINED = iota

	ARMED   = iota
	GO      = iota
	RUNNING = iota
	DONE    = iota
	ABORTED = iota
)

type state int

func (d state) MarshalJSON() ([]byte, error) {
	field := "UNKNOWN"
	switch d {
	case ARMED:
		field = "ARMED"
	case GO:
		field = "GO"
	case RUNNING:
		field = "RUNNING"
	case DONE:
		field = "DONE"
	case ABORTED:
		field = "ABORTED"
	}
	return []byte(field), nil
}

type plan struct {
	state        state
	start        JSONTime
	test_type    string
	plot_command string
	input        input
	points       []point
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
	return &Stepper{shutdownChan: make(chan interface{}), plan: plan{state: UNDEFINED}}
}

type message struct {
	step     float64
	duration time.Duration
}

func NewStep(step float64, duration time.Duration) interface{} {
	return message{step: step, duration: duration}
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

	if d.plan.state == UNDEFINED {
		d.plan.state = ARMED
		d.plan.test_type = fmt.Sprintf("bump test of %f", m.step)
		d.plan.input.duration = JSONDuration(m.duration)
		d.plan.input.step = m.step
	}

}

func (d *Stepper) processGPSMessage(m pilot.GPSFeedBackAction) {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()

	switch d.plan.state {

	case GO:
		d.plan.input.heading = m.Heading
		d.plan.state = RUNNING
		d.plan.points = []point{{
			timestamp:      JSONTime(now),
			fix_date:       m.Date,
			fix_time:       m.Time,
			course:         m.Heading,
			speed:          m.Speed,
			delta_steering: d.plan.input.step,
			latitude:       m.Latitude,
			longitude:      m.Longitude,
			validity:       m.Validity,
		}}
		d.plan.start = JSONTime(now)

		// send message to steering
		d.steeringChan <- steering.NewMessage(d.plan.input.step)

	case RUNNING:
		d.plan.points = append(d.plan.points, point{
			timestamp:      JSONTime(now),
			course:         m.Heading,
			speed:          m.Speed,
			delta_steering: 0,
			latitude:       m.Latitude,
			longitude:      m.Longitude,
			validity:       m.Validity,
		})

		if now.After(time.Time(d.plan.start).Add(time.Duration(d.plan.input.duration))) {
			// Time is up
			d.plan.state = DONE
		}
	case DONE:
		// TODO(ssoudan) write to file
		log.Info("Done with %#v", d.plan)

		// TODO(ssoudan) tell the pilot the calibration test is over and data can be collected
	}

}

func (d *Stepper) enable() {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()

	switch d.plan.state {
	case ARMED:
		d.plan.state = GO
	}
}

func (d *Stepper) disable() {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()

	d.plan.state = ABORTED

}

// Shutdown sets all the state to down and notify the handlers
func (d *Stepper) Shutdown() {

	d.shutdownChan <- 1
	<-d.shutdownChan
}

func (d *Stepper) shutdown() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.plan.state != DONE {
		d.plan.state = ABORTED
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
	d.mu.RLock()
	defer d.mu.RUnlock()
	json.NewEncoder(w).Encode(d.plan)
}
