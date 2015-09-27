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
* @Date:   2015-09-18 14:10:18
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-22 13:08:03
 */

package pwm

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
)

const (
	maxPeriodNanoSec = 218453000
	MaxPeriod        = maxPeriodNanoSec * time.Nanosecond
	minPeriodNanoSec = 104
	MinPeriod        = minPeriodNanoSec * time.Nanosecond
)

type Pwm struct {
	pwmId  uint8
	pwmPin byte
	gpio   gpio.Gpio
}

// New returns a new PWM for a given pwmId - See Edison Breakout documentation to figure out which one you want.
func New(pwmId uint8, pwmPin byte) (*Pwm, error) {
	var err error
	var g = gpio.New(pwmPin)
	if !g.IsExported() {
		err = g.Export()
		if err != nil {
			return nil, err
		}
	}

	err = g.SetDirection(gpio.OUT)
	if err != nil {
		return nil, err
	}

	err = g.Disable()
	if err != nil {
		return nil, err
	}

	return &Pwm{pwmId: pwmId, pwmPin: pwmPin, gpio: g}, nil
}

func writeTo(filename string, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

// IsExported returns true with the pwm is already exported and usable from sysfs.
func (p Pwm) IsExported() bool {
	if _, err := os.Stat(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/", p.pwmId)); os.IsNotExist(err) {
		return false
	}
	return true
}

// Export the pwm to be usable from sysfs.
func (p Pwm) Export() error {

	return writeTo("/sys/class/pwm/pwmchip0/export", fmt.Sprintf("%d", p.pwmId))
}

// Unexport the pwm from sysfs.
func (p Pwm) Unexport() error {
	return writeTo("/sys/class/pwm/pwmchip0/unexport", fmt.Sprintf("%d", p.pwmId))
}

// SetPeriodAndDutyCycle configures the pwm for a given period and duty cycle ratio.
// Note this might go through a transient state if the pwm is Enabled
func (p *Pwm) SetPeriodAndDutyCycle(period time.Duration, duty_cycle float32) error {
	if period < MinPeriod || period > MaxPeriod {
		return fmt.Errorf("must be in 104:218453000 ns range")
	}

	if duty_cycle < 0 || duty_cycle > 1 {
		return fmt.Errorf("must be in 0:1 range")
	}

	if err := p.setDutyCycleNanoSec(1); err != nil {
		return err
	}
	if err := p.setPeriodNanoSecond(period.Nanoseconds()); err != nil {
		return err
	}

	dutyCycle := (int64)(float32(period.Nanoseconds()) * duty_cycle)

	if err := p.setDutyCycleNanoSec(dutyCycle); err != nil {
		return err
	}

	return nil
}

func (p Pwm) setDutyCycleNanoSec(duty_cycle int64) error {
	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/duty_cycle", p.pwmId), fmt.Sprintf("%d", duty_cycle))
}

func (p Pwm) setPeriodNanoSecond(period int64) error {

	if period > maxPeriodNanoSec || period < minPeriodNanoSec {
		return fmt.Errorf("must be in 104:218453000 range")
	}

	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/period", p.pwmId), fmt.Sprintf("%d", period))
}

// Enable this pwm
func (p Pwm) Enable() error {
	err := gpio.EnablePWM(p.pwmPin)
	if err != nil {
		return err
	}

	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/enable", p.pwmId), "1")
}

// Disable this pwm
func (p Pwm) Disable() error {
	err := gpio.EnableGPIO(p.pwmPin)
	if err != nil {
		return err
	}

	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/enable", p.pwmId), "0")
}