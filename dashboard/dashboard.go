/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 16:30:19
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:00:53
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
