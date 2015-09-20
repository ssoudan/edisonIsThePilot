/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 09:58:02
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 16:52:46
 */

package pilot

import (
	"github.com/ssoudan/edisonIsThePilot/dashboard"
)

type Pilot struct {
	alarm   Alarm
	heading float64
	bound   float64
	enabled bool

	dashboardChannel chan dashboard.Message

	leds map[string]bool
}

type Leds map[Led]bool
type Led string

type GPSFeedBack struct {
	Heading  float64
	Validity bool
	Speed    float64
}

type FixStatus byte

func checkFixStatus(fix FixStatus) (alarm Alarm, ledEnabled bool) {
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

const (
	NOFIX    = 0
	FIX      = 1
	DGPS_FIX = 2
)

func (p *Pilot) checkHeadingError(headingError float64) Alarm {

	inputStatus := validateInput(p.bound, headingError)
	return computeAlarmState(p.alarm, inputStatus)

}

func computeHeadingError(heading float64, gpsHeading float64) float64 {

	headingError := gpsHeading - heading
	if headingError > 180. {
		headingError -= 360.
	}

	if headingError <= -180. {
		headingError += 360.
	}

	return headingError
}

func (p *Pilot) UpdateFixStatus(fix FixStatus) {
	// compute the update for fix status
	fixAlarm, fixLed := checkFixStatus(fix)

	/////////////////////////
	// Update pilot state from previous checks
	////////////////////////
	p.alarm = p.alarm || fixAlarm

	// TODO(ssoudan) wrap this stuff in something that can be tested
	if fixLed {
		p.leds[dashboard.NoGPSFix] = true
	} else {
		delete(p.leds, dashboard.NoGPSFix)
	}

	/////////////////////////
	// Tell the world
	/////////////////////////
	p.dashboardChannel <- dashboard.NewMessage(p.leds)
}

func (p *Pilot) UpdateFeedback(gpsHeading GPSFeedBack) {

	// TODO(ssoudan) do something with the validity

	// TODO(ssoudan) do something with the speed

	headingError := computeHeadingError(p.heading, gpsHeading.Heading)

	headingAlarm := p.checkHeadingError(headingError)

	/////////////////////////
	// Update pilot state from previous checks
	////////////////////////

	// Update alarm state from the previously computed alarms
	p.alarm = Alarm(p.enabled) && headingAlarm // || blah

	// Update alarm state from the previously computed alarms
	steeringEnabled := p.computeSteeringState()

	/////////////////////////
	// Tell the world
	/////////////////////////
	if steeringEnabled {
		// TODO(ssoudan) do something with the heading error

		// TODO(ssoudan) call the PID

		// TODO(ssoudan) check the PID output
	}

	// p.tellTheWorld()
	p.dashboardChannel <- dashboard.NewMessage(p.leds)
}

const (
	INVALID = false
	VALID   = true
)

type InputStatus bool

func (p Pilot) computeSteeringState() bool {
	return p.enabled && !bool(p.alarm)
}

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

func computeAlarmState(previousState Alarm, input InputStatus) Alarm {
	if previousState == RAISED {
		return RAISED
	}

	if input == INVALID {
		return RAISED
	}

	return UNRAISED
}
