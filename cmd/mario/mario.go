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
* @Date:   2015-09-22 13:24:54
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-03 01:58:40
 */

package main

import (
	"github.com/jessevdk/go-flags"

	"time"

	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("motorControl")

type Options struct {
	Clockwise bool `short:"c" long:"clockwise" description:"clockwise rotation" default:"false"`
	// Speed     uint32  `short:"s" long:"speed" description:"rotation speed pps" default:"600"`
	// Duration  float64 `short:"d" long:"duration" description:"duration (seconds)" default:"1"`
}

var opts Options

var parser = flags.NewParser(&opts, flags.Default)

func step(motor *motor.Motor, clockwise bool, stepsBySecond uint32, duration time.Duration) {
	motor.Enable()
	log.Info("Moving clockwise[%v] for %v at %v[steps/s]", clockwise, duration, stepsBySecond)
	motor.Move(clockwise, stepsBySecond, duration)
	motor.Disable()
}

type step_ struct {
	stepsBySecond uint32
	duration      time.Duration
	pause         time.Duration
}

func main() {

	if _, err := parser.Parse(); err != nil {
		log.Fatalf("failed to parse options: %v", err)
	}

	log.Info("%v", opts)
	motor := motor.New(
		conf.MotorStepPin,
		conf.MotorStepPwm,
		conf.MotorDirPin,
		conf.MotorSleepPin)

	sol := step_{392, time.Duration(250 * time.Millisecond), time.Duration(0 * time.Millisecond)}
	sol__ := step_{392, time.Duration(1000 * time.Millisecond), time.Duration(0 * time.Millisecond)}
	la := step_{440, time.Duration(250 * time.Millisecond), time.Duration(0 * time.Millisecond)}
	la_ := step_{440, time.Duration(500 * time.Millisecond), time.Duration(0 * time.Millisecond)}
	si := step_{494, time.Duration(500 * time.Millisecond), time.Duration(0 * time.Millisecond)}
	si_ := step_{494, time.Duration(250 * time.Millisecond), time.Duration(0 * time.Millisecond)}

	steps := []step_{
		sol,
		sol,
		sol,
		la,
		si,
		la_,
		sol,
		si_,
		la,
		la,
		sol__,
		// {587, time.Duration(125 * time.Millisecond), time.Duration(125 * time.Millisecond)},
		// {587, time.Duration(125 * time.Millisecond), time.Duration(250 * time.Millisecond)},
		// {587, time.Duration(125 * time.Millisecond), time.Duration(250 * time.Millisecond)},
		// {494, time.Duration(200 * time.Millisecond), time.Duration(250 * time.Millisecond)},
		// {587, time.Duration(200 * time.Millisecond), time.Duration(250 * time.Millisecond)},
		// {699, time.Duration(200 * time.Millisecond), time.Duration(250 * time.Millisecond)},
		// {587, time.Duration(125 * time.Millisecond), time.Duration(125 * time.Millisecond)},
	}

	for _, s := range steps {
		step(motor, opts.Clockwise, s.stepsBySecond, s.duration)
		time.Sleep(250 * time.Millisecond)
	}

}
