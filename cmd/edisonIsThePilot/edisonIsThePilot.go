/*
* @Author: Sebastien Soudan
* @Date:   2015-09-18 12:20:59
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-18 17:37:39
 */

package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/gps"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/pwm"

	// "bitbucket.org/gmcbay/i2c"

	"github.com/adrianmo/go-nmea"
)

var log = logger.Log("edisonIsThePilot")

func main() {
	var err error

	////////////////////////////////////////
	// PWM stuffs
	////////////////////////////////////////
	var pwm = pwm.Pwm{Pin: 2}
	if !pwm.IsExported() {
		err = pwm.Export()
		if err != nil {
			log.Fatal(err)
		}
	}

	pwm.Disable()

	if err = pwm.SetPeriodAndDutyCycle(50*time.Millisecond, 0.5); err != nil {
		log.Fatal(err)
	}

	if err = pwm.Enable(); err != nil {
		log.Fatal(err)
	}
	log.Info("pwm configured")

	////////////////////////////////////////
	// I2C stuffs
	////////////////////////////////////////

	// bp, err := i2c.Bus(1)
	// if err != nil {
	// 	log.Panicf("failed to create bus: %v\n", err)
	// }
	// addr := 0x12
	// reg := 0x13
	// len := 0x1
	// data, err := bp.ReadByteBlock(addr, reg, length)

	////////////////////////////////////////
	// gps stuffs
	////////////////////////////////////////
	gps := gps.New("/dev/ttyMFD1")
	messagesChan, errorChan := gps.Stream()

	for {
		select {
		case m := <-messagesChan:
			switch t := m.(type) {
			default:
				// don't care
				log.Debug("%+v\n", m)
			case nmea.GPRMC:
				log.Info("[GPRMC] validity: %v heading: %v[˚] speed: %v[knots] \n", t.Validity == "A", t.Course, t.Speed)
			}
		case err := <-errorChan:
			log.Warning("Error while processing GPS message: %v", err)
		}
	}

	pwm.Disable()
	pwm.Unexport()

}
