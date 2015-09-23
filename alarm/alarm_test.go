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
* @Date:   2015-09-21 15:42:21
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-23 07:01:41
 */

package alarm

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alarm", func() {
	var (
		handler *testHandler
		alarm   *Alarm
		c       chan interface{}
	)

	BeforeEach(func() {
		handler = &testHandler{}

		alarm = New(handler)

		c = make(chan interface{})

		alarm.SetInputChan(c)

		alarm.Start()

	})

	It("must be false", func() {

		Expect(handler.state).To(BeFalse())
	})

	Describe("after start with no channels", func() {

		It("should enable alarm on message telling so", func() {

			m := NewMessage(true)

			alarm.processMessage(m.(message))

			Expect(handler.state).To(BeTrue())
		})

		It("should disable alarm on message telling so", func() {

			m := NewMessage(false)

			alarm.processMessage(m.(message))

			Expect(handler.state).To(BeFalse())
		})

		It("should enable and disable alarm on messages received in sequence", func() {
			By("receiving a message at true and processing it")
			m := NewMessage(true)

			alarm.processMessage(m.(message))

			Expect(handler.state).To(BeTrue())

			By("receiving a message at false and processing it")
			m = NewMessage(false)

			alarm.processMessage(m.(message))

			Expect(handler.state).To(BeFalse())
		})
	})

	Describe("after start with channels", func() {

		It("should be enabled by a message at true", func() {

			Expect(handler.state).To(BeFalse())
			m := NewMessage(true)

			c <- m

			Eventually(handler.state).Should(BeTrue())
		})

		It("should be disabled by a message at false", func() {
			m := NewMessage(false)

			c <- m

			Eventually(handler.state).Should(BeFalse())
		})

		It("shutdown should disable the alarm", func() {

			m := NewMessage(true)

			c <- m

			Eventually(handler.state).Should(BeTrue())

			alarm.Shutdown()

			Eventually(func() bool { return handler.state }).Should(BeFalse())
		})
	})
})

type testHandler struct {
	mu    sync.Mutex // protects state
	state bool
}

func (h *testHandler) Enable() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.state = true
	return nil
}
func (h *testHandler) Disable() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.state = false
	return nil
}
