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
* @Last Modified time: 2015-09-29 10:26:24
 */

package pilot

import (
	"github.com/adrianmo/go-nmea"
)

type GPSFeedBackAction struct {
	Heading   float64
	Validity  bool
	Speed     float64
	Latitude  nmea.LatLong
	Longitude nmea.LatLong
	Date      string
	Time      string
}

type EnableAction struct {
}

type DisableAction struct {
}

// Enable the autopilot
func (p *Pilot) Enable() error {
	p.inputChan <- EnableAction{}
	return nil
}

// Disable the autopilot
func (p *Pilot) Disable() error {
	p.inputChan <- DisableAction{}
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

func (p Pilot) Shutdown() {
	p.shutdownChan <- 1
	<-p.shutdownChan
}

func (p *Pilot) shutdown() {
	p.enabled = false
	p.headingSet = false
	close(p.shutdownChan)
}
