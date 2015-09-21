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
* @Date:   2015-09-20 09:58:18
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 23:12:03
 */

package pilot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConstroller struct {
	sp        float64
	lastValue float64
}

func (c *testConstroller) Set(sp float64) {
	log.Info("Set has been called with %v", sp)
	c.sp = sp
}

func (c *testConstroller) Update(value float64) float64 {
	log.Info("Update has been called with %v", value)
	c.lastValue = value
	return 2.
}

func TestThatTellTheWorldSendTheAlarmFirst(t *testing.T) {
	d := make(chan interface{})

	// In this test nobody reads on the dashboard channel, first message will block the sender
	// We test that alarms are still being delivered
	a := make(chan interface{})

	pilot := Pilot{
		alarm: UNRAISED}

	pilot.SetAlarmChan(a)
	pilot.SetDashboardChan(d)
	pilot.SetSteeringChan(d)

	go func() {
		pilot.tellTheWorld()
	}()

	m := <-a

	// Yup not very elegant but don't want to leak the internals of the alarm.message
	assert.EqualValues(t, "{false}", fmt.Sprintf("%v", m), "tell the world should propagates the alarm")

	pilot.alarm = RAISED

	go func() {
		pilot.tellTheWorld()
	}()

	m2 := <-a

	// Yup not very elegant but don't want to leak the internals of the alarm.message
	assert.EqualValues(t, "{true}", fmt.Sprintf("%v", m2), "tell the world should propagates the alarm")

}

func TestThatTellTheWorldPropagatesTheAlarmState(t *testing.T) {
	d := make(chan interface{})

	// This is required as we are likely to send messages on the dashboard chan too
	go func() {
		for true {
			<-d
		}
	}()

	a := make(chan interface{})

	pilot := Pilot{
		alarm: UNRAISED}

	pilot.SetAlarmChan(a)
	pilot.SetDashboardChan(d)
	pilot.SetSteeringChan(d)

	go func() {
		pilot.tellTheWorld()
	}()

	m := <-a

	// Yup not very elegant but don't want to leak the internals of the alarm.message
	assert.EqualValues(t, "{false}", fmt.Sprintf("%v", m), "tell the world should propagates the alarm")

	pilot.alarm = RAISED

	go func() {
		pilot.tellTheWorld()
	}()

	m2 := <-a

	// Yup not very elegant but don't want to leak the internals of the alarm.message
	assert.EqualValues(t, "{true}", fmt.Sprintf("%v", m2), "tell the world should propagates the alarm")

}

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
		alarmChan:     c,
		steeringChan:  c,
		inputChan:     make(chan interface{}),
		pid:           &testConstroller{}}

	// pilot.Start()
	pilot.enable()

	assert.EqualValues(t, false, pilot.headingSet, "heading need to be set during first updateFeedback")

	gpsHeadingStep1 := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeadingStep1, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

	assert.EqualValues(t, true, pilot.headingSet, "heading has been set to first gpsHeading")
	assert.EqualValues(t, gpsHeadingStep1, pilot.heading, "heading has been set to first gpsHeading")

	gpsHeading := 170.

	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

	assert.EqualValues(t, true, pilot.headingSet, "heading has been set to first gpsHeading")
	assert.EqualValues(t, gpsHeadingStep1, pilot.heading, "heading has been set to first gpsHeading")

}

func TestThatPIDControllerIsUpdatedWhenThePilotIsEnabled(t *testing.T) {

	c := make(chan interface{})

	go func() {
		for true {
			<-c
		}
	}()

	INIT_SP := -2.
	INIT_VALUE := -1.

	controller := testConstroller{sp: INIT_SP, lastValue: INIT_VALUE}

	pilot := Pilot{
		alarm:         UNRAISED,
		bound:         45,
		leds:          make(map[string]bool),
		dashboardChan: c,
		alarmChan:     c,
		steeringChan:  c,
		inputChan:     make(chan interface{}),
		pid:           &controller}

	assert.EqualValues(t, INIT_SP, controller.sp, "sp has not yet been modified")
	assert.EqualValues(t, INIT_VALUE, controller.lastValue, "error (aka PID input value) has not yet been updated")

	pilot.Start()

	assert.EqualValues(t, INIT_SP, controller.sp, "sp has not yet been modified by Start()")
	assert.EqualValues(t, INIT_VALUE, controller.lastValue, "error (aka PID input value) has not yet been updated by Start()")

	pilot.enable() // We call the internal synchronous version here

	assert.EqualValues(t, INIT_SP, controller.sp, "sp has not yet been modified by enable()")
	assert.EqualValues(t, INIT_VALUE, controller.lastValue, "error (aka PID input value) has not yet been updated by enable()")

	assert.EqualValues(t, false, pilot.headingSet, "heading need to be set during first updateFeedback")

	gpsHeadingStep1 := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeadingStep1, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

	assert.EqualValues(t, true, pilot.headingSet, "heading has been set to first gpsHeading")
	assert.EqualValues(t, gpsHeadingStep1, pilot.heading, "heading has been set to first gpsHeading")

	assert.EqualValues(t, 0., controller.sp, "sp has been set to 0 by the first updateFeedback() call after enable()")
	assert.EqualValues(t, 0., controller.lastValue, "error (aka PID input value) is zero at this point")

	gpsHeading := 170.

	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

	assert.EqualValues(t, true, pilot.headingSet, "heading has been set to first gpsHeading")
	assert.EqualValues(t, gpsHeadingStep1, pilot.heading, "heading has been set to first gpsHeading")

	assert.EqualValues(t, 0., controller.sp, "sp has been set to 0 by the first updateFeedback() call after enable()")
	assert.EqualValues(t, gpsHeading-gpsHeadingStep1, controller.lastValue, "error (aka PID input value) is now the difference between the two headings")

	pilot.disable()

	gpsHeadingStep3 := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeadingStep3, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

	assert.EqualValues(t, 0., controller.sp, "sp has been set to 0 by the first updateFeedback() call after enable()")
	assert.EqualValues(t, gpsHeading-gpsHeadingStep1, controller.lastValue, "error has not changed - since Update() has not been called cause the pilot is disabled")

}

func TestThatTimeoutRaiseTheAlarm(t *testing.T) {

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
		alarmChan:     c,
		steeringChan:  c,
		inputChan:     make(chan interface{}),
		pid:           &testConstroller{}}

	pilot.disable()

	expected := UNRAISED
	result := pilot.alarm
	assert.EqualValues(t, expected, result)

	pilot.updateAfterTimeout()

	expected = UNRAISED
	result = pilot.alarm
	assert.EqualValues(t, expected, result)

	pilot.enable()

	expected = UNRAISED
	result = pilot.alarm
	assert.EqualValues(t, expected, result)

	pilot.updateAfterTimeout()

	expected = RAISED
	result = pilot.alarm
	assert.EqualValues(t, expected, result)

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
		alarmChan:     c,
		steeringChan:  c,
		inputChan:     make(chan interface{}),
		pid:           &testConstroller{}}

	// pilot.Start()
	pilot.enable()

	gpsHeading := 180.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

	gpsHeading = 110.
	pilot.updateFeedback(GPSFeedBackAction{Heading: gpsHeading, Validity: true, Speed: MinimumSpeedInKnots * 1.1})

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
