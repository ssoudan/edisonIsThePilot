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
* @Last Modified time: 2015-10-21 13:13:43
 */

package dashboard

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
)

// Warning/Error leds
const (
	NoGPSFix                = "NoGPSFix"
	InvalidGPSData          = "InvalidGPSData"
	SpeedTooLow             = "SpeedTooLow"
	HeadingErrorOutOfBounds = "HeadingErrorOutOfBounds"
	CorrectionAtLimit       = "CorrectionAtLimit"
)

var log = logger.Log("dashboard")

// Dashboard is the component that receives the alarms and warnings and enable or disable the corresponding LEDs
type Dashboard struct {
	leds       map[string]bool
	ledHandler map[string]types.Enablable

	// channels
	inputChan    chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

// New creates a new Dashboard component
func New() *Dashboard {
	return &Dashboard{
		leds:         make(map[string]bool),
		ledHandler:   make(map[string]types.Enablable),
		shutdownChan: make(chan interface{}),
	}
}

// RegisterMessageHandler registers the handler for a particular message - there can only be one per message
func (d *Dashboard) RegisterMessageHandler(message string, handler types.Enablable) {
	d.ledHandler[message] = handler
}

type message struct {
	Leds map[string]bool
}

// NewMessage creates a LED update message
func NewMessage(leds map[string]bool) interface{} {
	return message{Leds: leds}
}

// SetInputChan sets the channel where the LED update messages are sent
func (d *Dashboard) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

// SetPanicChan sets the channel where the panics are sent
func (d *Dashboard) SetPanicChan(c chan interface{}) {
	d.panicChan = c
}

func (d Dashboard) processLedState() {
	for led, state := range d.leds {
		gpio, ok := d.ledHandler[led]
		if ok {
			if state {
				err := gpio.Enable()
				if err != nil {
					log.Panicf("Failed to change led state for [%s = %v]: %v", led, state, err)
				}
			} else {
				err := gpio.Disable()
				if err != nil {
					log.Panicf("Failed to change led state for [%s = %v]: %v", led, state, err)
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

func (d Dashboard) processQueryActionMessage(m queryMessage) {
	m.backChannel <- d.leds
}

// Shutdown sets all the state to down and notify the handlers
func (d Dashboard) Shutdown() {
	d.shutdownChan <- 1
	<-d.shutdownChan
}

func (d Dashboard) shutdown() {
	for k := range d.leds {
		d.leds[k] = false
	}

	// Update the LEDs
	d.processLedState()

	close(d.shutdownChan)
}

type queryMessage struct {
	backChannel chan (map[string]bool)
}

// GetDashboardInfoAction returns the current dashboard info
func (d Dashboard) GetDashboardInfoAction() map[string]bool {
	c := make(chan map[string]bool)
	defer close(c)
	d.inputChan <- queryMessage{backChannel: c}
	return <-c
}

// Start the event loop of the Dashboard component
func (d Dashboard) Start() {

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
				case queryMessage:
					d.processQueryActionMessage(m)
				}
			case <-d.shutdownChan:
				d.shutdown()
				return
			}

		}
	}()

}
