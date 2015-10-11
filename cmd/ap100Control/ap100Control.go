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
* @Date:   2015-10-10 17:19:10
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-11 13:19:04
 */
package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/drivers/mcp4725"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("ap100Control")

func main() {

	bus := byte(6)
	address1 := byte(0x62) // or 0x63
	address2 := byte(0x63)

	mcp, err := mcp4725.New(bus, address1)
	if err != nil {
		log.Panic("%v", err)
	}

	mcp2, err := mcp4725.New(bus, address2)
	if err != nil {
		log.Panic("%v", err)
	}

	for i := uint16(0); i < 0xfff; i++ {
		err = mcp.SetValue(i)
		if err != nil {
			log.Info("mcp1: %v", err)
		}

		err = mcp2.SetValue((i + 0x800) % 0xfff)
		if err != nil {
			log.Info("mcp2: %v", err)
		}
		time.Sleep(10 * time.Microsecond)
	}

}
