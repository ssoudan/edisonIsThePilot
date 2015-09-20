/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 09:58:18
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 21:44:39
 */

package pilot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThatHeadingIsSetWithFirstGPSHeadingAfterItHasBeenEnabled(t *testing.T) {

	c := make(chan interface{})

	go func() {
		for true {
			<-c
		}
	}()

	pilot := Pilot{
		alarm:         UNRAISED,
		bound:         45,
		leds:          make(map[string]bool),
		dashboardChan: c,
		inputChan:     make(chan interface{})}

	// pilot.Start()
	pilot.enable()

	assert.EqualValues(t, false, pilot.headingSet, "heading need to be set during first updateFeedback")

	gpsHeadingStep1 := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeadingStep1})

	assert.EqualValues(t, true, pilot.headingSet, "heading has been set to first gpsHeading")
	assert.EqualValues(t, gpsHeadingStep1, pilot.heading, "heading has been set to first gpsHeading")

	gpsHeading := 170.

	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading})

	assert.EqualValues(t, true, pilot.headingSet, "heading has been set to first gpsHeading")
	assert.EqualValues(t, gpsHeadingStep1, pilot.heading, "heading has been set to first gpsHeading")

}

func TestThatOutOfBoundsGPSInputRaisesAnAlarm(t *testing.T) {

	c := make(chan interface{})

	go func() {
		for true {
			<-c
		}
	}()

	pilot := Pilot{
		alarm:         UNRAISED,
		bound:         45,
		leds:          make(map[string]bool),
		dashboardChan: c,
		inputChan:     make(chan interface{})}

	// pilot.Start()
	pilot.enable()

	gpsHeading := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading})

	gpsHeading = 110.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading})

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
	result := computeHeadingError(pilot.heading, c.gpsHeading)

	assert.EqualValues(t, expected, result, fmt.Sprintf("\"%s\" case failed", c.description))
}

func TestHeadingErrorComputation(t *testing.T) {
	cases := []headingCase{
		headingCase{ /* heading not set*/ gpsHeading: 140., expected: 140., description: "heading not explicit set"},
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
