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
* @Date:   2015-09-20 22:08:20
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:17:28
 */

package pilot

type FixStatus byte

const (
	NOFIX    = 0
	FIX      = 1
	DGPS_FIX = 2
)

func validateFixStatus(fix FixStatus) (alarm Alarm, ledEnabled bool) {
	alarm = Alarm(UNRAISED)
	ledEnabled = false

	switch fix {
	case NOFIX:

		alarm = RAISED
		ledEnabled = true

	case FIX, DGPS_FIX:
		// blah
	}

	return
}
