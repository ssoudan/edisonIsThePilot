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
* @Last Modified time: 2015-10-21 16:48:53
 */

package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/alarm"
	"github.com/ssoudan/edisonIsThePilot/ap100"
	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/control"
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/drivers/pwm"
	"github.com/ssoudan/edisonIsThePilot/drivers/sincos"
	"github.com/ssoudan/edisonIsThePilot/gps"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/pid"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/utils"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/webserver"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	"github.com/ssoudan/edisonIsThePilot/steering"
	"github.com/ssoudan/edisonIsThePilot/tracer"
)

var log = logger.Log("edisonIsThePilot")

// Version is the version of this code -- sets at compilation time
var Version = "unknown"

func main() {
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
			// The motor
			motor := motor.New(
				conf.MotorStepPin,
				conf.MotorStepPwm,
				conf.MotorDirPin,
				conf.MotorSleepPin)
			if err := motor.Disable(); err != nil {
				log.Error("Failed to stop the motor")
			}
			motor.Unexport()

			log.Fatalf("Version %v -- Received a panic error -- exiting: %v", Version, m)
		}
	}()

	ws := webserver.New(Version)
	ws.SetPanicChan(panicChan)
	ws.Start()

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

		err = g.SetDirection(gpio.OutDirection)
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
	for _, v := range conf.MessageToPin {
		g := mapMessageToGPIO(v.Message, v.Pin)
		dashboardGPIOs[v.Message] = g
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

		err = g.SetDirection(gpio.InDirection)
		if err != nil {
			log.Panic(err)
		}

		err = g.SetActiveLevel(gpio.ActiveHigh)
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
	defer motor.Disable()
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

	// The compass sincos output interface
	compass := sincos.New(conf.I2CBus, conf.SinAddress, conf.CosAddress)

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
	ws.SetDashboard(dashboard)

	////////////////////////////////////////
	// an astonishing steering
	////////////////////////////////////////
	steering := steering.New(motor)
	steeringChan := make(chan interface{})
	steering.SetInputChan(steeringChan)
	steering.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a stunning tracer
	////////////////////////////////////////
	tracer := tracer.New(conf.Conf.TraceSize)
	tracerChan := make(chan interface{})
	tracer.SetInputChan(tracerChan)
	tracer.SetPanicChan(panicChan)
	ws.SetTracer(tracer)

	////////////////////////////////////////
	// an amazing PID
	////////////////////////////////////////
	pidController := pid.New(
		conf.Conf.P,
		conf.Conf.I,
		conf.Conf.D,
		conf.Conf.N,
		conf.Conf.MinPIDOutputLimits,
		conf.Conf.MaxPIDOutputLimits)

	////////////////////////////////////////
	// a great pilot
	////////////////////////////////////////
	thePilot := pilot.New(pidController, conf.Conf.Bounds)
	pilotChan := make(chan interface{})
	thePilot.SetInputChan(pilotChan)
	thePilot.SetDashboardChan(dashboardChan)
	thePilot.SetAlarmChan(alarmChan)
	thePilot.SetSteeringChan(steeringChan)
	thePilot.SetPanicChan(panicChan)
	ws.SetPilot(thePilot)

	////////////////////////////////////////
	// a surprising input
	////////////////////////////////////////
	control := control.New(switchGpio, thePilot)
	control.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a friendly interface to the AP100
	////////////////////////////////////////
	ap100 := ap100.New(compass)
	headingChan := make(chan interface{})
	ap100.SetInputChan(headingChan)
	ap100.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a wonderful gps
	////////////////////////////////////////
	gps := gps.New(conf.Conf.GpsSerialPort)
	gps.SetMessagesChan(pilotChan)
	gps.SetHeadingChan(headingChan)
	gps.SetErrorChan(pilotChan)
	gps.SetPanicChan(panicChan)
	gps.SetTracerChan(tracerChan)

	tracer.Start()
	defer tracer.Shutdown()
	ap100.Start()
	defer ap100.Shutdown()
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
