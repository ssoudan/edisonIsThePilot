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
* @Date:   2015-09-20 09:58:02
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 19:30:50
 */

package pilot

import (
	"github.com/ssoudan/edisonIsThePilot/alarm"
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/steering"
)

var log = logger.Log("pilot")

type Pilot struct {
	heading float64 // target heading (set point)
	bound   float64

	alarm      Alarm
	enabled    bool
	headingSet bool

	leds map[string]bool
	pid  Controller

	// channels with the other components
	dashboardChan chan interface{}
	inputChan     chan interface{}
	alarmChan     chan interface{}
	steeringChan  chan interface{}
}

type Controller interface {
	Set(sp float64)
	Update(value float64) float64
}

type Leds map[Led]bool
type Led string

func (p *Pilot) checkHeadingError(headingError float64) Alarm {

	inputStatus := validateInput(p.bound, headingError)
	return computeAlarmStateForInputStatus(p.alarm, inputStatus)

}

func computeHeadingError(heading float64, gpsHeading float64) float64 {

	headingError := gpsHeading - heading
	if headingError > 180. {
		headingError -= 360.
	}

	if headingError <= -180. {
		headingError += 360.
	}

	return headingError
}

func New(controller Controller, bound float64) *Pilot {

	return &Pilot{
		leds:  make(map[string]bool),
		bound: bound,
		pid:   controller}
}

func (p *Pilot) SetDashboardChan(c chan interface{}) {
	p.dashboardChan = c
}

func (p *Pilot) SetInputChan(c chan interface{}) {
	p.inputChan = c
}

func (p *Pilot) SetAlarmChan(c chan interface{}) {
	p.alarmChan = c
}

func (p *Pilot) SetSteeringChan(c chan interface{}) {
	p.steeringChan = c
}

func (p *Pilot) updateFixStatus(fix FixStatus) {
	// compute the update for fix status
	fixAlarm, fixLed := validateFixStatus(fix)

	/////////////////////////
	// Update pilot state from previous checks
	////////////////////////
	p.alarm = p.alarm || fixAlarm

	// TODO(ssoudan) wrap this stuff in something that can be tested
	if fixLed {
		p.leds[dashboard.NoGPSFix] = true
	} else {
		p.leds[dashboard.NoGPSFix] = false
	}

	/////////////////////////
	// Tell the world
	/////////////////////////
	p.tellTheWorld()
}

func (p Pilot) tellTheWorld() {
	// Keep the alarm first - so at least we get notified something is wrong
	p.alarmChan <- alarm.NewMessage(bool(p.alarm))
	p.dashboardChan <- dashboard.NewMessage(p.leds)
}

func (p *Pilot) updateFeedback(gpsHeading GPSFeedBackAction) {

	// Set the heading with the current GPS heading if it has not been set before
	if p.enabled && !p.headingSet {
		log.Info("Heading to %v", gpsHeading.Heading)
		p.heading = gpsHeading.Heading
		p.pid.Set(0) // Reference is always 0 for us
		p.headingSet = true
	}

	// TODO(ssoudan) do something with the validity

	// TODO(ssoudan) do something with the speed

	headingError := computeHeadingError(p.heading, gpsHeading.Heading)

	headingAlarm := p.checkHeadingError(headingError)

	/////////////////////////
	// Update pilot state from previous checks
	////////////////////////
	if p.enabled {
		log.Notice("Heading error is %v", headingError)

		headingControl := p.pid.Update(headingError)
		log.Notice("Heading control is %v", headingControl)

		// TODO(ssoudan) check if the control is not too large to raise the alarm or dashboard notification

		////////////////////////
		// <This section is updated even when the pilot is Enabled>
		////////////////////////
		// Update alarm state from the previously computed alarms
		p.alarm = headingAlarm // || blah

		// Update alarm state from the previously computed alarms
		if bool(headingAlarm) {
			p.leds[dashboard.HeadingErrorOutOfBounds] = true
		}

		steeringEnabled := p.computeSteeringState()
		////////////////////////
		// </This section is updated even when the pilot is Enabled>
		////////////////////////

		/////////////////////////
		// Tell the world
		/////////////////////////
		if steeringEnabled {
			log.Notice("Steering Enabled")

			// TODO(ssoudan) do something with the heading error

			// TODO(ssoudan) call the PID

			// TODO(ssoudan) check the PID output

			p.steeringChan <- steering.NewMessage(headingControl)

		} else {
			log.Notice("Steering Disabled")
		}
	} else {
		////////////////////////
		// <This section is updated even when the pilot is not enabled>
		////////////////////////
		// Alarms are UNRAISED
		p.alarm = UNRAISED

		// Update alarm state from the previously computed alarms
		p.leds[dashboard.HeadingErrorOutOfBounds] = false

		////////////////////////
		// </This section is updated even when the pilot is not enabled>
		////////////////////////
	}

	/////////////////////////
	// Tell the world
	/////////////////////////
	p.tellTheWorld()
}

// Start the event loop of the Pilot component
func (p Pilot) Start() {
	go func() {

		for {
			select {
			case m := <-p.inputChan:
				switch m := m.(type) {
				case FixStatus:
					p.updateFixStatus(m)
				case GPSFeedBackAction:
					p.updateFeedback(m)
				case EnableAction:
					p.enable()
				case DisableAction:
					p.disable()
				case error:
					log.Error("Received an error: %v", m)
				}
			}
		}

	}()

}

func (p Pilot) computeSteeringState() bool {
	return p.enabled && !bool(p.alarm)
}

type Alarm bool

const (
	RAISED   = true
	UNRAISED = false
)
