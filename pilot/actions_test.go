/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 21:48:32
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 21:58:34
 */

package pilot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnableDisable(t *testing.T) {

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

	assert.EqualValues(t, false, pilot.enabled, "pilot is initially not enabled")

	assert.EqualValues(t, false, pilot.headingSet, "heading still need to be set")

	pilot.enable()

	assert.EqualValues(t, true, pilot.enabled, "enable() has enabled the pilot")
	assert.EqualValues(t, false, pilot.headingSet, "heading need to be set during first updateFeedback")

	gpsHeading1 := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading1})

	assert.EqualValues(t, true, pilot.headingSet, "headingSet is set")
	assert.EqualValues(t, gpsHeading1, pilot.heading, "heading has been set to first gpsHeading")

	gpsHeading2 := 110.

	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading2})

	assert.EqualValues(t, true, pilot.headingSet, "another update does not change headingSet")
	assert.EqualValues(t, gpsHeading1, pilot.heading, "another update does not change the heading")
	assert.EqualValues(t, true, pilot.alarm, "make sure the alarm is set so we can later check disable() will reset it")

	pilot.disable()

	assert.EqualValues(t, true, pilot.headingSet, "disable does not change the headingSet")
	assert.EqualValues(t, false, pilot.enabled, "disable does update the enabled state")

	assert.EqualValues(t, false, pilot.alarm, "disable does reset the alarm")

	pilot.enable()

	assert.EqualValues(t, true, pilot.enabled, "enable() has again enabled the pilot")
	assert.EqualValues(t, false, pilot.headingSet, "heading again need to be set during next updateFeedback")

	gpsHeading3 := 90.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading3})

	assert.EqualValues(t, true, pilot.headingSet, "headingSet is set")
	assert.EqualValues(t, gpsHeading3, pilot.heading, "heading has been set to another value")

}
