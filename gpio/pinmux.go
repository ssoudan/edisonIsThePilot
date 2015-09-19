/*
* @Author: Sebastien Soudan
* @Date:   2015-09-19 12:30:26
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-19 14:37:03
 */

package gpio

import (
	"fmt"
)

// TODO(ssoudan) need a map to translate between the different name of the pins

// EnablePWM enables PWM on a mux-ed pin
// See http://www.emutexlabs.com/project/215-intel-edison-gpio-pin-multiplexing-guide
// TODO(ssoudan) complete this
func EnablePWM(pin byte) error {
	return writeTo(fmt.Sprintf("/sys/kernel/debug/gpio_debug/gpio%d/current_pinmux", pin), "mode1")
}

// EnableGPIO enables GPIO mode on a mux-ed pin
func EnableGPIO(pin byte) error {
	return writeTo(fmt.Sprintf("/sys/kernel/debug/gpio_debug/gpio%d/current_pinmux", pin), "mode0")
}
