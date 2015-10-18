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
* @Author: Philippe Martinez
* @Date:   2015-10-18 19:25:31
* @Last Modified by:   Philippe Martinez
* @Last Modified time: 2015-10-18 19:25:31
 */

package conf

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/viper"
)


func TestUnMarshallingDefaultConfig(t *testing.T){
	assert.EqualValues(t, 2.23108985822891,Conf.N)
	assert.EqualValues(t, "/dev/ttyMFD1",Conf.GpsSerialPort)
	assert.EqualValues(t,25.,Conf.Bounds)
	assert.EqualValues(t,-25*(380/25),Conf.MinPIDOutputLimits)
}

func TestUnMarshallingFromFile(t *testing.T){
	viper.AddConfigPath(".")
	var localConf = loadConfiguration()
	assert.EqualValues(t, 2.13108985822891,localConf.N)
	assert.EqualValues(t, "/dev/bla",localConf.GpsSerialPort)
	assert.EqualValues(t,5.,localConf.Bounds)
	//check a default value
	assert.EqualValues(t,0.104659039843542, localConf.P)
}
