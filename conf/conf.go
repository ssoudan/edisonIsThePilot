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
* @Date:   2015-09-22 13:18:01
* @Last Modified by:   Philippe Martinez
* @Last Modified time: 2015-10-18 19:33:58
 */

package conf

import (
	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/spf13/viper"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"

)
var log = logger.Log("conf")


type MessagePin struct {
	Message string
	Pin     byte
}

var MessageToPin = []MessagePin{
	{dashboard.NoGPSFix, 43},                // J19 - pin 11
	{dashboard.InvalidGPSData, 48},          // J19 - pin 6
	{dashboard.SpeedTooLow, 40},             // J19 - pin 10
	{dashboard.HeadingErrorOutOfBounds, 82}, // J19 - pin 13
	{dashboard.CorrectionAtLimit, 83},       // J19 - pin 14
}

const (
	AlarmGpioPin  = 183  // J18 - pin 8
	AlarmGpioPWM  = 3    // corresponding pwm number
	MotorDirPin   = 165  // J18 - pin 2
	MotorSleepPin = 12   // J18 - pin 7
	MotorStepPin  = 182  // J17 - pin 1
	MotorStepPwm  = 2    // corresponding pwm number
	SwitchGpioPin = 46   // J19 - pin 5
	I2CBus        = 6    // I2C bus where the Sin/Cos DAC are
	SinAddress    = 0x62 // Address on I2CBus of the MCP4725 used for the Sine
	CosAddress    = 0x63 // Address on I2CBus of the MCP4725 used for the Cosine
)

type Configuration struct{
	Bounds                         float64// error bound in degree
	SteeringReductionRatio         float64// reduction ratio between the motor and the steering wheel
	MaxPIDOutputLimits             float64// maximum pid output value (in degree)
	MinPIDOutputLimits             float64// minimum pid output value (in degree)
	P                              float64// Proportional coefficient
	I                              float64// Integrative coefficient
	D                              float64// Derivative coefficient
	N                              float64// Derivative filter coefficient
	GpsSerialPort                  string
	NoInputMessageTimeoutInSeconds int64
	MinimumSpeedInKnots            float64
}


func postLoadConfiguration(configuration Configuration) Configuration {

	//set relative values  
	configuration.MaxPIDOutputLimits = configuration.Bounds * configuration.SteeringReductionRatio
	configuration.MinPIDOutputLimits = -(configuration.Bounds) * configuration.SteeringReductionRatio
	
	return configuration
}

func setDefaultValues() {
	viper.SetDefault("Bounds",25.)
	viper.SetDefault("SteeringReductionRatio",380/25)
	viper.SetDefault("P",0.104659039843542)
	viper.SetDefault("I",8.06799673280568e-05)
	viper.SetDefault("D",27.8353089535829)
	viper.SetDefault("N",2.23108985822891)
	viper.SetDefault("GpsSerialPort","/dev/ttyMFD1")
	viper.SetDefault("NoInputMessageTimeoutInSeconds",10)
	viper.SetDefault("MinimumSpeedInKnots",3)
}

func loadConfiguration() Configuration {
	var conf Configuration
	
	viper.SetConfigType("properties") 
	viper.SetConfigName("edisonIsThePilot") // name of config file (without extension)
	viper.AddConfigPath("/etc")   // path to look for the config file in	
	setDefaultValues()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
	     log.Warning("Unable to read the configuration file (looking for edisonIsThePilot.properties in /etc/ path). Using default values")
	    
	}
	err = viper.Unmarshal(&conf) 
	if err != nil { // Handle errors reading the config file
	    log.Panic("Unable to load configuration. Check edisonIsThePilot.properties file")
	}
	
	conf = postLoadConfiguration(conf)
	log.Info("Configuration is: %#v", conf)
	return conf
}	
// Conf contains the configuration parameters loaded at the initialization from the edisonIsThePilot.properties or defaults values
var Conf = loadConfiguration() 