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
* @Date:   2015-09-20 23:11:55
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 23:44:25
 */

package dashboard

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

func TestDashboardChangesLEDStates(t *testing.T) {

	dashboard := New()

	aaaHandler := testHandler{}
	dashboard.RegisterMessageHandler("AAA", &aaaHandler)

	bbbHandler := testHandler{}
	dashboard.RegisterMessageHandler("BBB", &bbbHandler)

	c := make(chan interface{})

	dashboard.SetInputChan(c)

	dashboard.Start()

	assert.EqualValues(t, false, aaaHandler.state, "state is false by default")
	assert.EqualValues(t, false, bbbHandler.state, "state is false by default (2)")

	m := NewMessage(map[string]bool{"AAA": true})

	dashboard.processMessage(m.(message))

	assert.EqualValues(t, true, aaaHandler.state, "previous message should have set it to true")
	assert.EqualValues(t, false, bbbHandler.state, "state should not have changed for BBB")

	m = NewMessage(map[string]bool{"BBB": true})

	dashboard.processMessage(m.(message))

	assert.EqualValues(t, true, aaaHandler.state, "state should not have changed for AAA")
	assert.EqualValues(t, true, bbbHandler.state, "BBB should have changed to true")

	m = NewMessage(map[string]bool{"BBB": false})

	dashboard.processMessage(m.(message))

	assert.EqualValues(t, true, aaaHandler.state, "still true")
	assert.EqualValues(t, false, bbbHandler.state, "and back to false")

	m = NewMessage(map[string]bool{"BBB": false, "AAA": false})

	dashboard.processMessage(m.(message))

	assert.EqualValues(t, false, aaaHandler.state, "false now")
	assert.EqualValues(t, false, bbbHandler.state, "false too")

}

func TestDashboardWithChannel(t *testing.T) {
	dashboard := New()

	aaaHandler := testHandler{}
	dashboard.RegisterMessageHandler("AAA", &aaaHandler)

	bbbHandler := testHandler{}
	dashboard.RegisterMessageHandler("BBB", &bbbHandler)

	c := make(chan interface{})

	dashboard.SetInputChan(c)

	dashboard.Start()

	assert.EqualValues(t, false, aaaHandler.state, "state is false by default")
	assert.EqualValues(t, false, bbbHandler.state, "state is false by default (2)")

	m := NewMessage(map[string]bool{"AAA": true})

	c <- m

	// Have to wait for the update to be propagated
	for !aaaHandler.state {
		time.Sleep(1 * time.Second)
	}

	assert.EqualValues(t, true, aaaHandler.state, "previous message should have set it to true")
	assert.EqualValues(t, false, bbbHandler.state, "state should not have changed for BBB")

}

func TestDashboardShutdown(t *testing.T) {
	dashboard := New()

	aaaHandler := testHandler{}
	dashboard.RegisterMessageHandler("AAA", &aaaHandler)

	bbbHandler := testHandler{}
	dashboard.RegisterMessageHandler("BBB", &bbbHandler)

	c := make(chan interface{})

	dashboard.SetInputChan(c)

	dashboard.Start()

	assert.EqualValues(t, false, aaaHandler.state, "state is false by default")
	assert.EqualValues(t, false, bbbHandler.state, "state is false by default (2)")

	m := NewMessage(map[string]bool{
		"BBB": true,
		"AAA": true,
	})

	c <- m

	// Have to wait for the update to be propagated
	for !(aaaHandler.state && bbbHandler.state) {
		log.Debug("aaaHandler.state = %v", aaaHandler.state)
		log.Debug("bbbHandler.state = %v", bbbHandler.state)
		time.Sleep(1 * time.Second)
	}

	assert.EqualValues(t, true, aaaHandler.state, "previous message should have set it to true")
	assert.EqualValues(t, true, bbbHandler.state, "state should have changed for BBB too")

	dashboard.Shutdown()

	for !(!aaaHandler.state && !bbbHandler.state) {
		log.Debug("aaaHandler.state = %v", aaaHandler.state)
		log.Debug("bbbHandler.state = %v", bbbHandler.state)
		time.Sleep(1 * time.Second)
	}

	assert.EqualValues(t, false, aaaHandler.state, "shutdown effect")
	assert.EqualValues(t, false, bbbHandler.state, "second shutdown effect")
}
