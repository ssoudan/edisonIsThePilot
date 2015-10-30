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
* @Date:   2015-09-23 11:37:24
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 14:21:42
 */

package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/alarm"
	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/drivers/pwm"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("alarmControl")

func main() {

	panicChan := make(chan interface{})
	go func() {
		select {
		case m := <-panicChan:
			log.Fatal("Received a panic error - exiting: %v", m)
		}
	}()

	////////////////////////////////////////
	// a nice and delicate alarm
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

		return pwm
	}(conf.AlarmGpioPin, conf.AlarmGpioPWM)
	// defer alarmPwm.Unexport() // Don't do that it can disable the alarms for the autopilot program

	theAlarm := alarm.New(alarmPwm)
	alarmChan := make(chan interface{})
	theAlarm.SetInputChan(alarmChan)
	theAlarm.SetPanicChan(panicChan)

	theAlarm.Start()

	alarmChan <- alarm.NewMessage(true)

	for !theAlarm.Enabled() {
		log.Info("Waiting for the alarm to come")
		time.Sleep(10 * time.Millisecond)
	}

	theAlarm.Shutdown()

}
