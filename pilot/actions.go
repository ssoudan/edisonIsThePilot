/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 21:45:21
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 21:59:17
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
