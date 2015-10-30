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
* @Last Modified time: 2015-10-29 23:05:01
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
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	"github.com/ssoudan/edisonIsThePilot/steering"

	"github.com/adrianmo/go-nmea"
)

var log = logger.Log("stepper")

// Input is the definition of a step
type Input struct {
	Duration types.JSONDuration
	Heading  float64
	Step     float64
}

// Point is the structure collected when a step is RUNNING
type Point struct {
	Timestamp     types.JSONTime
	FixTime       string
	FixDate       string
	Course        float64
	Speed         float64
	DeltaSteering float64
	Latitude      nmea.LatLong
	Longitude     nmea.LatLong
	Validity      bool
}

// state value
const (
	UNDEFINED = iota
	ARMED     = iota
	GO        = iota
	PREPARING = iota
	RUNNING   = iota
	DONE      = iota
	ABORTED   = iota
)

// State of a step
type State int

// MarshalJSON does the JSON serialization of a State
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
	case PREPARING:
		field = "\"PREPARING\""
	}
	return []byte(field), nil
}

type plan struct {
	State       State
	Start       types.JSONTime
	TestType    string
	PlotCommand string
	Input       Input
	Points      []Point
	Description string
}

// Stepper is a component that execute steering plans and collect the position as it runs
type Stepper struct {
	mu   sync.RWMutex
	plan plan // protected by mu

	// channels
	inputChan    chan interface{}
	steeringChan chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

// New creates a new Stepper component
func New() *Stepper {
	return &Stepper{shutdownChan: make(chan interface{}), plan: plan{State: UNDEFINED}}
}

type message struct {
	step        float64
	duration    time.Duration
	description string
}

// NewStep creates a new step
func (s Stepper) NewStep(step float64, duration time.Duration, description string) {
	s.inputChan <- newStep(step, duration, description)
}

func newStep(step float64, duration time.Duration, description string) interface{} {
	return message{step: step, duration: duration, description: description}
}

// SetInputChan sets the channel where the Stepper will be getting new step messages from
func (s *Stepper) SetInputChan(c chan interface{}) {
	s.inputChan = c
}

// SetPanicChan sets the channel where panics are sent
func (s *Stepper) SetPanicChan(c chan interface{}) {
	s.panicChan = c
}

// SetSteeringChan sets the channel where the Stepper will send steering order to
func (s *Stepper) SetSteeringChan(c chan interface{}) {
	s.steeringChan = c
}

type enableAction struct {
}

type disableAction struct {
}

// Enable the autopilot
func (s *Stepper) Enable() error {
	s.inputChan <- enableAction{}
	return nil
}

// Disable the autopilot
func (s *Stepper) Disable() error {
	s.inputChan <- disableAction{}
	return nil
}

func (s *Stepper) processNewStepMessage(m message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.plan.State == UNDEFINED {
		s.plan.State = ARMED
		s.plan.TestType = fmt.Sprintf("bump test of %f", m.step)
		s.plan.Input.Duration = types.JSONDuration(m.duration)
		s.plan.Input.Step = m.step
		s.plan.Description = m.description
	}

}

func (s *Stepper) processGPSMessage(m pilot.GPSFeedBackAction) {
	// Update the state
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	switch s.plan.State {

	case GO:
		s.plan.Input.Heading = m.Heading
		s.plan.State = PREPARING
		s.plan.Start = types.JSONTime(now)
		s.plan.Points = []Point{{
			Timestamp:     types.JSONTime(now),
			FixDate:       m.Date,
			FixTime:       m.Time,
			Course:        m.Heading,
			Speed:         m.Speed,
			DeltaSteering: 0,
			Latitude:      m.Latitude,
			Longitude:     m.Longitude,
			Validity:      m.Validity,
		}}

	case PREPARING:
		if time.Time(s.plan.Start).Add(time.Duration(s.plan.Input.Duration)).Before(now) {
			// send message to steering -- that's where we punch the system
			s.steeringChan <- steering.NewMessage(s.plan.Input.Step, true)

			s.plan.Points = append(s.plan.Points, Point{
				Timestamp:     types.JSONTime(now),
				FixDate:       m.Date,
				FixTime:       m.Time,
				Course:        m.Heading,
				Speed:         m.Speed,
				DeltaSteering: s.plan.Input.Step,
				Latitude:      m.Latitude,
				Longitude:     m.Longitude,
				Validity:      m.Validity,
			})

			s.plan.State = RUNNING
		} else {
			s.plan.Points = append(s.plan.Points, Point{
				Timestamp:     types.JSONTime(now),
				FixDate:       m.Date,
				FixTime:       m.Time,
				Course:        m.Heading,
				Speed:         m.Speed,
				DeltaSteering: 0,
				Latitude:      m.Latitude,
				Longitude:     m.Longitude,
				Validity:      m.Validity,
			})
		}

	case RUNNING:
		s.plan.Points = append(s.plan.Points, Point{
			Timestamp:     types.JSONTime(now),
			FixDate:       m.Date,
			FixTime:       m.Time,
			Course:        m.Heading,
			Speed:         m.Speed,
			DeltaSteering: 0,
			Latitude:      m.Latitude,
			Longitude:     m.Longitude,
			Validity:      m.Validity,
		})

		if time.Time(s.plan.Start).Add(time.Duration(2 * s.plan.Input.Duration)).Before(now) {
			// Time is up
			// write to file
			f, err := os.Create("/tmp/systemCalibration-" + time.Now().Format(time.RFC3339))
			if err != nil {
				log.Error("Failed to open experiment log file: %v", err)
				break
			}
			defer f.Close()
			json.NewEncoder(f).Encode(s.plan)

			// tell the pilot the calibration test is over and data can be collected
			log.Notice("The bump test is over. You can disable the autopilot, stop the program and start another test.")

			s.steeringChan <- steering.NewMessage(0, false)

			s.plan.State = DONE
		}
	case DONE:
		// Nothing
		log.Notice("DONE")
	}

}

func (s *Stepper) enable() {
	// Update the state
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.plan.State {
	case ARMED:
		s.plan.State = GO
	}
}

func (s *Stepper) disable() {
	// Update the state
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.plan.State != DONE {
		s.plan.State = ABORTED
	}

}

// Shutdown sets all the state to down and notify the handlers
func (s *Stepper) Shutdown() {

	s.shutdownChan <- 1
	<-s.shutdownChan
}

func (s *Stepper) shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.plan.State != DONE {
		s.plan.State = ABORTED
	}

	close(s.shutdownChan)
}

// Start the event loop of the Stepper component
func (s *Stepper) Start() {

	go func() {

		defer func() {
			if r := recover(); r != nil {
				s.panicChan <- r
			}
		}()

		for {
			select {
			case m := <-s.inputChan:
				switch m := m.(type) {
				case pilot.GPSFeedBackAction:
					s.processGPSMessage(m)
				case message:
					s.processNewStepMessage(m)
				case enableAction:
					s.enable()
				case disableAction:
					s.disable()
				case error:
					log.Error("Received an error: %v", m)
				}
			case <-s.shutdownChan:
				s.shutdown()

				return
			}

		}
	}()

}

// CalibrationEndpoint is a rest endpoint to get the current status of a plan
func (s *Stepper) CalibrationEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := json.NewEncoder(w).Encode(s.plan)
	if err != nil {
		log.Error("Failed to encode %v", err)
	}
}
