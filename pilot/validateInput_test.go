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
* @Date:   2015-09-20 12:09:28
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:17:27
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
