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
* @Date:   2015-09-24 14:30:13
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 12:19:32
 */

package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForInterrupt is a blocking function that wait until SIGTERM or SIGQUIT signals are received
func WaitForInterrupt(f func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigChan:
		f()
	}
}
