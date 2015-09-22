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
* @Last Modified time: 2015-09-22 20:49:49
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
	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/control"
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
	for k, v := range conf.MessageToPin {
		dashboard.RegisterMessageHandler(k, mapMessageToGPIO(k, v))
	}

	////////////////////////////////////////
	// a nice alarm
	////////////////////////////////////////
	alarmPwm := func(pin byte, pwmId byte) *pwm.Pwm {

		// kill the process (via log.Fatal) in case we can't create the PWM
		pwm, err := pwm.New(pwmId, pin)
		if err != nil {
			log.Fatal(err)
		}

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
	}(conf.AlarmGpioPin, conf.AlarmGpioPWM)

	alarm := alarm.New(alarmPwm)
	alarmChan := make(chan interface{})
	alarm.SetInputChan(alarmChan)

	////////////////////////////////////////
	// an astonishing steering
	////////////////////////////////////////
	motor := motor.New(
		conf.MotorStepPin,
		conf.MotorStepPwm,
		conf.MotorDirPin,
		conf.MotorSleepPin)

	steering := steering.New(motor)
	steeringChan := make(chan interface{})
	steering.SetInputChan(steeringChan)

	////////////////////////////////////////
	// PID stuffs
	////////////////////////////////////////
	pidController := pidctrl.NewPIDController(conf.P, conf.I, conf.D)
	pidController.SetOutputLimits(conf.MinPIDOutputLimits, conf.MaxPIDOutputLimits)

	////////////////////////////////////////
	// pilot stuffs
	////////////////////////////////////////
	thePilot := pilot.New(pidController, conf.Bounds)
	pilotChan := make(chan interface{})
	thePilot.SetInputChan(pilotChan)
	thePilot.SetDashboardChan(dashboardChan)
	thePilot.SetAlarmChan(alarmChan)
	thePilot.SetSteeringChan(steeringChan)

	////////////////////////////////////////
	// input stuffs
	////////////////////////////////////////
	switchGpio := func(pin byte) gpio.Gpio {

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

		err = g.SetDirection(gpio.IN)
		if err != nil {
			log.Fatal(err)
		}

		err = g.SetActiveLevel(gpio.ACTIVE_HIGH)
		if err != nil {
			log.Fatal(err)
		}

		// Test we can read it
		value, err := g.Value()
		if err != nil {
			log.Fatal(err)
		}

		switchState := "OFF"
		if value {
			switchState = "ON"
		}

		log.Info("[AUTOTEST] current switch position is %s", switchState)

		return g
	}(conf.SwitchGpioPin)
	control := control.New(switchGpio, thePilot)

	// TODO(ssoudan) if the value of the switch at start is true, the trigger the alarm, we have rebooted while the pilot was enabled

	////////////////////////////////////////
	// gps stuffs
	////////////////////////////////////////
	gps := gps.New(conf.GpsSerialPort)
	gps.SetMessagesChan(pilotChan)
	gps.SetErrorChan(pilotChan)

	gps.Start()
	dashboard.Start()
	alarm.Start()
	thePilot.Start()
	steering.Start()
	control.Start()

	// Wait until we receive a signal
	waitForInterrupt()

	dashboard.Shutdown()
	alarm.Shutdown()
	steering.Shutdown()
	control.Shutdown()
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
