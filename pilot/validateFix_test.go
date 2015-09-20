/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 14:27:34
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:07:47
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
