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
* @Date:   2015-09-20 22:05:43
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 12:28:50
 */

package pilot

// InputStatus is the course status
type InputStatus bool

const (
	// INVALID course
	INVALID = false
	// VALID course
	VALID = true
)

func validateInput(bound, headingError float64) InputStatus {
	if -bound > headingError || bound < headingError {
		return INVALID
	}

	return VALID
}

func computeAlarmStateForInputStatus(previousState Alarm, input InputStatus) Alarm {
	if previousState == RAISED {
		return RAISED
	}

	if input == INVALID {
		return RAISED
	}

	return UNRAISED
}
