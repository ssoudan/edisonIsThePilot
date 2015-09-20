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
* @Date:   2015-09-20 21:48:32
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:26:28
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
	bound := 45.
	pilot := Pilot{
		alarm:         UNRAISED,
		bound:         bound,
		leds:          make(map[string]bool),
		dashboardChan: c,
		inputChan:     make(chan interface{})}

	assert.EqualValues(t, false, pilot.enabled, "pilot is initially not enabled")
	assert.EqualValues(t, false, pilot.headingSet, "heading still need to be set")

	// Enable
	pilot.enable()

	assert.EqualValues(t, true, pilot.enabled, "enable() has enabled the pilot")
	assert.EqualValues(t, false, pilot.headingSet, "heading need to be set during first updateFeedback")

	gpsHeading1 := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading1})

	assert.EqualValues(t, true, pilot.headingSet, "headingSet is set")
	assert.EqualValues(t, gpsHeading1, pilot.heading, "heading has been set to first gpsHeading")

	assert.EqualValues(t, true, pilot.computeSteeringState(), "ready to go")

	gpsHeading2 := gpsHeading1 - bound - 10.
	// will cause the state to become out of bound and thus the alarm to be raised
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading2})

	assert.EqualValues(t, true, pilot.headingSet, "another update does not change headingSet")
	assert.EqualValues(t, gpsHeading1, pilot.heading, "another update does not change the heading")
	assert.EqualValues(t, true, pilot.alarm, "make sure the alarm is set so we can later check disable() will reset it")

	// Disable
	pilot.disable()

	assert.EqualValues(t, true, pilot.headingSet, "disable does not change the headingSet")
	assert.EqualValues(t, false, pilot.enabled, "disable does update the enabled state")
	assert.EqualValues(t, false, pilot.computeSteeringState(), "steering state is disable when the pilot is disabled")
	assert.EqualValues(t, false, pilot.alarm, "disable does reset the alarm")

	// Enable
	pilot.enable()

	assert.EqualValues(t, true, pilot.enabled, "enable() has again enabled the pilot")
	assert.EqualValues(t, false, pilot.headingSet, "heading again need to be set during next updateFeedback")

	gpsHeading3 := 90.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading3})

	assert.EqualValues(t, true, pilot.headingSet, "headingSet is set")
	assert.EqualValues(t, gpsHeading3, pilot.heading, "heading has been set to another value")

}
