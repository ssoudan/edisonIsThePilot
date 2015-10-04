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
* @Last Modified time: 2015-10-04 00:00:47
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
	rotationSpeedInStepPerSeconds = 2 * numberOfSteps // aka 2 rotation per second
)

type Motor struct {
	actionner Actionner

	// channels
	inputChan    chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

type Actionner interface {
	Enable() error
	Disable() error
	Move(clockwise bool, speedInStepBySeconds uint32, duration time.Duration) error
}

func New(actionner Actionner) *Motor {
	return &Motor{actionner: actionner, shutdownChan: make(chan interface{})}
}

type message struct {
	rotationInDegree float64
	stayEnabled      bool
}

func NewMessage(rotationInDegree float64, stayEnabled bool) interface{} {
	return message{rotationInDegree: rotationInDegree, stayEnabled: stayEnabled}
}

func (m *Motor) SetInputChan(c chan interface{}) {
	m.inputChan = c
}

func (m *Motor) SetPanicChan(c chan interface{}) {
	m.panicChan = c
}

func rotationInDegreeToMove(rotationInDegree float64) (clockwise bool, speed uint32, duration time.Duration) {
	clockwise = rotationInDegree > 0.
	speed = uint32(rotationSpeedInStepPerSeconds)
	duration = time.Duration(
		math.Abs(rotationInDegree/360.*numberOfSteps/rotationSpeedInStepPerSeconds) * float64(time.Second))

	return
}

func (m *Motor) processMotorState(msg message) {

	rotationInDegree := msg.rotationInDegree

	if !msg.stayEnabled {
		defer m.actionner.Disable()
	}

	if rotationInDegree != 0. {
		m.actionner.Enable()
		clockwise, speed, duration := rotationInDegreeToMove(rotationInDegree)

		err := m.actionner.Move(clockwise, speed, duration)
		if err != nil {
			log.Panicf("Failed to move [clockwise=%v] for %v at %v", clockwise, duration, speed)
		}

	}
}

func (m Motor) processMessage(msg message) {
	// no state to update

	// move
	m.processMotorState(msg)
}

// Shutdown sets all the state to down and notify the handlers
func (m Motor) Shutdown() {
	m.shutdownChan <- 1
	<-m.shutdownChan
}

func (m Motor) shutdown() {
	// disable the steering -- should not be Enabled()
	m.actionner.Disable()
	close(m.shutdownChan)
}

// Start the event loop of the Motor component
func (m Motor) Start() {

	go func() {

		defer func() {
			if r := recover(); r != nil {
				m.panicChan <- r
			}
		}()

		for {
			select {
			case msg := <-m.inputChan:
				switch msg := msg.(type) {
				case message:
					m.processMessage(msg)
				}
			case <-m.shutdownChan:
				m.shutdown()
				return
			}

		}
	}()

}
