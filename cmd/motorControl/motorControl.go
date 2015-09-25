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
* @Last Modified time: 2015-09-24 14:58:11
 */

package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("motorControl")

func step(motor *motor.Motor, clockwise bool, stepsBySecond uint32, duration time.Duration) {
	motor.Enable()
	log.Info("Moving clockwise[%v] for %v at %v[steps/s]", clockwise, duration, stepsBySecond)
	motor.Move(clockwise, stepsBySecond, duration)
	motor.Disable()
}

func main() {

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
		{true, 200, 5 * time.Second},
		{true, 400, 5 * time.Second},
		{false, 200, 5 * time.Second},
		{false, 400, 5 * time.Second},
	}

	for _, s := range steps {
		step(motor, s.clockwise, s.stepsBySecond, s.duration)
		time.Sleep(10 * time.Second)
	}

}
