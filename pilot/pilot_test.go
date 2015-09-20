/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 09:58:18
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 12:09:46
 */

package pilot

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatOutOfBoundsGPSInputRaisesAnAlarm(t *testing.T) {

	pilot := Pilot{alarm: UNRAISED, heading: 112., bound: 45}
	gpsHeading := 180.
	pilot.UpdateInput(gpsHeading)

	expected := RAISED
	result := pilot.alarm

	assert.EqualValues(t, expected, result)
}

type headingCase struct {
	heading     float64
	gpsHeading  float64
	expected    float64
	description string
}

func checkHeadingCase(t *testing.T, c headingCase) {

	pilot := Pilot{heading: c.heading}

	expected := c.expected
	result := pilot.computeHeadingError(c.gpsHeading)

	assert.EqualValues(t, expected, result, fmt.Sprintf("\"%s\" case failed", c.description))
}

func TestHeadingErrorComputation(t *testing.T) {
	cases := []headingCase{
		headingCase{gpsHeading: 140., expected: 140., description: "heading not explicit set"},
		headingCase{heading: 0., gpsHeading: 0., expected: 0., description: "heading equals gpsHeading"},
		headingCase{heading: 1., gpsHeading: 2., expected: 1., description: "heading NE quadrant, gpsHeading NE quadrant"},
		headingCase{heading: 1., gpsHeading: 339., expected: -22., description: "heading NE quadrant, gpsHeading NW quadrant"},
		headingCase{heading: 349., gpsHeading: 359., expected: 10., description: "heading NW quadrant, gpsHeading NW quadrant"},
		headingCase{heading: 349., gpsHeading: 10., expected: 21., description: "heading NW quadrant, gpsHeading NE quadrant"},
		headingCase{heading: 1., gpsHeading: 181., expected: 180., description: "limit case"},
		headingCase{heading: 181., gpsHeading: 1., expected: 180., description: "limit case 2"},
		headingCase{heading: 130., gpsHeading: 350., expected: -140., description: "limit case 3"},
	}

	for _, c := range cases {
		checkHeadingCase(t, c)
	}
}
