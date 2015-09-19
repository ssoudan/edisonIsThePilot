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
* @Last Modified time: 2015-09-19 11:12:33
 */

package gps

import (
	"bufio"
	"strings"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"

	"github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
)

var log = logger.Log("gps")

type GPS struct {
	deviceName string
	baud       int
}

func New(deviceName string) GPS {
	return GPS{deviceName: deviceName, baud: 9600}
}

func (g GPS) Stream() (chan nmea.SentenceI, chan error) {
	messagesChan := make(chan nmea.SentenceI)
	errorChan := make(chan error)

	c := &serial.Config{Name: g.deviceName, Baud: g.baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		bufferedReader := bufio.NewReader(s)
		defer s.Close()

		for true {
			str, err := bufferedReader.ReadString('\n')
			if err != nil {
				log.Fatal(err) // TODO(ssoudan) we cannot fail like this -- need to recover on a new line
			}

			m, err := nmea.Parse(strings.TrimSuffix(str, "\r\n"))
			if err == nil {
				messagesChan <- m
			} else {
				errorChan <- err
			}
		}

	}()

	return messagesChan, errorChan
}
