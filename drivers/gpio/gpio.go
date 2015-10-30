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
* @Last Modified time: 2015-10-22 14:59:12
 */

package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	// InDirection is the direction for input GPIO
	InDirection = "in"
	// OutDirection is the direction for output GPIO
	OutDirection = "out"
)

const (
	// ActiveHigh is the state for active high input GPIO
	ActiveHigh = "0"
	// ActiveLow is the state for active low input GPIO
	ActiveLow = "1"
)

const (
	high = "1\n"
	low  = "0\n"
)

const (
	sysfsGpioValue     = "/sys/class/gpio/gpio%d/value"
	sysfsGpioDirection = "/sys/class/gpio/gpio%d/direction"
	sysfsGpioDir       = "/sys/class/gpio/gpio%d/"
	sysfsGpioExport    = "/sys/class/gpio/export"
	sysfsGpioUnexport  = "/sys/class/gpio/unexport"
	sysfsGpioActiveLow = "/sys/class/gpio/gpio%d/active_low"
)

// Gpio is a general purpose I/O of the Edison
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
	if _, err := os.Stat(fmt.Sprintf(sysfsGpioDir, p.pin)); os.IsNotExist(err) {
		return false
	}
	return true
}

// Export the gpio to be usable from sysfs.
func (p Gpio) Export() error {
	return writeTo(sysfsGpioExport, fmt.Sprintf("%d", p.pin))
}

// Unexport the gpio from sysfs.
func (p Gpio) Unexport() error {
	return writeTo(sysfsGpioUnexport, fmt.Sprintf("%d", p.pin))
}

// SetDirection defines whether this particular GPIO is used for input or output (use constants InDirection and OutDirection).
func (p Gpio) SetDirection(dir string) error {
	if dir != InDirection && dir != OutDirection {
		return fmt.Errorf("Incorrect direction: %s", dir)
	}
	return writeTo(fmt.Sprintf(sysfsGpioDirection, p.pin), dir)
}

// SetActiveLevel set the ActiveLow or ActiveHigh level for IN direction.
func (p Gpio) SetActiveLevel(level string) error {
	if level != ActiveHigh && level != ActiveLow {
		return fmt.Errorf("Incorrect active level: %s", level)
	}
	return writeTo(fmt.Sprintf(sysfsGpioActiveLow, p.pin), level)
}

// Value returns the value of the GPIO
func (p Gpio) Value() (bool, error) {
	val, err := readfrom(fmt.Sprintf(sysfsGpioValue, p.pin))
	if err != nil {
		return false, err
	}
	switch val {
	case high:
		return true, nil
	case low:
		return false, nil
	default:
		return false, fmt.Errorf("invalid value: [%v]", val)
	}

}

// Enable this gpio
func (p Gpio) Enable() error {
	return writeTo(fmt.Sprintf(sysfsGpioValue, p.pin), high)
}

// Disable this gpio
func (p Gpio) Disable() error {
	return writeTo(fmt.Sprintf(sysfsGpioValue, p.pin), low)
}
