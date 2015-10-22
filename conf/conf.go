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
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-22 15:13:00
 */

package conf

import (
	"github.com/spf13/viper"

	"github.com/ssoudan/edisonIsThePilot/dashboard"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
)

var log = logger.Log("conf")

// MessagePin is the mapping of pin to a message
type MessagePin struct {
	Message string
	Pin     byte
}

// MessageToPin contains the mapping of the dashboard messages to the LED pin
var MessageToPin = []MessagePin{
	{dashboard.NoGPSFix, 43},                // J19 - pin 11
	{dashboard.InvalidGPSData, 48},          // J19 - pin 6
	{dashboard.SpeedTooLow, 40},             // J19 - pin 10
	{dashboard.HeadingErrorOutOfBounds, 82}, // J19 - pin 13
	{dashboard.CorrectionAtLimit, 83},       // J19 - pin 14
}

const (
	// AlarmGpioPin is the pin where the alarm is connected
	AlarmGpioPin = 183 // J18 - pin 8
	// AlarmGpioPWM is the pwm where the alarm is connected
	AlarmGpioPWM = 3

	// SwitchGpioPin is the pin where the autopilot switch is connected
	SwitchGpioPin = 46 // J19 - pin 5

	// MotorDirPin is the direction pin for the motor
	MotorDirPin = 165 // J18 - pin 2
	// MotorSleepPin is the sleep pin for the motor
	MotorSleepPin = 12 // J18 - pin 7
	// MotorStepPin is the step pin for the motor
	MotorStepPin = 182 // J17 - pin 1
	// MotorStepPwm is the pwm where the step pin of the motor is connected
	MotorStepPwm = 2

	// I2CBus of the Sin/Cos interface
	I2CBus = 6
	// SinAddress is the i2c address of the Sine DAC (MCP4725)
	SinAddress = 0x62
	// CosAddress is the i2c address of the Cosine DAC (MCP4725)
	CosAddress = 0x63
)

// Configuration is the type of the configuration loaded from the config file
type Configuration struct {
	Bounds                         float64 // error bound in degree
	SteeringReductionRatio         float64 // reduction ratio between the motor and the steering wheel
	MaxPIDOutputLimits             float64 // maximum pid output value (in degree)
	MinPIDOutputLimits             float64 // minimum pid output value (in degree)
	P                              float64 // Proportional coefficient
	I                              float64 // Integrative coefficient
	D                              float64 // Derivative coefficient
	N                              float64 // Derivative filter coefficient
	GpsSerialPort                  string
	NoInputMessageTimeoutInSeconds int64
	MinimumSpeedInKnots            float64
	TraceSize                      uint32
}

func setDefaultValues() {
	viper.SetDefault("Bounds", 25.)
	viper.SetDefault("SteeringReductionRatio", 380/25)
	viper.SetDefault("P", 0.104659039843542)
	viper.SetDefault("I", 8.06799673280568e-05)
	viper.SetDefault("D", 27.8353089535829)
	viper.SetDefault("N", 2.23108985822891)
	viper.SetDefault("GpsSerialPort", "/dev/ttyMFD1")
	viper.SetDefault("NoInputMessageTimeoutInSeconds", 10)
	viper.SetDefault("MinimumSpeedInKnots", 3)
	viper.SetDefault("TraceSize", 500)
	viper.SetDefault("MinPIDOutputLimits", -380.)
	viper.SetDefault("MaxPIDOutputLimits", 380.)
}

func loadConfiguration() Configuration {
	var conf Configuration

	viper.SetConfigType("properties")
	viper.SetConfigName("edisonIsThePilot") // name of config file (without extension)
	viper.AddConfigPath("/etc")             // path to look for the config file in
	setDefaultValues()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Warning("Unable to read the configuration file /etc/edisonIsThePilot.properties. Using default values.")

	}
	err = viper.Unmarshal(&conf)
	if err != nil { // Handle errors reading the config file
		log.Panic("Unable to load configuration. Check /etc/edisonIsThePilot.properties file")
	}

	log.Info("Configuration is: %#v", conf)
	return conf
}

// Conf contains the configuration parameters loaded at the initialization from the edisonIsThePilot.properties or defaults values
var Conf = loadConfiguration()
