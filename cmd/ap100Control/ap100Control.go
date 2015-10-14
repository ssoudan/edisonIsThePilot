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
* @Last Modified time: 2015-10-13 17:24:09
 */
package main

import (
	"time"

	"github.com/ssoudan/edisonIsThePilot/conf"
	"github.com/ssoudan/edisonIsThePilot/drivers/sincos"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("ap100Control")

func main() {

	compass := sincos.New(conf.I2CBus, conf.SinAddress, conf.CosAddress)

	for i := uint16(0); i <= 360; i++ {
		compass.UpdateCourse(i)
		time.Sleep(200 * time.Millisecond)
	}

}
