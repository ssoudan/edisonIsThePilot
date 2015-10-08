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
* @Last Modified time: 2015-09-23 09:58:41
 */

package dashboard

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dashboard", func() {

	var (
		dashboard  *Dashboard
		aaaHandler testHandler
		bbbHandler testHandler
		c          chan interface{}
	)

	BeforeEach(func() {
		dashboard = New()

		aaaHandler = testHandler{}
		dashboard.RegisterMessageHandler("AAA", &aaaHandler)

		bbbHandler = testHandler{}
		dashboard.RegisterMessageHandler("BBB", &bbbHandler)

		c = make(chan interface{})

		dashboard.SetInputChan(c)

		dashboard.Start()
	})

	Describe("without channels", func() {

		It("state is false by default", func() {
			Expect(aaaHandler.state).To(BeFalse())
			Expect(bbbHandler.state).To(BeFalse())
		})

		It("changes with messages", func() {
			m := NewMessage(map[string]bool{"AAA": true})

			dashboard.processMessage(m.(message))

			Expect(aaaHandler.state).To(BeTrue())

			Expect(bbbHandler.state).To(BeFalse())

			m = NewMessage(map[string]bool{"BBB": true})

			dashboard.processMessage(m.(message))

			Expect(aaaHandler.state).To(BeTrue())
			Expect(bbbHandler.state).To(BeTrue())

			m = NewMessage(map[string]bool{"BBB": false})

			dashboard.processMessage(m.(message))

			Expect(aaaHandler.state).To(BeTrue())
			Expect(bbbHandler.state).To(BeFalse())

			m = NewMessage(map[string]bool{"BBB": false, "AAA": false})

			dashboard.processMessage(m.(message))

			Expect(aaaHandler.state).To(BeFalse())
			Expect(bbbHandler.state).To(BeFalse())
		})

	})

	Describe("with channels", func() {

		It("state is false by default", func() {
			Expect(aaaHandler.state).To(BeFalse())
			Expect(bbbHandler.state).To(BeFalse())
		})

		It("changes with messages", func() {

			m := NewMessage(map[string]bool{"AAA": true})

			c <- m

			Eventually(aaaHandler.state).Should(BeTrue())
			Consistently(bbbHandler.state).Should(BeFalse())
		})

	})

	Describe("dashboard shutdown", func() {

		It("resets leds when shut down", func() {
			Expect(aaaHandler.state).To(BeFalse())
			Expect(bbbHandler.state).To(BeFalse())

			m := NewMessage(map[string]bool{
				"BBB": true,
				"AAA": true,
			})

			c <- m

			Eventually(aaaHandler.state).Should(BeTrue())
			Eventually(bbbHandler.state).Should(BeTrue())

			dashboard.Shutdown()

			Eventually(aaaHandler.state).Should(BeFalse())
			Eventually(bbbHandler.state).Should(BeFalse())
		})

	})
})

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
