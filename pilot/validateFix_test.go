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
* @Date:   2015-09-20 14:27:34
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:17:28
 */

package pilot

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixGPS(t *testing.T) {

	cases := []fixCase{
		fixCase{fix: NOFIX, expectedAlarm: RAISED, expectedLedStatus: true, description: "No fix "},
		fixCase{fix: FIX, expectedAlarm: UNRAISED, expectedLedStatus: false, description: "Fix"},
		fixCase{fix: DGPS_FIX, expectedAlarm: UNRAISED, expectedLedStatus: false, description: "DGPS fix"},
	}

	for _, c := range cases {
		checkNoFixCase(t, c)
	}
}

type fixCase struct {
	fix               byte
	expectedAlarm     Alarm
	expectedLedStatus bool
	description       string
}

func checkNoFixCase(t *testing.T, c fixCase) {

	fix := c.fix

	alarm, led := validateFixStatus(FixStatus(fix))

	assert.EqualValues(t, c.expectedAlarm, alarm, fmt.Sprintf("\"%s\" [alarm] case failed", c.description))

	assert.EqualValues(t, led, c.expectedLedStatus, fmt.Sprintf("\"%s\" [led] case failed", c.description))

}
