/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 22:05:43
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:06:20
 */

package pilot

type InputStatus bool

const (
	INVALID = false
	VALID   = true
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
