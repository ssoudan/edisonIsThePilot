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
* @Last Modified time: 2015-09-25 11:29:32
 */

package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/alarm"
	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/control"
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/drivers/pwm"
	"github.com/ssoudan/edisonIsThePilot/gps"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/pid"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/utils"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	"github.com/ssoudan/edisonIsThePilot/steering"
)

var log = logger.Log("edisonIsThePilot")

var Version = "unknown"

func main() {
	// TODO(ssoudan) need to make sure there is a single process running! --> a socket could do the job

	log.Info("Starting -- version %s", Version)

	panicChan := make(chan interface{})
	defer func() {
		if r := recover(); r != nil {
			panicChan <- r
		}
	}()

	go func() {
		select {
		case m := <-panicChan:

			// kill the process (via log.Fatal) in case we can't create the PWM
			if pwm, err := pwm.New(conf.AlarmGpioPWM, conf.AlarmGpioPin); err == nil {
				if !pwm.IsExported() {
					err = pwm.Export()
					if err != nil {
						log.Error("Failed to raise the alarm")
					}
				}

				pwm.Enable()
			} else {
				log.Error("Failed to raise the alarm")
			}

			log.Fatalf("Version %v -- Received a panic error -- exiting: %v", Version, m)
		}
	}()

	////////////////////////////////////////
	// Init the IO
	////////////////////////////////////////
	// the LEDs
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

		log.Info("[AUTOTEST] %s LED is ON", message)
		time.Sleep(1 * time.Second)

		err = g.Enable()
		if err != nil {
			log.Panic(err)
		}

		log.Info("[AUTOTEST] %s LED is OFF", message)
		time.Sleep(1 * time.Second)

		err = g.Disable()
		if err != nil {
			log.Panic(err)
		}

		return g
	}
	dashboardGPIOs := make(map[string]gpio.Gpio, len(conf.MessageToPin))
	for k, v := range conf.MessageToPin {
		g := mapMessageToGPIO(k, v)
		dashboardGPIOs[k] = g
	}
	defer func() {
		for _, g := range dashboardGPIOs {
			g.Disable()
			g.Unexport()
		}
	}()

	// the input button
	switchGpio := func(pin byte) gpio.Gpio {

		// kill the process (via log.Panic) in case we can't create the GPIO
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

		err = g.SetDirection(gpio.IN)
		if err != nil {
			log.Panic(err)
		}

		err = g.SetActiveLevel(gpio.ACTIVE_HIGH)
		if err != nil {
			log.Panic(err)
		}

		// Test we can read it and make we we don't go beyond this point until the switch is OFF
		// to prevent the autopilot to be re-enabled when the system restart.
		// Since the alarm has not been initialized yet, after a reboot it will be ON.
		for value, err := g.Value(); value; value, err = g.Value() {
			if err != nil {
				log.Panic(err)
			}

			log.Info("[AUTOTEST] current autopilot switch position is ON - switch it OFF to proceed.")

			time.Sleep(200 * time.Millisecond)
		}

		return g
	}(conf.SwitchGpioPin)
	defer switchGpio.Unexport()

	// The motor
	motor := motor.New(
		conf.MotorStepPin,
		conf.MotorStepPwm,
		conf.MotorDirPin,
		conf.MotorSleepPin)
	defer motor.Unexport()

	// The alarm
	alarmPwm := func(pin byte, pwmId byte) *pwm.Pwm {

		// kill the process (via log.Panic) in case we can't create the PWM
		pwm, err := pwm.New(pwmId, pin)
		if err != nil {
			log.Panic(err)
		}

		if !pwm.IsExported() {
			err = pwm.Export()
			if err != nil {
				log.Panic(err)
			}
		}

		pwm.Disable()

		if err = pwm.SetPeriodAndDutyCycle(200*time.Millisecond, 0.5); err != nil {
			log.Panic(err)
		}

		if err = pwm.Enable(); err != nil {
			log.Panic(err)
		}
		log.Info("[AUTOTEST] alarm is ON")

		time.Sleep(2 * time.Second)
		if err = pwm.Disable(); err != nil {
			log.Panic(err)
		}
		log.Info("[AUTOTEST] alarm is OFF")

		return pwm
	}(conf.AlarmGpioPin, conf.AlarmGpioPWM)
	defer alarmPwm.Unexport()

	////////////////////////////////////////
	// a nice and delicate alarm
	////////////////////////////////////////

	alarm := alarm.New(alarmPwm)
	alarmChan := make(chan interface{})
	alarm.SetInputChan(alarmChan)
	alarm.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a beautiful dashboard
	////////////////////////////////////////
	dashboard := dashboard.New()
	dashboardChan := make(chan interface{})
	dashboard.SetInputChan(dashboardChan)
	dashboard.SetPanicChan(panicChan)
	for m, g := range dashboardGPIOs {
		dashboard.RegisterMessageHandler(m, g)
	}

	////////////////////////////////////////
	// an astonishing steering
	////////////////////////////////////////
	steering := steering.New(motor)
	steeringChan := make(chan interface{})
	steering.SetInputChan(steeringChan)
	steering.SetPanicChan(panicChan)

	////////////////////////////////////////
	// an amazing PID
	////////////////////////////////////////
	pidController := pid.New(
		conf.P,
		conf.I,
		conf.D,
		conf.N,
		conf.MinPIDOutputLimits,
		conf.MaxPIDOutputLimits)

	////////////////////////////////////////
	// a great pilot
	////////////////////////////////////////
	thePilot := pilot.New(pidController, conf.Bounds)
	pilotChan := make(chan interface{})
	thePilot.SetInputChan(pilotChan)
	thePilot.SetDashboardChan(dashboardChan)
	thePilot.SetAlarmChan(alarmChan)
	thePilot.SetSteeringChan(steeringChan)
	thePilot.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a surprising input
	////////////////////////////////////////
	control := control.New(switchGpio, thePilot)
	control.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a wonderful gps
	////////////////////////////////////////
	gps := gps.New(conf.GpsSerialPort)
	gps.SetMessagesChan(pilotChan)
	gps.SetErrorChan(pilotChan)
	gps.SetPanicChan(panicChan)

	gps.Start()
	control.Start()
	defer control.Shutdown()
	alarm.Start()
	defer alarm.Shutdown()
	dashboard.Start()
	defer dashboard.Shutdown()
	steering.Start()
	defer steering.Shutdown()
	thePilot.Start()
	defer thePilot.Shutdown()

	// Wait until we receive a signal
	utils.WaitForInterrupt(func() {
		log.Info("Interrupted - exiting")
		log.Info("Exiting -- version %v", Version)
	})

}
