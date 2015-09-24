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
* @Last Modified time: 2015-09-24 13:12:42
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
	inputChan    chan interface{}
	shutdownChan chan interface{}
}

type Enablable interface {
	Enable() error
	Disable() error
}

func New(alarmHandler Enablable) *Alarm {
	return &Alarm{alarmHandler: alarmHandler, shutdownChan: make(chan interface{})}
}

type message struct {
	alarm bool
}

func NewMessage(alarm bool) interface{} {
	return message{alarm: alarm}
}

func (d *Alarm) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

func (d Alarm) processAlarmState() {

	alarm := d.alarmHandler
	state := d.alarmState

	if state {
		err := alarm.Enable()
		if err != nil {
			log.Error("Failed to change alarm state for [%s = %v]: %v", state, err)
		}
	} else {
		err := alarm.Disable()
		if err != nil {
			log.Error("Failed to change alarm state for [%s = %v]: %v", state, err)
		}
	}

}

func (d *Alarm) processMessage(m message) {
	// Update the state
	d.alarmState = m.alarm

	// Update the alarm
	d.processAlarmState()
}

// Shutdown sets all the state to down and notify the handlers
func (d Alarm) Shutdown() {

	d.shutdownChan <- 1
	<-d.shutdownChan
}

func (d *Alarm) shutdown() {
	d.alarmState = false

	// Update the alarm
	d.processAlarmState()

	close(d.shutdownChan)
}

// Start the event loop of the Alarm component
func (d *Alarm) Start() {

	go func() {
		for {
			select {
			case m := <-d.inputChan:
				switch m := m.(type) {
				case message:
					d.processMessage(m)
				}
			case <-d.shutdownChan:
				d.shutdown()

				return
			}

		}
	}()

}
