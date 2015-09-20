/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 09:58:02
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 12:05:28
 */

package pilot

import ()

type Pilot struct {
	alarm   Alarm
	heading float64
	bound   float64
}

func (p *Pilot) checkHeadingError(headingError float64) {

	inputStatus := validateInput(p.bound, headingError)
	p.alarm = updateAlarmState(p.alarm, inputStatus)

}

func (p Pilot) computeHeadingError(gpsHeading float64) float64 {

	headingError := gpsHeading - p.heading
	if headingError > 180. {
		headingError -= 360.
	}

	if headingError <= -180. {
		headingError += 360.
	}

	return headingError
}

func (p *Pilot) UpdateInput(gpsHeading float64) {
	headingError := p.computeHeadingError(gpsHeading)
	p.checkHeadingError(headingError)
}

const (
	INVALID = false
	VALID   = true
)

type InputStatus bool

func validateInput(bound, headingError float64) InputStatus {
	if -bound > headingError || bound < headingError {
		return INVALID
	}

	return VALID
}

type Alarm bool

const (
	RAISED   = true
	UNRAISED = false
)

func updateAlarmState(previousState Alarm, input InputStatus) Alarm {
	if previousState == RAISED {
		return RAISED
	}

	if input == INVALID {
		return RAISED
	}

	return UNRAISED
}
