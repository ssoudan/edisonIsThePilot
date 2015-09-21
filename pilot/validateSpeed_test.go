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
* @Date:   2015-09-21 22:49:13
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 23:08:32
 */

package pilot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckSpeedErrorRaisesAnAlarmWhenTooSlow(t *testing.T) {

	input := MinimumSpeedInKnots * 0.9
	alarm := checkSpeedError(input)

	expected := RAISED

	assert.EqualValues(t, expected, alarm)
}

func TestCheckSpeedErrorDoNotRaiseAnAlarmWhenLargeEnougt(t *testing.T) {

	input := MinimumSpeedInKnots * 1.1
	alarm := checkSpeedError(input)

	expected := UNRAISED

	assert.EqualValues(t, expected, alarm)
}
