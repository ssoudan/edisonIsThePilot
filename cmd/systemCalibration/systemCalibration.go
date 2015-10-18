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
* @Date:   2015-09-27 22:18:56
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-29 22:43:45
 */

package main

import (
	"github.com/jessevdk/go-flags"

	"net/http"
	"time"

	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/control"
	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/gps"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/utils"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/webserver"
	"github.com/ssoudan/edisonIsThePilot/steering"
	"github.com/ssoudan/edisonIsThePilot/stepper"
)

var log = logger.Log("systemCalibration")

var Version = "unknown"

type Options struct {
	Step        float64 `short:"s" long:"step" description:"step intensity (motor rotation in degree)" required:"true"`
	Duration    int64   `short:"d" long:"duration" description:"duration (seconds)" required:"true"`
	Description string  `short:"D" long:"description" description:"description of the test" required:"true"`
}

var opts Options

var parser = flags.NewParser(&opts, flags.Default)

func main() {

	// parse inputs
	if _, err := parser.Parse(); err != nil {
		log.Fatalf("failed to parse options: %v", err)
	}

	// Scenario input format
	// we want scenario with different speed
	// we want scenario with positive and negative steering changes
	// we want different values of those changes
	// we want different initial steering (do we?)
	// we want a scenario where we just go straight to record the noise
	// ---> all that can be defined with the step height and the duration

	log.Info("Starting -- version %s", Version)

	log.Info("Opts: %#v", opts)

	panicChan := make(chan interface{})
	defer func() {
		if r := recover(); r != nil {
			panicChan <- r
		}
	}()

	stepperChan := make(chan interface{})

	// The motor
	motor := motor.New(
		conf.MotorStepPin,
		conf.MotorStepPwm,
		conf.MotorDirPin,
		conf.MotorSleepPin)
	defer motor.Unexport()

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

	////////////////////////////////////////
	// an astonishing steering
	////////////////////////////////////////
	steering := steering.New(motor)
	steeringChan := make(chan interface{})
	steering.SetInputChan(steeringChan)
	steering.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a wonderful gps
	////////////////////////////////////////
	gps := gps.New(conf.Conf.GpsSerialPort)
	gps.SetMessagesChan(stepperChan)
	gps.SetErrorChan(stepperChan)
	gps.SetPanicChan(panicChan)

	////////////////////////////////////////
	// a crazy stepper
	////////////////////////////////////////
	stepper_ := stepper.New()
	stepper_.SetInputChan(stepperChan)
	stepper_.SetPanicChan(panicChan)
	stepper_.SetSteeringChan(steeringChan)

	////////////////////////////////////////
	// a surprising input
	////////////////////////////////////////
	control := control.New(switchGpio, stepper_)
	control.SetPanicChan(panicChan)

	// tell the pilot what we are going to do
	log.Notice(`When the autopilot button is switched to ON, we are going to acquire the 
current heading and start a steering step of %f degree and hold it for %v. As a pilot, 
you'll have to make sure there is enough place for that and the vessel is safe. 
Changing the steering during the test will make it invalide. 
When the test is over your are free to resume normal operations.`, opts.Step, opts.Duration)

	gps.Start()
	control.Start()
	steering.Start()
	stepper_.Start()
	defer steering.Shutdown()
	defer stepper_.Shutdown()
	defer control.Shutdown()

	go func() {
		select {
		case m := <-panicChan:

			// make sure the motor is stopped
			err := motor.Disable()
			if err != nil {
				log.Error("Failed to disable the motor - exiting anyway", err)
			}

			log.Fatalf("Version %v -- Received a panic error -- exiting: %v", Version, m)
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}
		}()

		http.HandleFunc("/", webserver.VersionEndpoint(Version))
		http.HandleFunc("/calibration", stepper_.CalibrationEndpoint)
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Panic("Already running! or something is living on port 8000 - exiting")
		}
	}()

	// populate the scenario
	stepperChan <- stepper.NewStep(opts.Step, time.Duration(opts.Duration)*time.Second, opts.Description)

	// FUTURE(ssoudan) We can imagine to generate the input from matlab? :)

	// FUTURE(ssoudan) support periodic response testing

	// Wait until we receive a signal
	utils.WaitForInterrupt(func() {
		log.Info("Interrupted - exiting")
		log.Info("Exiting -- version %v", Version)
	})
}
