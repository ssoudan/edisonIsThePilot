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
* @Date:   2015-09-21 15:46:55
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-21 15:51:42
 */

package alarm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testHandler struct {
	state bool
}

func (h *testHandler) Enable() error {
	h.state = true
	return nil
}
func (h *testHandler) Disable() error {
	h.state = false
	return nil
}

func TestAlarmChangesState(t *testing.T) {

	alarmHandler := &testHandler{}
	alarm := New(alarmHandler)

	c := make(chan interface{})

	alarm.SetInputChan(c)

	alarm.Start()

	assert.EqualValues(t, false, alarmHandler.state, "state is false by default")

	m := NewMessage(true)

	alarm.processMessage(m.(message))

	assert.EqualValues(t, true, alarmHandler.state, "previous message should have set it to true")

	m = NewMessage(false)

	alarm.processMessage(m.(message))

	assert.EqualValues(t, false, alarmHandler.state, "should be false")

	m = NewMessage(false)

	alarm.processMessage(m.(message))

	assert.EqualValues(t, false, alarmHandler.state, "should be false again")

}

func TestAlarmWithChannel(t *testing.T) {
	alarmHandler := &testHandler{}
	alarm := New(alarmHandler)

	c := make(chan interface{})

	alarm.SetInputChan(c)

	alarm.Start()

	assert.EqualValues(t, false, alarmHandler.state, "state is false by default")

	m := NewMessage(true)

	c <- m

	// Have to wait for the update to be propagated
	for !alarmHandler.state {
		time.Sleep(1 * time.Second)
	}

	assert.EqualValues(t, true, alarmHandler.state, "previous message should have set it to true")

}

func TestAlarmShutdown(t *testing.T) {
	alarmHandler := &testHandler{}
	alarm := New(alarmHandler)

	c := make(chan interface{})

	alarm.SetInputChan(c)

	alarm.Start()

	assert.EqualValues(t, false, alarmHandler.state, "state is false by default")

	m := NewMessage(true)

	c <- m

	// Have to wait for the update to be propagated
	for !alarmHandler.state {
		log.Debug("alarmHandler.state = %v", alarmHandler.state)
		time.Sleep(1 * time.Second)
	}

	assert.EqualValues(t, true, alarmHandler.state, "previous message should have set it to true")

	alarm.Shutdown()

	for alarmHandler.state {
		log.Debug("alarmHandler.state = %v", alarmHandler.state)
		time.Sleep(1 * time.Second)
	}

	assert.EqualValues(t, false, alarmHandler.state, "shutdown effect")
}
