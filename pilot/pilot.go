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
* @Last Modified time: 2015-09-20 22:37:30
 */

package pilot

import (
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("pilot")

type Pilot struct {
	heading float64 // target heading (set point)
	bound   float64

	alarm      Alarm
	enabled    bool
	headingSet bool

	leds map[string]bool

	// channels with the other components
	dashboardChan chan interface{}
	inputChan     chan interface{}
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

func New(bound float64) *Pilot {
	return &Pilot{
		leds:  make(map[string]bool),
		bound: bound}
}

func (p *Pilot) SetDashboardChan(c chan interface{}) {
	p.dashboardChan = c
}

func (p *Pilot) SetInputChan(c chan interface{}) {
	p.inputChan = c
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
	p.dashboardChan <- dashboard.NewMessage(p.leds)
}

func (p *Pilot) updateFeedback(gpsHeading GPSFeedBackAction) {

	// Set the heading with the current GPS heading if it has not been set before
	if p.enabled && !p.headingSet {
		log.Info("Heading to %v", gpsHeading.Heading)
		p.heading = gpsHeading.Heading
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
	}

	// Update alarm state from the previously computed alarms
	p.alarm = Alarm(p.enabled) && headingAlarm // || blah

	// Update alarm state from the previously computed alarms
	steeringEnabled := p.computeSteeringState()

	if bool(headingAlarm) && p.enabled {
		p.leds[dashboard.HeadingErrorOutOfBounds] = true
	} else {
		p.leds[dashboard.HeadingErrorOutOfBounds] = false
	}

	/////////////////////////
	// Tell the world
	/////////////////////////
	if steeringEnabled {
		log.Notice("Steering Enabled")

		// TODO(ssoudan) do something with the heading error

		// TODO(ssoudan) call the PID

		// TODO(ssoudan) check the PID output
	} else {
		log.Notice("Steering Disabled")
	}

	// p.tellTheWorld()
	p.dashboardChan <- dashboard.NewMessage(p.leds)
}

// Start the event loop of the Pilot component
func (p Pilot) Start() chan interface{} {
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

	return p.inputChan
}

func (p Pilot) computeSteeringState() bool {
	return p.enabled && !bool(p.alarm)
}

type Alarm bool

const (
	RAISED   = true
	UNRAISED = false
)
