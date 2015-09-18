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
* @Last Modified time: 2015-09-18 17:53:48
 */

package pwm

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Pwm struct {
	Pin uint8
}

func writeTo(filename string, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func (p Pwm) IsExported() bool {
	if _, err := os.Stat(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/", p.Pin)); os.IsNotExist(err) {
		return false
	}
	return true
}

func (p Pwm) Export() error {
	return writeTo("/sys/class/pwm/pwmchip0/export", fmt.Sprintf("%d", p.Pin))
}

func (p Pwm) Unexport() error {
	return writeTo("/sys/class/pwm/pwmchip0/unexport", fmt.Sprintf("%d", p.Pin))
}

func (p Pwm) SetPeriodAndDutyCycle(period time.Duration, duty_cycle float32) error {
	if period < 104*time.Nanosecond || period > 218453000*time.Nanosecond {
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
	if err := p.setDutyCycleNanoSec((int64)(float32(period.Nanoseconds()) * duty_cycle)); err != nil {
		return err
	}

	return nil
}

func (p Pwm) setDutyCycleNanoSec(duty_cycle int64) error {
	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/duty_cycle", p.Pin), fmt.Sprintf("%d", duty_cycle))
}

func (p Pwm) setPeriodNanoSecond(period int64) error {

	if period > 218453000 || period < 104 {
		return fmt.Errorf("must be in 104:218453000 range")
	}

	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/period", p.Pin), fmt.Sprintf("%d", period))
}

func (p Pwm) Enable() error {
	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/enable", p.Pin), "1")
}

func (p Pwm) Disable() error {
	return writeTo(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%d/enable", p.Pin), "0")
}
