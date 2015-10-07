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
* @Date:   2015-09-25 16:06:30
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-30 14:34:47
 */

package pid

import (
	"github.com/ssoudan/edisonIsThePilot/conf"

	"testing"

	"github.com/stretchr/testify/assert"
)

var simin = []float64{
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
	1,
}

var simout = []float64{
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	5.7477,
	4.2972,
	3.2184,
	2.416,
	1.8193,
	1.3756,
	1.0455,
	0.80008,
	0.61755,
	0.48182,
	0.38089,
	0.30584,
	0.25004,
	0.20856,
	0.17773,
	0.15482,
	0.1378,
	0.12516,
	0.11577,
	0.10881,
	0.10366,
	0.099839,
	0.097019,
	0.09494,
	0.093413,
	0.092296,
	0.091484,
	0.090899,
	0.090482,
	0.090191,
	0.089994,
	0.089865,
	0.089789,
	0.089751,
	0.089741,
	0.089753,
	0.08978,
	0.089819,
	0.089867,
	0.089921,
	0.089981,
	0.090044,
	0.090109,
	0.090177,
	0.090245,
	0.090316,
	0.090386,
	0.090458,
	0.09053,
	0.090602,
	0.090675,
	0.090747,
	0.09082,
	0.090893,
	0.090966,
	0.091039,
	0.091113,
	0.091186,
	0.091259,
	0.091332,
	0.091405,
	0.091479,
	0.091552,
	0.091625,
	0.091698,
	0.091772,
	0.091845,
	0.091918,
	0.091991,
	0.092065,
	0.092138,
	0.092211,
	0.092284,
	0.092358,
	0.092431,
	0.092504,
	0.092577,
	0.092651,
	0.092724,
	0.092797,
	0.09287,
	0.092944,
	0.093017,
	0.09309,
	0.093163,
	0.093237,
	0.09331,
	0.093383,
	0.093457,
	0.09353,
	0.093603,
	0.093676,
	0.09375,
	0.093823,
	0.093896,
	0.093969,
	0.094043,
	0.094116,
	0.094189,
	0.094262,
}

func TestThatPIDOutputMatchesTheSimulation(t *testing.T) {

	P := 0.0870095459081994   // Proportional coefficient
	I := 7.32612847120554e-05 // Integrative coefficient
	D := 22.0896577752675     // Derivative coefficient
	N := 0.25625893108953     // Derivative filter coefficient

	pidController := New(
		P,
		I,
		D,
		N,
		conf.MinPIDOutputLimits,
		conf.MaxPIDOutputLimits)

	pidController.Set(0)

	for i := 0; i < len(simin); i++ {
		in := simin[i]
		expectedOut := simout[i]

		output := pidController.updateWithDuration(-in, 1.)

		assert.InDelta(t, expectedOut, output, 1e-4, "supposed to be the same")
	}

}

func TestThatPIDOutputDoNotExceedMaxOutput(t *testing.T) {

	pidController := New(
		1,
		1,
		1,
		1,
		conf.MinPIDOutputLimits,
		conf.MaxPIDOutputLimits)

	pidController.Set(10 * conf.MaxPIDOutputLimits)

	for i := 0; i < 100; i++ {
		output := pidController.updateWithDuration(0, 1.)

		assert.True(t, output <= conf.MaxPIDOutputLimits, "output should never be larger than maxOutput")
	}

}

func TestThatPIDOutputDoNotExceedMinOutput(t *testing.T) {

	pidController := New(
		1,
		1,
		1,
		1,
		conf.MinPIDOutputLimits,
		conf.MaxPIDOutputLimits)

	pidController.Set(10 * conf.MinPIDOutputLimits)

	for i := 0; i < 100; i++ {
		output := pidController.updateWithDuration(0, 1.)
		assert.True(t, output >= conf.MinPIDOutputLimits, "output should never be larger than maxOutput")
	}

}
