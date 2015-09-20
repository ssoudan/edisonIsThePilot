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
* @Date:   2015-09-18 12:20:59
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 21:42:36
 */

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/ssoudan/edisonIsThePilot/compass/hmc"
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	// "github.com/ssoudan/edisonIsThePilot/gpio"
	"github.com/ssoudan/edisonIsThePilot/gps"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	// "github.com/ssoudan/edisonIsThePilot/pwm"
)

var log = logger.Log("edisonIsThePilot")

func main() {
	// var err error

	////////////////////////////////////////
	// PWM stuffs
	////////////////////////////////////////
	// err = gpio.EnablePWM(182)
	// if err != nil {
	// 	log.Fatal("Failed to set pin 182 to pwm2: %v", err)
	// }
	// var pwm = pwm.New(2)
	// if !pwm.IsExported() {
	// 	err = pwm.Export()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// pwm.Disable()

	// if err = pwm.SetPeriodAndDutyCycle(10*time.Millisecond, 0.5); err != nil {
	// 	log.Fatal(err)
	// }

	// if err = pwm.Enable(); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info("pwm configured")
	// time.Sleep(100 * time.Second)

	// pwm.Disable()
	// pwm.Unexport()

	////////////////////////////////////////
	// GPIO stuffs
	////////////////////////////////////////
	// err = gpio.EnableGPIO(165)
	// if err != nil {
	// 	log.Fatal("Failed to set pin 182 to pwm2: %v", err)
	// }
	// var gpio165 = gpio.New(165)
	// if !gpio165.IsExported() {
	// 	log.Debug("GPIO165 not exported")
	// 	err = gpio165.Export()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// err = gpio165.SetDirection(gpio.OUT)
	// if err != nil {
	// 	log.Fatal("Failed to set direction: %v", err)
	// }

	// for i := 0; i < 4; i++ {
	// 	err = gpio165.Enable()
	// 	if err != nil {
	// 		log.Fatal("Failed to enable: %v", err)
	// 	}
	// 	log.Debug("GPIO up")
	// 	time.Sleep(1 * time.Second)
	// 	err = gpio165.Disable()
	// 	if err != nil {
	// 		log.Fatal("Failed to disable: %v", err)
	// 	}
	// 	log.Debug("GPIO down")
	// 	time.Sleep(1 * time.Second)
	// }
	// err = gpio165.Disable()
	// if err != nil {
	// 	log.Fatal("Failed to disable: %v", err)
	// }

	// err = gpio165.Unexport()
	// if err != nil {
	// 	log.Fatal("Failed to unexport: %v", err)
	// }

	////////////////////////////////////////
	// HMC5883 stuffs
	////////////////////////////////////////
	// compass := hmc.New(6)
	// for !compass.Begin() {

	// }

	// // Set measurement range
	// compass.SetRange(hmc.HMC5883L_RANGE_1_3GA)

	// // Set measurement mode
	// compass.SetMeasurementMode(hmc.HMC5883L_CONTINOUS)

	// // Set data rate
	// compass.SetDataRate(hmc.HMC5883L_DATARATE_3HZ)

	// // Set number of samples averaged
	// compass.SetSamples(hmc.HMC5883L_SAMPLES_8)

	// // Set calibration offset. See HMC5883L_calibration.ino
	// compass.SetOffset(-82, 72)

	// mag, err := compass.ReadNormalize()
	// if err == nil {
	// 	log.Info("Compass reading is %v", mag)
	// }

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
	// a beautiful dashboard
	////////////////////////////////////////
	dashboard := dashboard.New()
	dashboardChan := make(chan interface{})
	dashboard.SetInputChan(dashboardChan)

	////////////////////////////////////////
	// pilot stuffs
	////////////////////////////////////////
	thePilot := pilot.New(15.)
	pilotChan := make(chan interface{})
	thePilot.SetInputChan(pilotChan)
	thePilot.SetDashboardChan(dashboardChan)

	////////////////////////////////////////
	// gps stuffs
	////////////////////////////////////////
	gps := gps.New("/dev/ttyMFD1")
	gps.SetMessagesChan(pilotChan)
	gps.SetErrorChan(pilotChan)

	gps.Start()
	dashboard.Start()
	thePilot.Start()

	// For tests
	go func() {
		for {
			log.Notice("Disabling the pilot")
			thePilot.Disable()
			time.Sleep(35 * time.Second)
			log.Notice("Enabling the pilot")
			thePilot.Enable()
			time.Sleep(35 * time.Second)
		}
	}()

	// Wait until we receive a signal
	waitForInterrupt()

}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigChan:
		log.Info("Interrupted - exiting")
		os.Exit(255)
	}
}
