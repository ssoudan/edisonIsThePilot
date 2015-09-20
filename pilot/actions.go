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
* @Last Modified time: 2015-09-20 22:17:07
 */

package pilot

type GPSFeedBackAction struct {
	Heading  float64
	Validity bool
	Speed    float64
}

type EnableAction struct {
}

type DisableAction struct {
}

// Enable the autopilot
func (p *Pilot) Enable() {
	p.inputChan <- EnableAction{}
}

// Disable the autopilot
func (p *Pilot) Disable() {
	p.inputChan <- DisableAction{}
}

func (p *Pilot) enable() {
	p.enabled = true
	p.headingSet = false
}

func (p *Pilot) disable() {
	p.enabled = false
	p.alarm = UNRAISED
}
