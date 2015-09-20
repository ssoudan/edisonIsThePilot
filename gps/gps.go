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
* @Date:   2015-09-18 17:13:41
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 21:45:01
 */

package gps

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/pilot"

	"github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
)

var log = logger.Log("gps")

type GPS struct {
	deviceName string
	baud       int

	// channels
	messagesChan chan interface{}
	errorChan    chan interface{}
}

// New creates a new GPS component
func New(deviceName string) GPS {
	return GPS{deviceName: deviceName, baud: 9600}
}

func (g *GPS) SetMessagesChan(c chan interface{}) {
	g.messagesChan = c
}

func (g *GPS) SetErrorChan(c chan interface{}) {
	g.errorChan = c
}

func (g GPS) doReceiveGPSMessages() {
	c := &serial.Config{Name: g.deviceName, Baud: g.baud}

	var err error
	var s *serial.Port

	//
	for s == nil {
		s, err = serial.OpenPort(c)
		if err != nil {
			log.Error("Failed to open serial port: %v", err)
			g.errorChan <- err

			time.Sleep(1 * time.Second)
		}
	}

	bufferedReader := bufio.NewReader(s)
	// Close the serial port when we have to leave this method
	defer s.Close()

	for true {
		str, err := bufferedReader.ReadString('\n')
		if err != nil {
			log.Error("Failed to read from serial port: %v", err)
			g.errorChan <- err
			// Exit this method to close the port, and re-open it later
			return
		}

		m, err := nmea.Parse(strings.TrimSuffix(str, "\r\n"))
		if err != nil {
			g.errorChan <- err
			// Here we don't return as it is a non-fatal error and the next line
			// will be better
		} else {

			switch t := m.(type) {
			default:
				// don't care
				// log.Debug("%+v\n", m)

			case nmea.GPGGA:
				fix, err := strconv.Atoi(t.FixQuality)
				if err != nil {
					log.Error("Failed to parse FixQuality [%s] : %v", t.FixQuality, err)

				} else {
					log.Info("[GPGGA] fixQuality: %v \n", fix)
					g.messagesChan <- pilot.FixStatus(fix)
				}

			case nmea.GPRMC:
				log.Info("[GPRMC] validity: %v heading: %v[˚] speed: %v[knots] \n", t.Validity == "A", t.Course, t.Speed)
				g.messagesChan <- pilot.GPSFeedBackAction{t.Course, t.Validity == "A", t.Speed}
			}
		}
	}
}

// Start creates an infinite go routine which will try to open the serial port to the GPS
// and parse the input to delive sentences or errors on the respective channels.
func (g GPS) Start() {

	go func() {

		for true {
			g.doReceiveGPSMessages()
			time.Sleep(1 * time.Second)
		}

	}()

}
