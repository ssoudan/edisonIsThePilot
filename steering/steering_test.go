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
* @Date:   2015-09-21 18:47:13
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-24 17:25:31
 */

package steering

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThatRotationDirectionIsClockwise(t *testing.T) {
	clockwise, _, _ := rotationInDegreeToMove(12)

	assert.EqualValues(t, true, clockwise, "supposed to be clockwise")

	clockwise, _, _ = rotationInDegreeToMove(-12)

	assert.EqualValues(t, false, clockwise, "supposed to be anticlockwise")

}

func TestThatSpeedIsConstant(t *testing.T) {
	_, speed1, _ := rotationInDegreeToMove(12)

	_, speed2, _ := rotationInDegreeToMove(140)

	assert.EqualValues(t, speed1, speed2, "speed is contant")

}

func TestThatNothingMovesAtNullSpeed(t *testing.T) {
	clockwise, _, duration := rotationInDegreeToMove(0)

	assert.EqualValues(t, false, clockwise, "supposed to be clockwise")

	assert.EqualValues(t, time.Duration(0), duration, "supposed to be clockwise")

}
