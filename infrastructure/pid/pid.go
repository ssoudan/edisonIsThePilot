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
* @Date:   2015-09-24 21:35:33
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 12:15:26
 */

package pid

import (
	"time"
)

// PID is a Proportional-Integral-Derivative controller
type PID struct {
	setPoint float64

	kp float64
	ki float64
	kd float64
	n  float64

	integratorState float64
	filterState     float64
	lastUpdate      time.Time

	minOutput float64
	maxOutput float64
}

// New creates a new PID with specific parameters
func New(kp, ki, kd, n, minOutput, maxOutput float64) *PID {
	return &PID{kp: kp, ki: ki, kd: kd, n: n, minOutput: minOutput, maxOutput: maxOutput}
}

// Set sets the setpoint
func (p *PID) Set(sp float64) {
	p.setPoint = sp
}

// Update takes an error and returns the correction to be applied
func (p *PID) Update(input float64) float64 {

	// time difference
	var duration time.Duration
	if !p.lastUpdate.IsZero() {
		duration = time.Since(p.lastUpdate)
	}
	p.lastUpdate = time.Now()
	timeDifference := duration.Seconds()

	return p.updateWithDuration(input, timeDifference)
}

func (p *PID) updateWithDuration(input float64, timeDifference float64) float64 {

	// error
	u := p.setPoint - input

	// output computation
	filterCoefficient := (p.kd*u - p.filterState) * p.n
	output := (p.kp*u + p.integratorState) + filterCoefficient

	if timeDifference > 0 {
		p.integratorState += p.ki * u * timeDifference
		p.filterState += timeDifference * filterCoefficient
	}

	// saturation
	if output > p.maxOutput {
		p.integratorState -= output - p.maxOutput
		output = p.maxOutput
	} else if output < p.minOutput {
		p.integratorState += p.minOutput - output
		output = p.minOutput
	}

	return output
}

// OutputLimits returns the correction limits
func (p PID) OutputLimits() (float64, float64) {
	return p.minOutput, p.maxOutput
}
