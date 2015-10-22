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
* @Date:   2015-09-27 22:18:56
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 16:45:46
 */

package main

import (
	"math/rand"
	"time"

	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/utils"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/webserver"
	"github.com/ssoudan/edisonIsThePilot/pilot"
)

var log = logger.Log("webserver")

// Version is the version of this code -- sets at compilation time
var Version = "unknown"

type queryable struct {
	leds map[string]bool
}

type fakeTracer struct {
	points []types.Point
}

// GetPoints returns the trace as an array of types.Point
func (f fakeTracer) GetPoints() []types.Point {
	return f.points
}

// GetDashboardInfoAction returns the LEDs status in this mock implementation
func (q queryable) GetDashboardInfoAction() map[string]bool {
	return q.leds
}

type fakePilot struct {
	enabled       bool
	headingOffset float64
	setPoint      float64
	course        float64
	speed         float64
}

func (p *fakePilot) GetInfoAction() pilot.Info {
	p.speed = p.speed + (r.Float64()*1 - 0.5)
	p.course = p.course + (r.Float64()*5 - 2.5)

	pi := pilot.Info{
		Enabled:       p.enabled,
		HeadingOffset: p.headingOffset,
		SetPoint:      p.setPoint,
		Course:        p.course,
		Speed:         p.speed,
	}
	return pi
}
func (p *fakePilot) Enable() error {
	p.enabled = true
	return nil
}
func (p *fakePilot) Disable() error {
	p.enabled = false
	return nil
}
func (p *fakePilot) SetOffset(headingOffset float64) error {
	p.headingOffset = headingOffset
	return nil
}

var r = rand.New(rand.NewSource(99))

func main() {
	panicChan := make(chan interface{})
	defer func() {
		if r := recover(); r != nil {
			panicChan <- r
		}
	}()

	go func() {
		select {
		case m := <-panicChan:

			log.Fatalf("Version %v -- Received a panic error -- exiting: %v", Version, m)
		}
	}()

	ws := webserver.New(Version)
	ws.SetPanicChan(panicChan)
	ws.SetPilot(&fakePilot{})
	ws.SetTracer(&fakeTracer{points: []types.Point{{
		Latitude:  45.,
		Longitude: 5.,
		Time:      types.JSONTime(time.Now())}}})
	ws.SetDashboard(queryable{leds: map[string]bool{
		dashboard.NoGPSFix:                true,
		dashboard.InvalidGPSData:          true,
		dashboard.SpeedTooLow:             false,
		dashboard.HeadingErrorOutOfBounds: true,
		dashboard.CorrectionAtLimit:       true,
	}})
	ws.Start()

	// Wait until we receive a signal
	utils.WaitForInterrupt(func() {
		log.Info("Interrupted - exiting")
		log.Info("Exiting -- version %v", Version)
	})

}
