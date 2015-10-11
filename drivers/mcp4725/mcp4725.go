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
* @Last Modified time: 2015-10-11 13:15:31
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
	WRITEDAC        = 0x40
	WRITEDAC_EEPROM = 0x60
)

type MCP4725 struct {
	bus     byte
	address byte
	i2c     *i2c.I2CBus
}

const (
	I2C_6_SCL = 27
	I2C_6_SDA = 28
)

func New(bus byte, address byte) (*MCP4725, error) {

	switch bus {
	case 6:
		gpio.EnableI2C(I2C_6_SCL)
		gpio.EnableI2C(I2C_6_SDA)
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

func (dac MCP4725) SetValue(value uint16) error {
	return dac.writeRegister16(WRITEDAC, value)
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
