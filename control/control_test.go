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
* @Last Modified time: 2015-09-23 11:27:23
 */
package control

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Control", func() {

	var (
		control *Control
		input   *testHandler
		pilot   *testPilot
	)

	BeforeEach(func() {
		input = &testHandler{value: false, err: nil}
		pilot = &testPilot{}
		control = New(input, pilot)

		control.Start()

	})

	Describe("when input is high", func() {

		BeforeEach(func() {
			input.value = true
		})

		It("enables the control", func() {
			Eventually(func() bool { return pilot.state }).Should(BeTrue())
		})

	})

	Describe("when input is down", func() {

		BeforeEach(func() {
			input.value = false
		})

		It("disables the control", func() {
			Eventually(func() bool { return pilot.state }).Should(BeFalse())
		})

	})

	Describe("when shutting down and input was high", func() {

		BeforeEach(func() {
			input.value = true

		})

		It("doen't change the state", func() {
			Eventually(func() bool { return pilot.state }).Should(BeTrue())

			control.Shutdown()

			Consistently(func() bool { return pilot.state }).Should(BeTrue())
		})

	})

	Describe("when shutting down and input was low", func() {

		BeforeEach(func() {
			input.value = false

		})

		It("doen't change the state", func() {
			Eventually(func() bool { return pilot.state }).Should(BeFalse())

			control.Shutdown()

			Consistently(func() bool { return pilot.state }).Should(BeFalse())
		})

	})

})

type testPilot struct {
	state bool
}

func (t *testPilot) Enable() error {
	t.state = true
	return nil
}

func (t *testPilot) Disable() error {
	t.state = false
	return nil
}

type testHandler struct {
	value bool
	err   error
}

func (t testHandler) Value() (bool, error) {
	return t.value, t.err
}
