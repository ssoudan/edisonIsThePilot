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
* @Date:   2015-09-22 13:18:01
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-03 11:01:36
 */

package conf

import (
	"github.com/ssoudan/edisonIsThePilot/dashboard"
)

type MessagePin struct {
	Message string
	Pin     byte
}

var MessageToPin = []MessagePin{
	{dashboard.NoGPSFix, 43},                // J19 - pin 11
	{dashboard.InvalidGPSData, 48},          // J19 - pin 6
	{dashboard.SpeedTooLow, 40},             // J19 - pin 10
	{dashboard.HeadingErrorOutOfBounds, 82}, // J19 - pin 13
	{dashboard.CorrectionAtLimit, 83},       // J19 - pin 14
}

const (
	AlarmGpioPin  = 183 // J18 - pin 8
	AlarmGpioPWM  = 3   // corresponding pwm number
	MotorDirPin   = 165 // J18 - pin 2
	MotorSleepPin = 12  // J18 - pin 7
	MotorStepPin  = 182 // J17 - pin 1
	MotorStepPwm  = 2   // corresponding pwm number
	SwitchGpioPin = 46  // J19 - pin 5
)

const (
	Bounds                 = 15.                          // error bound in degree
	SteeringReductionRatio = 380 / 25                     // reduction ratio between the motor and the steering wheel
	MaxPIDOutputLimits     = 15 * SteeringReductionRatio  // maximum pid output value (in degree)
	MinPIDOutputLimits     = -15 * SteeringReductionRatio // minimum pid output value (in degree)
	P                      = 0.0869670447979224           // Proportional coefficient
	I                      = 7.00335119238032e-05         // Integrative coefficient
	D                      = 21.9331758752948             // Derivative coefficient
	N                      = 0.136075524385029            // Derivative filter coefficient

	GpsSerialPort                  = "/dev/ttyMFD1"
	NoInputMessageTimeoutInSeconds = 10
	MinimumSpeedInKnots            = 3
)
