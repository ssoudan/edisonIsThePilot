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
* @Last Modified time: 2015-09-25 11:59:31
 */

package main

import (
	"fmt"
	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/drivers/motor"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"math"
	"time"
)

var log = logger.Log("motorCalibration")

var numberOfSteps = float64(200)

func rotationInDegreeToMove(speed uint32, rotationInDegree float64) (clockwise bool, duration time.Duration) {
	clockwise = rotationInDegree > 0.
	duration = time.Duration(
		math.Abs(rotationInDegree/360.*numberOfSteps/float64(speed)) * float64(time.Second))

	return
}
func doStep(motor *motor.Motor, s step) {

	clockwise, duration := rotationInDegreeToMove(s.stepsBySecond, s.rotationInDegree)

	motor.Enable()
	log.Info("[%3d] Moving [%6s] -- clockwise[%v] at %v[steps/s] for %v", s.id, fmt.Sprintf("%3.2f", s.rotationInDegree), clockwise, s.stepsBySecond, duration)
	motor.Move(clockwise, s.stepsBySecond, duration)
	motor.Disable()
}

type step struct {
	id               int
	stepsBySecond    uint32
	rotationInDegree float64
}

func main() {

	motor := motor.New(
		conf.MotorStepPin,
		conf.MotorStepPwm,
		conf.MotorDirPin,
		conf.MotorSleepPin)

	stepCount := 201

	steps := make([]step, stepCount)

	for i := 1; i < stepCount; i++ {
		steps[i] = step{id: i, stepsBySecond: 200, rotationInDegree: float64(i) * 1.8}
	}
	fmt.Println(steps)

	for _, s := range steps {
		doStep(motor, s)
		time.Sleep(510 * time.Millisecond)
	}

}
