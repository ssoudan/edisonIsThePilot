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
* @Last Modified time: 2015-10-03 23:08:40
 */

package pilot

import (
	"time"
	"github.com/ssoudan/edisonIsThePilot/conf"
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
	shutdownChan  chan interface{}
	panicChan     chan interface{}
}

type Controller interface {
	Set(sp float64)
	Update(value float64) float64
	OutputLimits() (float64, float64)
}

type Leds map[Led]bool
type Led string

func (p *Pilot) checkHeadingError(headingError float64) Alarm {

	inputStatus := validateInput(p.bound, headingError)
	return computeAlarmStateForInputStatus(p.alarm, inputStatus)

}

func ComputeHeadingError(heading float64, gpsHeading float64) float64 {

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
		leds:         make(map[string]bool),
		bound:        bound,
		pid:          controller,
		shutdownChan: make(chan interface{})}
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

func (p *Pilot) SetPanicChan(c chan interface{}) {
	p.panicChan = c
}

func (p *Pilot) updateFixStatus(fix FixStatus) {
	// compute the update for fix status
	fixAlarm, fixLed := validateFixStatus(fix)

	/////////////////////////
	// Update pilot state from previous checks
	////////////////////////
	p.alarm = p.alarm || fixAlarm

	p.leds[dashboard.NoGPSFix] = fixLed
}

func (p Pilot) tellTheWorld() {
	// Keep the alarm first - so at least we get notified something is wrong
	p.alarmChan <- alarm.NewMessage(bool(p.alarm))
	p.dashboardChan <- dashboard.NewMessage(p.leds)
}

func (p *Pilot) updateFeedback(gpsHeading GPSFeedBackAction) {

	// Set the heading with the current GPS heading if it has not been set before
	if p.enabled && !p.headingSet && gpsHeading.Validity {
		log.Info("Heading to %v", gpsHeading.Heading)
		p.heading = gpsHeading.Heading
		p.pid.Set(0) // Reference is always 0 for us
		p.headingSet = true
	}

	// check the validity of the message validity of the gps message
	validityAlarm := checkValidityError(gpsHeading.Validity)

	// check the speed
	speedAlarm := checkSpeedError(gpsHeading.Speed)

	headingError := ComputeHeadingError(p.heading, gpsHeading.Heading)

	headingAlarm := !validityAlarm && !speedAlarm && p.checkHeadingError(headingError)

	/////////////////////////
	// Update pilot state from previous checks
	////////////////////////
	if p.enabled {

		// Update alarm state from the previously computed alarms
		p.alarm = p.alarm || headingAlarm || validityAlarm || speedAlarm

		if p.alarm == UNRAISED {
			log.Notice("Heading error is %v", headingError)
		}

		// Update alarm state from the previously computed alarms
		if bool(headingAlarm) {
			p.leds[dashboard.HeadingErrorOutOfBounds] = true
		}
		if bool(validityAlarm) {
			p.leds[dashboard.InvalidGPSData] = true
		}
		if bool(speedAlarm) {
			p.leds[dashboard.SpeedTooLow] = true
		}

		headingControl := p.pid.Update(headingError)

		steeringEnabled := p.computeSteeringState()

		if steeringEnabled {
			log.Notice("Heading control is %v", headingControl)

			log.Notice("Steering Enabled")

			// check the PID output
			minPIDoutput, maxPIDOutput := p.pid.OutputLimits()
			if headingControl <= minPIDoutput || headingControl >= maxPIDOutput {
				p.leds[dashboard.CorrectionAtLimit] = true
			}

			p.steeringChan <- steering.NewMessage(headingControl, true)

		} else {
			log.Notice("Steering Disabled")
			p.steeringChan <- steering.NewMessage(0, false)
		}
	} else {
		////////////////////////
		// <This section is updated when the pilot is not enabled>
		////////////////////////
		// Alarms are UNRAISED
		p.alarm = UNRAISED

		// Update alarm state from the previously computed alarms
		p.leds[dashboard.HeadingErrorOutOfBounds] = false // Doesn't make sense when disabled
		p.leds[dashboard.InvalidGPSData] = bool(validityAlarm)
		p.leds[dashboard.SpeedTooLow] = bool(speedAlarm)
		p.leds[dashboard.CorrectionAtLimit] = false // Doesn't make sense when disabled

		// make sure the steering is disabled
		p.steeringChan <- steering.NewMessage(0, false)
		////////////////////////
		// </This section is updated when the pilot is not enabled>
		////////////////////////
	}
}

func (p *Pilot) updateAfterTimeout() {
	if p.enabled {
		p.alarm = RAISED
		p.leds[dashboard.NoGPSFix] = true
		// make sure the steering is disabled
		p.steeringChan <- steering.NewMessage(0, false)
	}
}

func (p *Pilot) updateAfterError() {
	if p.enabled {
		p.alarm = RAISED
		// make sure the steering is disabled
		p.steeringChan <- steering.NewMessage(0, false)
	}
}

// Start the event loop of the Pilot component
func (p Pilot) Start() {

	go func() {

		defer func() {
			if r := recover(); r != nil {
				p.panicChan <- r
			}
		}()

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
					p.updateAfterError()
				}
			case <-time.After(time.Duration(conf.Conf.NoInputMessageTimeoutInSeconds) * time.Second):
				p.updateAfterTimeout()
			case <-p.shutdownChan:
				p.shutdown()
				/////////////////////////
				// Tell the world
				/////////////////////////
				p.tellTheWorld()
				return
			}

			/////////////////////////
			// Tell the world
			/////////////////////////
			p.tellTheWorld()
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
