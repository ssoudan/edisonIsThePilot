/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 12:09:28
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:10:27
 */

package pilot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatInvalidStatusLeadsToAnAlarm(t *testing.T) {
	previousState := Alarm(UNRAISED)
	input := InputStatus(INVALID)
	nextState := computeAlarmStateForInputStatus(previousState, input)

	expected := RAISED

	assert.EqualValues(t, expected, nextState)
}

func TestThatInvalidStatusLeadsToAnAlarmWhenAlreadyRAISED(t *testing.T) {
	previousState := Alarm(RAISED)
	input := InputStatus(INVALID)
	nextState := computeAlarmStateForInputStatus(previousState, input)

	expected := RAISED

	assert.EqualValues(t, expected, nextState)
}

func TestThatValidStatusLeadsToAnAlarmWhenAlreadyRAISED(t *testing.T) {
	previousState := Alarm(RAISED)
	input := InputStatus(VALID)
	nextState := computeAlarmStateForInputStatus(previousState, input)

	expected := RAISED

	assert.EqualValues(t, expected, nextState)
}

func TestThatValidInputLeavesUnraisedAlarmUnraised(t *testing.T) {

	var previousState Alarm
	previousState = UNRAISED

	var input InputStatus
	input = VALID

	expected := UNRAISED
	result := computeAlarmStateForInputStatus(previousState, input)

	assert.EqualValues(t, expected, result)
}

func TestThatOutOfBoundsGPSInputIsInvalid(t *testing.T) {

	headingError := 4.

	bound := 3.

	expected := INVALID
	result := validateInput(bound, headingError)

	assert.EqualValues(t, expected, result)
}

func TestThatValidGPSInputIsValid(t *testing.T) {

	headingError := 2.

	bound := 3.

	expected := VALID
	result := validateInput(bound, headingError)

	assert.EqualValues(t, expected, result)
}

func TestThatUpperBoundGPSInputIsValid(t *testing.T) {

	headingError := 3.

	bound := 3.

	expected := VALID
	result := validateInput(bound, headingError)

	assert.EqualValues(t, expected, result)
}

func TestThatLowerBoundGPSInputIsValid(t *testing.T) {

	headingError := 3.

	bound := 3.

	expected := VALID
	result := validateInput(bound, headingError)

	assert.EqualValues(t, expected, result)
}
