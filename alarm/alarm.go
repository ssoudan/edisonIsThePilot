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
* @Last Modified time: 2015-10-21 14:21:19
 */

package alarm

import (
	"sync"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
)

var log = logger.Log("alarm")

// Alarm is the component that manage the external alarm
type Alarm struct {
	alarmHandler types.Enablable

	mu         sync.RWMutex
	alarmState bool // protected by mu

	// channels
	inputChan    chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

// New creates a new Alarm component for an types.Enablable
func New(alarmHandler types.Enablable) *Alarm {
	return &Alarm{alarmHandler: alarmHandler, shutdownChan: make(chan interface{})}
}

type message struct {
	alarm bool
}

// NewMessage creates a new alarm state update message
func NewMessage(alarm bool) interface{} {
	return message{alarm: alarm}
}

// SetInputChan sets the channel where the alarm updates are sent
func (d *Alarm) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

// SetPanicChan sets the channel where this component will send the panics
func (d *Alarm) SetPanicChan(c chan interface{}) {
	d.panicChan = c
}

// Enabled is true when the alarm is raised, false otherwise
func (d *Alarm) Enabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.alarmState
}

func (d *Alarm) processAlarmState() {

	// NOTE: already hold the WLock

	alarm := d.alarmHandler
	state := d.alarmState

	if state {
		err := alarm.Enable()
		if err != nil {
			log.Panicf("Failed to change alarm state to %v: %v", state, err)
		}
	} else {
		err := alarm.Disable()
		if err != nil {
			log.Panicf("Failed to change alarm state to %v: %v", state, err)
		}
	}

}

func (d *Alarm) processMessage(m message) {
	// Update the state
	d.mu.Lock()
	defer d.mu.Unlock()
	d.alarmState = m.alarm

	// Update the alarm
	d.processAlarmState()
}

// Shutdown sets all the state to down and notify the handlers
func (d *Alarm) Shutdown() {

	d.shutdownChan <- 1
	<-d.shutdownChan
}

func (d *Alarm) shutdown() {

	// Update the alarm
	d.processAlarmState()

	close(d.shutdownChan)
}

// Start the event loop of the Alarm component
func (d *Alarm) Start() {

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
