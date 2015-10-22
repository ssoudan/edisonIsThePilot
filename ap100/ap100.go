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
* @Date:   2015-10-13 17:12:30
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-19 15:34:13
 */

package ap100

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("ap100")

// AP100 is the component that send the course to an external device (via a sin/cos interface)
type AP100 struct {
	output Compass

	// channels
	inputChan    chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

// Compass is something on which we can update the course (in degree)
type Compass interface {
	UpdateCourse(course uint16) error
}

// New creates a new AP100 component for a Compass
func New(output Compass) *AP100 {
	return &AP100{output: output, shutdownChan: make(chan interface{})}
}

type message struct {
	course uint16
}

// NewMessage creates course update messages
func NewMessage(course uint16) interface{} {
	return message{course: course}
}

// SetInputChan sets the channel where this component will be listening to for the course updates
func (d *AP100) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

// SetPanicChan sets the channel where this component will send the panics
func (d *AP100) SetPanicChan(c chan interface{}) {
	d.panicChan = c
}

func (d *AP100) processMessage(m message) {
	d.output.UpdateCourse(m.course)
}

// Shutdown sets all the state to down and notify the handlers
func (d *AP100) Shutdown() {

	d.shutdownChan <- 1
	<-d.shutdownChan
}

func (d *AP100) shutdown() {

	close(d.shutdownChan)
}

// Start the event loop of the AP100 component
func (d *AP100) Start() {

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
