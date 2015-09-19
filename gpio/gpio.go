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
* @Last Modified time: 2015-09-19 13:34:51
 */

package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	IN          = "in"
	OUT         = "out"
	ACTIVE_HIGH = "0"
	ACTIVE_LOW  = "1"
	HIGH        = "1"
	LOW         = "0"
)

type Gpio struct {
	pin uint8
}

// New returns a new PWM for a given pin - See Edison Breakout documentation to figure out which one you want.
func New(pin uint8) Gpio {
	return Gpio{pin: pin}
}

func writeTo(filename string, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func readfrom(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// IsExported returns true with the gpio is already exported and usable from sysfs.
func (p Gpio) IsExported() bool {
	if _, err := os.Stat(fmt.Sprintf("/sys/class/gpio/gpio%d/", p.pin)); os.IsNotExist(err) {
		return false
	}
	return true
}

// Export the gpio to be usable from sysfs.
func (p Gpio) Export() error {
	return writeTo("/sys/class/gpio//export", fmt.Sprintf("%d", p.pin))
}

// Unexport the gpio from sysfs.
func (p Gpio) Unexport() error {
	return writeTo("/sys/class/gpio/unexport", fmt.Sprintf("%d", p.pin))
}

// SetDirection defines wether this particular GPIO is used for input or output (use constants IN and OUT).
func (p Gpio) SetDirection(dir string) error {
	if dir != IN && dir != OUT {
		return fmt.Errorf("Incorrect direction: %s", dir)
	}
	return writeTo(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.pin), dir)
}

// SetActiveLevel set the HIGH or LOW active level for IN direction.
func (p Gpio) SetActiveLevel(level string) error {
	if level != ACTIVE_HIGH && level != ACTIVE_LOW {
		return fmt.Errorf("Incorrect active level: %s", level)
	}
	return writeTo(fmt.Sprintf("/sys/class/gpio/gpio%d/active_low", p.pin), level)
}

// GetValue returns the value of the GPIO
func (p Gpio) GetValue() (bool, error) {
	val, err := readfrom(fmt.Sprintf("/sys/class/gpio/gpio%d/active_low", p.pin))
	if err != nil {
		return false, err
	}
	switch val {
	default:
		return false, fmt.Errorf("invalid value")
	case HIGH:
		return true, nil
	case LOW:
		return false, nil
	}

}

// Enable this gpio
func (p Gpio) Enable() error {
	return writeTo(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.pin), "1")
}

// Disable this gpio
func (p Gpio) Disable() error {
	return writeTo(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.pin), "0")
}
