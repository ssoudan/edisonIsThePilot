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
* @Date:   2015-09-19 12:30:26
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 13:13:17
 */

package gpio

import (
	"fmt"
)

const (
	sysfsPinmuxMode = "/sys/kernel/debug/gpio_debug/gpio%d/current_pinmux"
	sysfsI2CSpeed   = "/sys/class/i2c-adapter/i2c-%d/device/i2c_dw_sysnode/mode"
)

// EnablePWM enables PWM on a mux-ed pin
// See http://www.emutexlabs.com/project/215-intel-edison-gpio-pin-multiplexing-guide
func EnablePWM(pin byte) error {
	return writeTo(fmt.Sprintf(sysfsPinmuxMode, pin), "mode1")
}

// EnableGPIO enables GPIO mode on a mux-ed pin
func EnableGPIO(pin byte) error {
	return writeTo(fmt.Sprintf(sysfsPinmuxMode, pin), "mode0")
}

// EnableI2C enables i2c mode on a mux-ed pin
func EnableI2C(pin byte) error {
	return writeTo(fmt.Sprintf(sysfsPinmuxMode, pin), "mode1")
}

// EnableFastI2C enables i2c fast mode
func EnableFastI2C(bus byte) error {
	return writeTo(fmt.Sprintf(sysfsI2CSpeed, bus), "fast")
}

// EnableStdI2C enables i2c fast mode
func EnableStdI2C(bus byte) error {
	return writeTo(fmt.Sprintf(sysfsI2CSpeed, bus), "std")
}

// EnableHighI2C enables i2c fast mode
func EnableHighI2C(bus byte) error {
	return writeTo(fmt.Sprintf(sysfsI2CSpeed, bus), "high")
}
