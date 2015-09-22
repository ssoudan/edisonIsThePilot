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
* @Date:   2015-09-21 18:58:22
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-22 13:09:42
 */

package motor

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/drivers/pwm"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("motor")

type Motor struct {
	dirGPIO gpio.Gpio

	sleepGPIO gpio.Gpio
	stepPwm   *pwm.Pwm
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func New(stepPin, stepPwmId, dirPin, sleepPin byte) *Motor {

	// Create the dir GPIO
	err := gpio.EnableGPIO(dirPin)
	check(err)

	dirGPIO := gpio.New(dirPin)
	if !dirGPIO.IsExported() {
		err = dirGPIO.Export()
		check(err)
	}

	err = dirGPIO.SetDirection(gpio.OUT)
	check(err)

	// Test Disabled and Enabled state for each LEDs
	err = dirGPIO.Disable()
	check(err)

	// Create the sleep GPIO
	err = gpio.EnableGPIO(sleepPin)
	check(err)

	sleepGPIO := gpio.New(sleepPin)
	if !sleepGPIO.IsExported() {
		err = sleepGPIO.Export()
		check(err)
	}

	err = sleepGPIO.SetDirection(gpio.OUT)
	check(err)

	// Test Disabled and Enabled state for each LEDs
	err = sleepGPIO.Enable()
	check(err)

	// Create the Step pwm
	stepPwm, err := pwm.New(stepPwmId, stepPin)
	check(err)
	if !stepPwm.IsExported() {
		err = stepPwm.Export()
		check(err)
	}

	err = stepPwm.Disable()
	check(err)

	return &Motor{dirGPIO: dirGPIO, sleepGPIO: sleepGPIO, stepPwm: stepPwm}
}

func (m Motor) Enable() error {
	return m.sleepGPIO.Disable()
}

func (m Motor) Disable() error {
	return m.sleepGPIO.Enable()
}

func (m Motor) Move(clockwise bool, stepsBySecond uint32, duration time.Duration) error {
	var err error
	if clockwise {
		err = m.dirGPIO.Enable()
		if err != nil {
			log.Error("Failed to set direction: %v", err)
			return err
		}
	} else {
		err = m.dirGPIO.Disable()
		if err != nil {
			log.Error("Failed to set direction: %v", err)
			return err
		}
	}

	period := time.Duration(1. / float64(stepsBySecond) * float64(time.Second))
	if period < pwm.MinPeriod {
		originalPeriod := period
		period = pwm.MinPeriod
		log.Warning("period out of bounds: changed from %d to %d", originalPeriod, period)
	}
	if period > pwm.MaxPeriod {
		originalPeriod := period
		period = pwm.MaxPeriod
		log.Warning("period out of bounds: changed from %d to %d", originalPeriod, period)
	}

	err = m.stepPwm.SetPeriodAndDutyCycle(period, 0.5)
	if err != nil {
		return err
	}

	err = m.Enable()
	if err != nil {
		return err
	}
	log.Debug("rotation [clockwise=%v,period=%v,duration=%v]", clockwise, period, duration)
	time.Sleep(duration)
	err = m.Disable()
	if err != nil {
		return err
	}

	log.Debug("rotation stopped")
	return nil

}
