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
* @Date:   2015-09-26 17:50:09
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-30 12:45:51
 */

package main

import (
	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("ledControl")

func main() {

	mapMessageToGPIO := func(message string, pin byte) gpio.Gpio {

		// kill the process (via log.Panic -> recover -> panicChan -> go routine -> log.Fatal) in case we can't create the GPIO
		err := gpio.EnableGPIO(pin)
		if err != nil {
			log.Panic(err)
		}

		var g = gpio.New(pin)
		if !g.IsExported() {
			err = g.Export()
			if err != nil {
				log.Panic(err)
			}
		}

		err = g.SetDirection(gpio.OUT)
		if err != nil {
			log.Panic(err)
		}

		// Test Disabled and Enabled state for each LEDs
		err = g.Disable()
		if err != nil {
			log.Panic(err)
		}

		return g
	}
	dashboardGPIOs := make(map[string]gpio.Gpio, len(conf.MessageToPin))
	for _, v := range conf.MessageToPin {
		g := mapMessageToGPIO(v.Message, v.Pin)
		dashboardGPIOs[v.Message] = g
	}

	for _, g := range dashboardGPIOs {
		g.Enable()
	}

}
