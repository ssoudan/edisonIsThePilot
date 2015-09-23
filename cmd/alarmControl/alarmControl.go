/*
* @Author: Sebastien Soudan
* @Date:   2015-09-23 11:37:24
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-23 11:40:38
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

		if err = pwm.Enable(); err != nil {
			log.Fatal(err)
		}
		log.Info("[AUTOTEST] alarm is ON")

		time.Sleep(2 * time.Second)
		if err = pwm.Disable(); err != nil {
			log.Fatal(err)
		}
		log.Info("[AUTOTEST] alarm is OFF")

		return pwm
	}(conf.AlarmGpioPin, conf.AlarmGpioPWM)
	defer alarmPwm.Unexport()

	alarm_ := alarm.New(alarmPwm)
	alarmChan := make(chan interface{})
	alarm_.SetInputChan(alarmChan)

	alarm_.Start()

	alarmChan <- alarm.NewMessage(true)

}
