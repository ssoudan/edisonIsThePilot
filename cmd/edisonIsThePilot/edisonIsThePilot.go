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
* @Date:   2015-09-18 12:20:59
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 22:16:58
 */

package main

import (
	"github.com/felixge/pidctrl"

	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/ssoudan/edisonIsThePilot/compass/hmc"
	"github.com/ssoudan/edisonIsThePilot/alarm"
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/drivers/pwm"
	"github.com/ssoudan/edisonIsThePilot/gps"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	"github.com/ssoudan/edisonIsThePilot/steering"
)

var log = logger.Log("edisonIsThePilot")

var messageToPin = map[string]byte{
	dashboard.NoGPSFix:                40, // J19 - pin 10
	dashboard.InvalidGPSData:          43, // J19 - pin 11
	dashboard.SpeedTooLow:             48, // J19 - pin 6
	dashboard.HeadingErrorOutOfBounds: 82, // J19 - pin 13
	dashboard.CorrectionAtLimit:       83, // J19 - pin 14
}

const (
	alarmGpioPin  = 183 // J18 - pin 8
	alarmGpioPWM  = 3
	motorDirPin   = 165 // J18 - pin 2
	motorSleepPin = 12  // J18 - pin 7
	motorStepPin  = 182 // J17 - pin 1
	motorStepPwm  = 2
)

const (
	maxPIDOutputLimits = 15
	minPIDOutputLimits = -15
	p                  = 1
	i                  = 0.1
	d                  = 0.1
)

func main() {

	////////////////////////////////////////
	// HMC5883 stuffs
	////////////////////////////////////////
	// compass := hmc.New(6)
	// for !compass.Begin() {

	// }

	// // Set measurement range
	// compass.SetRange(hmc.HMC5883L_RANGE_1_3GA)

	// // Set measurement mode
	// compass.SetMeasurementMode(hmc.HMC5883L_CONTINOUS)

	// // Set data rate
	// compass.SetDataRate(hmc.HMC5883L_DATARATE_3HZ)

	// // Set number of samples averaged
	// compass.SetSamples(hmc.HMC5883L_SAMPLES_8)

	// // Set calibration offset. See HMC5883L_calibration.ino
	// compass.SetOffset(-82, 72)

	// mag, err := compass.ReadNormalize()
	// if err == nil {
	// 	log.Info("Compass reading is %v", mag)
	// }

	////////////////////////////////////////
	// I2C stuffs
	////////////////////////////////////////

	// bp, err := i2c.Bus(1)
	// if err != nil {
	// 	log.Panicf("failed to create bus: %v\n", err)
	// }
	// addr := 0x12
	// reg := 0x13
	// len := 0x1
	// data, err := bp.ReadByteBlock(addr, reg, length)

	////////////////////////////////////////
	// a beautiful dashboard
	////////////////////////////////////////
	dashboard := dashboard.New()
	dashboardChan := make(chan interface{})
	dashboard.SetInputChan(dashboardChan)
	mapMessageToGPIO := func(message string, pin byte) gpio.Gpio {

		// kill the process (via log.Fatal) in case we can't create the GPIO
		err := gpio.EnableGPIO(pin)
		if err != nil {
			log.Fatal(err)
		}

		var g = gpio.New(pin)
		if !g.IsExported() {
			err = g.Export()
			if err != nil {
				log.Fatal(err)
			}
		}

		err = g.SetDirection(gpio.OUT)
		if err != nil {
			log.Fatal(err)
		}

		// Test Disabled and Enabled state for each LEDs
		err = g.Disable()
		if err != nil {
			log.Fatal(err)
		}

		log.Info("[AUTOTEST] %s LED is ON", message)
		time.Sleep(1 * time.Second)

		err = g.Enable()
		if err != nil {
			log.Fatal(err)
		}

		log.Info("[AUTOTEST] %s LED is OFF", message)
		time.Sleep(1 * time.Second)

		err = g.Disable()
		if err != nil {
			log.Fatal(err)
		}

		return g
	}
	for k, v := range messageToPin {
		dashboard.RegisterMessageHandler(k, mapMessageToGPIO(k, v))
	}

	////////////////////////////////////////
	// a nice alarm
	////////////////////////////////////////
	alarmPwm := func(pin byte, pwmId byte) pwm.Pwm {

		// kill the process (via log.Fatal) in case we can't create the PWM
		err := gpio.EnablePWM(pin)
		if err != nil {
			log.Fatal(err)
		}

		var pwm = pwm.New(alarmGpioPWM)
		if !pwm.IsExported() {
			err = pwm.Export()
			if err != nil {
				log.Fatal(err)
			}
		}

		pwm.Disable()

		if err = pwm.SetPeriodAndDutyCycle(200*time.Millisecond, 0.5); err != nil {
			log.Fatal(err)
		}

		if err = pwm.Enable(); err != nil {
			log.Fatal(err)
		}
		log.Info("[AUTOTEST] alarm is ON")

		time.Sleep(2 * time.Second)
		if err = pwm.Disable(); err != nil {
			log.Fatal(err)
		}
		log.Info("[AUTOTEST] alarm is OFF")

		return pwm
	}(alarmGpioPin, alarmGpioPWM)

	alarm := alarm.New(alarmPwm)
	alarmChan := make(chan interface{})
	alarm.SetInputChan(alarmChan)

	////////////////////////////////////////
	// an astonishing steering
	////////////////////////////////////////
	motor := motor.New(motorStepPin, motorStepPwm, motorDirPin, motorSleepPin)

	steering := steering.New(motor)
	steeringChan := make(chan interface{})
	steering.SetInputChan(steeringChan)

	////////////////////////////////////////
	// PID stuffs
	////////////////////////////////////////
	pidController := pidctrl.NewPIDController(p, i, d)
	pidController.SetOutputLimits(minPIDOutputLimits, maxPIDOutputLimits)

	////////////////////////////////////////
	// pilot stuffs
	////////////////////////////////////////
	thePilot := pilot.New(pidController, 15.)
	pilotChan := make(chan interface{})
	thePilot.SetInputChan(pilotChan)
	thePilot.SetDashboardChan(dashboardChan)
	thePilot.SetAlarmChan(alarmChan)
	thePilot.SetSteeringChan(steeringChan)

	////////////////////////////////////////
	// gps stuffs
	////////////////////////////////////////
	gps := gps.New("/dev/ttyMFD1")
	gps.SetMessagesChan(pilotChan)
	gps.SetErrorChan(pilotChan)

	gps.Start()
	dashboard.Start()
	alarm.Start()
	thePilot.Start()
	steering.Start()

	// For tests
	go func() {
		for {
			log.Notice("Disabling the pilot")
			thePilot.Disable()
			time.Sleep(35 * time.Second)
			log.Notice("Enabling the pilot")
			thePilot.Enable()
			time.Sleep(35 * time.Second)
		}
	}()

	// Wait until we receive a signal
	waitForInterrupt()

	dashboard.Shutdown()
	alarm.Shutdown()
	steering.Shutdown()
}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigChan:
		log.Info("Interrupted - exiting")
	}
}

// TODO(ssoudan) set alarm if exited under an error condition
