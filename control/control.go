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
* @Date:   2015-09-22 11:55:49
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 12:53:12
 */

package control

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"

	"time"
)

var log = logger.Log("control")

// Control is a component that monitors change on a types.Readable and Enable or Disable its target
type Control struct {
	controlHandler types.Readable
	target         types.Enablable
	stateEnable    bool

	// channels
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

// New creates a new Control component
func New(controlHandler types.Readable, target types.Enablable) *Control {
	return &Control{controlHandler: controlHandler, target: target, shutdownChan: make(chan interface{})}
}

// SetPanicChan sets the channel where panics will be sent
func (c *Control) SetPanicChan(p chan interface{}) {
	c.panicChan = p
}

func (c *Control) updateControlState() error {

	control := c.controlHandler
	state, err := control.Value()
	if err != nil {
		log.Panicf("Failed to read switch value: %v", err)
		return err
	}

	if state != c.stateEnable {
		if state {
			log.Warning("Enabling the target")
			err = c.target.Enable()

		} else {
			log.Warning("Disabling the target")
			err = c.target.Disable()
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
		defer func() {
			if r := recover(); r != nil {
				c.panicChan <- r
			}
		}()

		for {
			select {
			case <-time.After(100 * time.Millisecond):
				err := c.updateControlState()
				if err != nil {
					log.Panicf("Error while updating control: %v", err)
				}
			case <-c.shutdownChan:
				c.shutdown()
				return
			}

		}
	}()

}
