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
* @Date:   2015-09-21 15:42:21
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 17:24:28
 */

package alarm

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("alarm")

type Alarm struct {
	alarmHandler Enablable
	alarmState   bool

	// channels
	inputChan chan interface{}
}

type Enablable interface {
	Enable() error
	Disable() error
}

func New(alarmHandler Enablable) *Alarm {
	return &Alarm{alarmHandler: alarmHandler}
}

type message struct {
	alarm bool
}

func (m message) IsAlarmRaised() bool {
	return m.alarm
}

func NewMessage(alarm bool) interface{} {
	return message{alarm: alarm}
}

func (d *Alarm) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

func (d *Alarm) processAlarmState() {

	gpio := d.alarmHandler
	state := d.alarmState

	if state {
		err := gpio.Enable()
		if err != nil {
			log.Error("Failed to change alarm state for [%s = %v]: %v", state, err)
		}
	} else {
		err := gpio.Disable()
		if err != nil {
			log.Error("Failed to change alarm state for [%s = %v]: %v", state, err)
		}
	}

}

func (d Alarm) processMessage(m message) {
	// Update the state
	d.alarmState = m.alarm

	// Update the LEDs
	d.processAlarmState()
}

// Shutdown sets all the state to down and notify the handlers
func (d Alarm) Shutdown() {

	d.alarmState = false

	// Update the LEDs
	d.processAlarmState()
}

// Start the event loop of the Alarm component
func (d Alarm) Start() {

	go func() {
		for true {
			select {
			case m := <-d.inputChan:
				switch m := m.(type) {
				case message:
					log.Info("Got an alarm message %v", m)
					d.processMessage(m)
					log.Info("Processed an alarm message %v", m)
				}

			}

		}
	}()

}
