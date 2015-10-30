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
* @Date:   2015-10-10 11:50:30
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 14:25:26
 */

package mcp4725

import (
	"bitbucket.org/gmcbay/i2c"

	"encoding/binary"

	"github.com/ssoudan/edisonIsThePilot/drivers/gpio"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("mcp4725")

const (
	// writeDacCommand is the MCP4725 command to set the output
	writeDacCommand = 0x40
	// writedacEepromCommand is the MCP4725 command to set the output and store the value in the eeprom
	writedacEepromCommand = 0x60
)

// MCP4725 is a driver for the MCP4725 i2c 12 bits DAC
type MCP4725 struct {
	bus     byte
	address byte
	i2c     *i2c.I2CBus
}

const (
	// i2c6SCL is the pin number of SCL line for i2c bus number 6
	i2c6SCL = 27
	// i2c6SDA is the pin number of SDA line for i2c bus number 6
	i2c6SDA = 28
)

// New creates a new MCP4725 driver on a i2c bus of the Edison
func New(bus byte, address byte) (*MCP4725, error) {

	switch bus {
	case 6:
		gpio.EnableI2C(i2c6SCL)
		gpio.EnableI2C(i2c6SDA)
		gpio.EnableFastI2C(6)
	default:
		log.Panic("Unknown i2c bus")
	}

	i2c, err := i2c.Bus(bus)
	if err != nil {
		return nil, err
	}

	return &MCP4725{bus: bus, address: address, i2c: i2c}, nil
}

// SetValue sets the output value of the DAC (only the 12 lower bits are used)
func (dac MCP4725) SetValue(value uint16) error {
	return dac.writeRegister16(writeDacCommand, value)
}

func (dac MCP4725) writeRegister16(reg uint8, value uint16) error {

	b := toBytes(value)
	return dac.i2c.WriteByteBlock(dac.address, reg, b)
}

func toBytes(value uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, (value<<4)&0xfff0)

	return b
}
