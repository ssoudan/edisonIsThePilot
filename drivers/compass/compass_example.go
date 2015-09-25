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
* @Date:   2015-09-20 09:58:02
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-23 07:32:49
 */

package compass

import (
	"github.com/ssoudan/edisonIsThePilot/drivers/compass/hmc"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("compass")

func main() {
	////////////////////////////////////////
	// HMC5883 stuffs
	////////////////////////////////////////
	compass := hmc.New(6)
	for !compass.Begin() {

	}

	// Set measurement range
	compass.SetRange(hmc.HMC5883L_RANGE_1_3GA)

	// Set measurement mode
	compass.SetMeasurementMode(hmc.HMC5883L_CONTINOUS)

	// Set data rate
	compass.SetDataRate(hmc.HMC5883L_DATARATE_3HZ)

	// Set number of samples averaged
	compass.SetSamples(hmc.HMC5883L_SAMPLES_8)

	// Set calibration offset. See HMC5883L_calibration.ino
	compass.SetOffset(-82, 72)

	mag, err := compass.ReadNormalize()
	if err == nil {
		log.Info("Compass reading is %v", mag)
	}
}
