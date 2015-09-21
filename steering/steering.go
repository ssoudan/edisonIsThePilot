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
* @Date:   2015-09-21 17:40:00
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 19:44:22
 */

package steering

import (
	"math"
	"time"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("steering")

const (
	numberOfSteps                 = 200
	rotationSpeedInStepPerSeconds = numberOfSteps // aka 1 rotation per second
)

type Motor struct {
	actionner Actionner

	// channels
	inputChan chan interface{}
}

type Actionner interface {
	Enable() error
	Disable() error
	Move(clockwise bool, speedInStepBySeconds uint32, duration time.Duration) error
}

func New(actionner Actionner) *Motor {
	return &Motor{actionner: actionner}
}

type message struct {
	rotationInDegree float64
}

func NewMessage(rotationInDegree float64) interface{} {
	return message{rotationInDegree: rotationInDegree}
}

func (d *Motor) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

func rotationInDegreeToMove(rotationInDegree float64) (clockwise bool, speed uint32, duration time.Duration) {
	clockwise = rotationInDegree > 0.
	speed = uint32(rotationSpeedInStepPerSeconds)
	duration = time.Duration(
		math.Abs(rotationInDegree/360.*numberOfSteps/rotationSpeedInStepPerSeconds) * float64(time.Second))

	return
}

func (d *Motor) processMotorState(m message) {

	rotationInDegree := m.rotationInDegree

	if rotationInDegree != 0. {
		d.actionner.Enable()
		defer d.actionner.Disable()

		clockwise, speed, duration := rotationInDegreeToMove(rotationInDegree)

		err := d.actionner.Move(clockwise, speed, duration)
		if err != nil {
			log.Error("Failed to move [clockwise=%v] for %v at %v", clockwise, duration, speed)
		}

	}
}

func (d Motor) processMessage(m message) {

	// move
	d.processMotorState(m)
}

// Shutdown sets all the state to down and notify the handlers
func (d Motor) Shutdown() {

	// disable the steering -- should not be Enabled()
	d.actionner.Disable()
}

// Start the event loop of the Motor component
func (d Motor) Start() {

	go func() {
		for true {
			select {
			case m := <-d.inputChan:
				switch m := m.(type) {
				case message:
					d.processMessage(m)
				}

			}

		}
	}()

}
