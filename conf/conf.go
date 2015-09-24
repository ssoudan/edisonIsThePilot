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
* @Last Modified time: 2015-09-24 14:12:11
 */

package conf

import (
	"github.com/ssoudan/edisonIsThePilot/dashboard"
)

var MessageToPin = map[string]byte{
	dashboard.NoGPSFix:                40, // J19 - pin 10
	dashboard.InvalidGPSData:          43, // J19 - pin 11
	dashboard.SpeedTooLow:             48, // J19 - pin 6
	dashboard.HeadingErrorOutOfBounds: 82, // J19 - pin 13
	dashboard.CorrectionAtLimit:       83, // J19 - pin 14
}

const (
	AlarmGpioPin  = 183 // J18 - pin 8
	AlarmGpioPWM  = 3
	MotorDirPin   = 165 // J18 - pin 2
	MotorSleepPin = 12  // J18 - pin 7
	MotorStepPin  = 182 // J17 - pin 1
	MotorStepPwm  = 2
	SwitchGpioPin = 46 // J19 - pin 5
)

const (
	MaxPIDOutputLimits             = 15
	MinPIDOutputLimits             = -15
	P                              = 1
	I                              = 0.1
	D                              = 0.1
	Bounds                         = 15.
	GpsSerialPort                  = "/dev/ttyMFD1"
	NoInputMessageTimeoutInSeconds = 10
	MinimumSpeedInKnots            = 3
)
