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
* @Last Modified time: 2015-10-21 12:59:58
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

// Options are the command line options for this tool
type Options struct {
	Clockwise  bool    `short:"c" long:"clockwise" description:"clockwise rotation" default:"false"`
	Speed      uint32  `short:"s" long:"speed" description:"rotation speed pps" default:"600"`
	Duration   float64 `short:"d" long:"duration" description:"duration (seconds)" default:"1"`
	Repetition int     `short:"r" long:"rep" description:"repetitions" default:"1"`
	Pause      float64 `short:"p" long:"pause" description:"pause" default:"0"`
}

var opts Options

var parser = flags.NewParser(&opts, flags.Default)

func step(motor *motor.Motor, clockwise bool, stepsBySecond uint32, duration time.Duration) {
	// motor.Enable()
	log.Info("Moving clockwise[%v] for %v at %v[steps/s]", clockwise, duration, stepsBySecond)
	motor.Move(clockwise, stepsBySecond, duration)
	// motor.Disable()
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

	steps := []struct {
		clockwise     bool
		stepsBySecond uint32
		duration      time.Duration
	}{
		// {true, 100, time.Duration(0.4 * float64(time.Second))},
		// {true, 200, time.Duration(0.8 * float64(time.Second))},

		{opts.Clockwise, opts.Speed, time.Duration(opts.Duration * float64(time.Second))},
		// {opts.Clockwise, opts.Speed / 2, time.Duration(opts.Duration * float64(time.Second) / 4)},
		// {true, 400, 5 * time.Second},
		// {false, 200, 5 * time.Second},
		// {false, 400, 5 * time.Second},
	}
	motor.Enable()
	for i := 0; i < opts.Repetition; i++ {
		for _, s := range steps {
			step(motor, s.clockwise, s.stepsBySecond, s.duration)
		}
		time.Sleep(time.Duration(opts.Pause * float64(time.Second)))
	}
	motor.Disable()

}
