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
* @Date:   2015-10-12 19:20:55
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-13 17:23:52
 */

package sincos

import (
	"math"

	"github.com/ssoudan/edisonIsThePilot/drivers/mcp4725"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("sincos")

// ToSinCos returns the Sine/Cosine values as expected by the Robertson autopilots (+/- 2V centered on 2.5V)
func ToSinCos(theta uint16) (uint16, uint16) {
	s, c := math.Sincos(float64(theta) * math.Pi / 180)

	// MCP4725 is a 12 bits DAC between 0 and Vcc(=5V)
	sf := float64(0xfff) * (0.5 + 0.4*s)
	cf := float64(0xfff) * (0.5 + 0.4*c)

	return uint16(sf), uint16(cf)
}

type SinCosInterface struct {
	sin *mcp4725.MCP4725
	cos *mcp4725.MCP4725
}

// New creates a new SinCosInterface interface (through 2 MCP4725 on an i2c bus)
func New(bus, sinAddr, cosAddr byte) *SinCosInterface {

	sin, err := mcp4725.New(bus, sinAddr)
	if err != nil {
		log.Panic(err)
	}

	cos, err := mcp4725.New(bus, cosAddr)
	if err != nil {
		log.Panic(err)
	}

	return &SinCosInterface{sin: sin, cos: cos}
}

// UpdateCourse sets the Sin/Cos outputs to the values that correspond to the provided course (in degree)
func (ap SinCosInterface) UpdateCourse(theta uint16) error {
	ss, cc := ToSinCos(theta)
	err := ap.sin.SetValue(ss)
	if err != nil {
		return err
	}
	err = ap.cos.SetValue(cc)
	if err != nil {
		return err
	}
	return nil
}
