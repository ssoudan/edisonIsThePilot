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
* @Last Modified time: 2015-10-21 16:37:18
 */

package gps

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/ssoudan/edisonIsThePilot/ap100"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
	"github.com/ssoudan/edisonIsThePilot/pilot"
	"github.com/ssoudan/edisonIsThePilot/tracer"

	"github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
)

var log = logger.Log("gps")

// GPS is a driver for a serial attached NMEA GPS
type GPS struct {
	deviceName string
	baud       int

	// channels
	messagesChan chan interface{}
	headingChan  chan interface{}
	errorChan    chan interface{}
	panicChan    chan interface{}
	tracerChan   chan interface{}
}

// New creates a new GPS component
func New(deviceName string) GPS {
	return GPS{deviceName: deviceName, baud: 9600}
}

// SetMessagesChan sets the channel where the GPS messages are delivered
func (g *GPS) SetMessagesChan(c chan interface{}) {
	g.messagesChan = c
}

// SetHeadingChan sets the channel to the AP100 interface
func (g *GPS) SetHeadingChan(c chan interface{}) {
	g.headingChan = c
}

// SetErrorChan sets the channel where parsing errors are posted
func (g *GPS) SetErrorChan(c chan interface{}) {
	g.errorChan = c
}

// SetPanicChan sets the channel where panics are sent
func (g *GPS) SetPanicChan(c chan interface{}) {
	g.panicChan = c
}

// SetTracerChan sets the channel to the tracer
func (g *GPS) SetTracerChan(c chan interface{}) {
	g.tracerChan = c
}

func (g GPS) doReceiveGPSMessages() {
	c := &serial.Config{Name: g.deviceName, Baud: g.baud}

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Panicf("Failed to open serial port: %v", err)

		return
	}

	// Close the serial port when we have to leave this method
	defer s.Close()

	bufferedReader := bufio.NewReader(s)

	defer func() {
		if r := recover(); r != nil {
			log.Warning("Recovered in f", r)
		}
	}()

	for {
		str, err := bufferedReader.ReadString('\n')
		// log.Debug("[%s]", str)
		if err != nil {
			log.Error("Failed to read from serial port: %v", err)
			g.errorChan <- err
			// Exit this method to close the port, and re-open it later
			return
		}

		m, err := nmea.Parse(strings.TrimSuffix(str, "\r\n"))
		if err != nil {
			if !strings.HasSuffix(err.Error(), "not implemented") {
				g.errorChan <- err
			}
			// Here we don't return as it is a non-fatal error and the next line
			// will be better
			continue
		}

		switch t := m.(type) {
		default:
			// don't care
			// log.Debug("%+v\n", m)
		case nmea.GPGGA:
			fix, err := strconv.Atoi(t.FixQuality)
			if err != nil {
				log.Error("Failed to parse FixQuality [%s] : %v", t.FixQuality, err)

			} else {
				log.Info("[GPGGA] fixQuality: %v hdop: %s sat: %s\n", fix, t.HDOP, t.NumSatellites)
				g.messagesChan <- pilot.FixStatus(fix)
			}
		case nmea.GPRMC:
			log.Info("[GPRMC] validity: %v heading: %v[˚] speed: %v[knots] \n", t.Validity == "A", t.Course, t.Speed)
			g.messagesChan <- pilot.GPSFeedBackAction{
				Heading:   t.Course,
				Validity:  t.Validity == "A",
				Speed:     t.Speed,
				Longitude: t.Longitude,
				Latitude:  t.Latitude,
				Date:      t.Date,
				Time:      t.Time,
			}
			if t.Validity == "A" {
				g.headingChan <- ap100.NewMessage(uint16(t.Course))
				g.tracerChan <- tracer.MkAddPointMessage(types.Point{
					Latitude:  float64(t.Latitude),
					Longitude: float64(t.Longitude),
					Time:      types.JSONTime(time.Now()),
				})
			}

		}
	}
}

// Start creates an infinite go routine which will try to open the serial port to the GPS
// and parse the input to deliver sentences or errors on the respective channels.
func (g GPS) Start() {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.panicChan <- r
			}
		}()

		for {
			g.doReceiveGPSMessages()
			time.Sleep(1 * time.Second) // Cooldown in case of repeated errors
		}

	}()

}
