/*
* @Author: Sebastien Soudan
* @Date:   2015-09-18 17:13:41
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-18 17:38:09
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
				log.Fatal(err)
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
