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
* @Last Modified time: 2015-10-21 16:28:00
 */

package tracer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
	"time"
)

var _ = Describe("tracer", func() {

	Describe("with no seat", func() {
		var (
			tracer *Tracer
			c      chan interface{}
		)

		BeforeEach(func() {
			tracer = New(0)

			c = make(chan interface{})

			tracer.SetInputChan(c)

			tracer.Start()
		})

		It("has no points at first", func() {
			Expect(tracer.GetPoints()).To(Equal([]types.Point{}))
		})

		It("doesn't store any point", func() {

			m := types.Point{
				Latitude:  45.,
				Longitude: 4.,
				Time:      types.JSONTime(time.Now()),
			}

			c <- MkAddPointMessage(m)

			Consistently(tracer.GetPoints).Should(Equal([]types.Point{}))

		})

	})

	Describe("with some seats", func() {
		var (
			tracer *Tracer
			c      chan interface{}
		)

		BeforeEach(func() {
			tracer = New(3)

			c = make(chan interface{})

			tracer.SetInputChan(c)

			tracer.Start()
		})

		It("has no points at first", func() {
			Expect(tracer.GetPoints()).To(Equal([]types.Point{}))
		})

		It("can store a point", func() {

			m := types.Point{
				Latitude:  45.,
				Longitude: 4.,
				Time:      types.JSONTime(time.Now()),
			}

			c <- MkAddPointMessage(m)

			Eventually(tracer.GetPoints).Should(Equal([]types.Point{m}))

		})

		It("accumulates points as they come", func() {
			stuff := make([]types.Point, 3)

			for i := 0; i < 3; i++ {
				m := types.Point{
					Latitude:  45. + float64(i),
					Longitude: 4.,
					Time:      types.JSONTime(time.Now()),
				}

				c <- MkAddPointMessage(m)
				stuff[i] = m
			}

			Eventually(tracer.GetPoints).Should(Equal(stuff))

		})

		It("wraps", func() {
			stuff := make([]types.Point, 3)

			for i := 0; i < 6; i++ {
				m := types.Point{
					Latitude:  45. + float64(i),
					Longitude: 4.,
					Time:      types.JSONTime(time.Now()),
				}

				c <- MkAddPointMessage(m)
				stuff[i%3] = m
			}

			Eventually(tracer.GetPoints).Should(Equal(stuff))

		})

	})

})
