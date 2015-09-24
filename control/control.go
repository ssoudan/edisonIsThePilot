/*
Copyright 2015 Sebastien Soudan

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy
of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliec. See the
License for the specific language governing permissions and limitations
under the License.
*/

/*
* @Author: Sebastien Soudan
* @Date:   2015-09-22 11:55:49
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-24 13:12:42
 */

package control

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"

	"time"
)

var log = logger.Log("control")

type Control struct {
	controlHandler Readable
	pilot          Enablable
	stateEnable    bool

	// channels
	pilotChan    chan interface{}
	shutdownChan chan interface{}
}

type Enablable interface {
	Enable() error
	Disable() error
}

type Readable interface {
	Value() (bool, error)
}

func New(controlHandler Readable, pilot Enablable) *Control {
	return &Control{controlHandler: controlHandler, pilot: pilot, shutdownChan: make(chan interface{})}
}

func (c *Control) updateControlState() error {

	control := c.controlHandler
	state, err := control.Value()
	if err != nil {
		log.Error("Failed to read switch value: %v", err)
		return err
	}

	if state != c.stateEnable {
		if state {
			log.Warning("Enabling the pilot")
			err = c.pilot.Enable()

		} else {
			log.Warning("Disabling the pilot")
			err = c.pilot.Disable()
		}

		if err == nil {
			c.stateEnable = state
		}
	}

	return err

}

// Shutdown sets all the state to down and notify the handlers
func (c Control) Shutdown() {
	c.shutdownChan <- 1
	<-c.shutdownChan
}

func (c Control) shutdown() {
	// Nothing
	close(c.shutdownChan)
}

// Start the event loop of the Control component
func (c Control) Start() {

	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				err := c.updateControlState()
				if err != nil {
					log.Error("Error while updating control: %v", err)
				}
			case <-c.shutdownChan:
				c.shutdown()
				return
			}

		}
	}()

}
