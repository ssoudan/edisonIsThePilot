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
* @Date:   2015-09-20 16:30:19
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-22 13:14:47
 */

package dashboard

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

const (
	NoGPSFix                = "NoGPSFix"
	InvalidGPSData          = "InvalidGPSData"
	SpeedTooLow             = "SpeedTooLow"
	HeadingErrorOutOfBounds = "HeadingErrorOutOfBounds"
	CorrectionAtLimit       = "CorrectionAtLimit"
)

var log = logger.Log("dashboard")

type Dashboard struct {
	leds       map[string]bool
	ledHandler map[string]Enablable

	// channels
	inputChan chan interface{}
}

type Enablable interface {
	Enable() error
	Disable() error
}

func New() *Dashboard {
	return &Dashboard{
		leds:       make(map[string]bool),
		ledHandler: make(map[string]Enablable),
	}
}

func (d *Dashboard) RegisterMessageHandler(message string, handler Enablable) {
	d.ledHandler[message] = handler
}

type message struct {
	Leds map[string]bool
}

func NewMessage(leds map[string]bool) interface{} {
	return message{Leds: leds}
}

func (d *Dashboard) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

func (d *Dashboard) processLedState() {
	for led, state := range d.leds {
		gpio, ok := d.ledHandler[led]
		if ok {
			if state {
				err := gpio.Enable()
				if err != nil {
					log.Error("Failed to change led state for [%s = %v]: %v", led, state, err)
					// TODO(ssoudan) probably a case for an alarm
				}
			} else {
				err := gpio.Disable()
				if err != nil {
					log.Error("Failed to change led state for [%s = %v]: %v", led, state, err)
					// TODO(ssoudan) probably a case for an alarm
				}
			}
		} else {
			log.Warning("No LED for [%s = %v]", led, state)
		}
	}

}

func (d Dashboard) processMessage(m message) {
	// Update the state
	for k, v := range m.Leds {
		d.leds[k] = v
	}

	// Update the LEDs
	d.processLedState()
}

// Shutdown sets all the state to down and notify the handlers
func (d Dashboard) Shutdown() {
	for k, _ := range d.leds {
		d.leds[k] = false
	}

	// Update the LEDs
	d.processLedState()
}

// Start the event loop of the Dashboard component
func (d Dashboard) Start() {

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
