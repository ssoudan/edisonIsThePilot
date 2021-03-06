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
* @Last Modified time: 2015-10-21 13:09:34
 */

package motor

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/drivers/pwm"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("motor")

// Motor is a driver for a stepper motor
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

// New creates a new Motor for a stepper motor drived with GPIOs
func New(stepPin, stepPwmID, dirPin, sleepPin byte) *Motor {

	// Create the dir GPIO
	err := gpio.EnableGPIO(dirPin)
	check(err)

	dirGPIO := gpio.New(dirPin)
	if !dirGPIO.IsExported() {
		err = dirGPIO.Export()
		check(err)
	}

	err = dirGPIO.SetDirection(gpio.OutDirection)
	check(err)

	// Test Disabled and Enabled state for each pin
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

	err = sleepGPIO.SetDirection(gpio.OutDirection)
	check(err)

	// Test Disabled and Enabled state for each pin
	err = sleepGPIO.Enable()
	check(err)

	// Create the Step pwm
	stepPwm, err := pwm.New(stepPwmID, stepPin)
	check(err)
	if !stepPwm.IsExported() {
		err = stepPwm.Export()
		check(err)
	}

	err = stepPwm.Disable()
	check(err)

	return &Motor{dirGPIO: dirGPIO, sleepGPIO: sleepGPIO, stepPwm: stepPwm}
}

// Enable enables the torque
func (m Motor) Enable() error {
	return m.sleepGPIO.Disable()
}

// Disable disables the torque
func (m Motor) Disable() error {
	return m.sleepGPIO.Enable()
}

// Move makes the motor rotate in the given direction at the specified speed for a given duration -- make sure to Enable() the motor first
func (m Motor) Move(clockwise bool, stepsBySecond uint32, duration time.Duration) error {
	if stepsBySecond == 0 || duration == 0 {
		return nil
	}

	var err error
	if clockwise {
		err = m.dirGPIO.Disable()
		if err != nil {
			log.Panicf("Failed to set direction: %v", err)
			return err // Not supposed to reach here
		}
	} else {
		err = m.dirGPIO.Enable()
		if err != nil {
			log.Panicf("Failed to set direction: %v", err)
			return err // Not supposed to reach here
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

	err = m.stepPwm.Enable()
	if err != nil {
		return err
	}

	// err = m.Enable()
	// if err != nil {
	// 	return err
	// }
	time.Sleep(duration)
	// err = m.Disable()
	// if err != nil {
	// 	return err
	// }
	err = m.stepPwm.Disable()
	if err != nil {
		return err
	}

	return nil

}

// Unexport unexports the GPIO used by to drive the motor
func (m Motor) Unexport() {
	m.dirGPIO.Unexport()
	m.sleepGPIO.Unexport()
	m.stepPwm.Unexport()
}
