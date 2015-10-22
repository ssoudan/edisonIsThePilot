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
* @Date:   2015-09-20 21:45:21
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 12:48:03
 */

package pilot

import (
	"github.com/adrianmo/go-nmea"
)

// GPSFeedBackAction is the message provided by the GPS component
type GPSFeedBackAction struct {
	Heading   float64
	Validity  bool
	Speed     float64
	Latitude  nmea.LatLong
	Longitude nmea.LatLong
	Date      string
	Time      string
}

type enableAction struct {
}

type disableAction struct {
}

type setOffsetAction struct {
	headingOffset float64
}

// Info contains the Pilot state information as used by the Webserver for example
type Info struct {
	Course        float64
	SetPoint      float64
	HeadingOffset float64
	Speed         float64
	Enabled       bool
}

type getInfoAction struct {
	backChannel chan Info
}

// GetInfoAction returns the current Info
func (p *Pilot) GetInfoAction() Info {
	c := make(chan Info)
	defer close(c)
	p.inputChan <- getInfoAction{backChannel: c}
	return <-c
}

func (p *Pilot) getInfoAction(c chan Info) {
	i := Info{
		Course:        p.course,
		SetPoint:      p.heading,
		HeadingOffset: p.headingOffset,
		Enabled:       p.enabled,
		Speed:         p.speed,
	}
	c <- i
}

// SetOffset changes the heading offset
func (p *Pilot) SetOffset(headingOffset float64) error {
	p.inputChan <- setOffsetAction{headingOffset: headingOffset}
	return nil
}

func (p *Pilot) setOffset(headingOffset float64) {
	p.headingOffset = headingOffset
}

// Enable the autopilot
func (p *Pilot) Enable() error {
	p.inputChan <- enableAction{}
	return nil
}

// Disable the autopilot
func (p *Pilot) Disable() error {
	p.inputChan <- disableAction{}
	return nil
}

func (p *Pilot) enable() {
	p.enabled = true
	p.headingSet = false
}

func (p *Pilot) disable() {
	p.enabled = false
	p.alarm = UNRAISED
}

// Shutdown the event loop of the Pilot
func (p Pilot) Shutdown() {
	p.shutdownChan <- 1
	<-p.shutdownChan
}

func (p *Pilot) shutdown() {
	p.enabled = false
	p.headingSet = false
	close(p.shutdownChan)
}
