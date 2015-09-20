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
* @Last Modified time: 2015-09-20 22:17:27
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
	NoLed                   = ""
)

var log = logger.Log("dashboard")

type Dashboard struct {
	leds map[string]bool

	// channels
	inputChan chan interface{}
}

func New() *Dashboard {
	return &Dashboard{leds: make(map[string]bool)}
}

type Message struct {
	Leds map[string]bool
}

func NewMessage(leds map[string]bool) Message {
	return Message{Leds: leds}
}

func (d *Dashboard) SetInputChan(c chan interface{}) {
	d.inputChan = c
}

// Start the event loop of the Dashboard component
func (d Dashboard) Start() {

	go func() {
		for true {
			select {
			case m := <-d.inputChan:
				switch m := m.(type) {
				case Message:
					for k, v := range m.Leds {
						d.leds[k] = v
					}

					for k, v := range d.leds {
						// TODO(ssoudan) update the LEDs
						log.Info(" %s = %v", k, v)
					}
				}

			}

		}
	}()

}
