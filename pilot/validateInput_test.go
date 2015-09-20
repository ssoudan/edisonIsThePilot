/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 12:09:28
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 12:09:57
 */

package pilot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
