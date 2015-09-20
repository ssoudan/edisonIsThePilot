/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 12:07:14
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 12:09:09
 */

package pilot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatInvalidStatusLeadsToAnAlarm(t *testing.T) {
	previousState := Alarm(UNRAISED)
	input := InputStatus(INVALID)
	nextState := updateAlarmState(previousState, input)

	expected := RAISED

	assert.EqualValues(t, expected, nextState)
}

func TestThatInvalidStatusLeadsToAnAlarmWhenAlreadyRAISED(t *testing.T) {
	previousState := Alarm(RAISED)
	input := InputStatus(INVALID)
	nextState := updateAlarmState(previousState, input)

	expected := RAISED

	assert.EqualValues(t, expected, nextState)
}

func TestThatValidStatusLeadsToAnAlarmWhenAlreadyRAISED(t *testing.T) {
	previousState := Alarm(RAISED)
	input := InputStatus(VALID)
	nextState := updateAlarmState(previousState, input)

	expected := RAISED

	assert.EqualValues(t, expected, nextState)
}

func TestThatValidInputLeavesUnraisedAlarmUnraised(t *testing.T) {

	var previousState Alarm
	previousState = UNRAISED

	var input InputStatus
	input = VALID

	expected := UNRAISED
	result := updateAlarmState(previousState, input)

	assert.EqualValues(t, expected, result)
}
