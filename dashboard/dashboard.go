/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 16:30:19
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 16:51:20
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
}

type Message struct {
	Leds map[string]bool
}

func NewMessage(leds map[string]bool) Message {
	return Message{Leds: leds}
}

func (d Dashboard) Run() chan Message {
	c := make(chan Message)

	d.leds = make(map[string]bool)

	go func() {
		for true {
			m := <-c

			for k, v := range m.Leds {
				d.leds[k] = v
			}

			// TODO(ssoudan) update the LEDs
			for k, v := range d.leds {
				log.Info(" %s = %v", k, v)
			}

		}
	}()

	return c
}
