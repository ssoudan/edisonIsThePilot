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
* @Date:   2015-10-21 15:37:36
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 16:29:47
 */

package tracer

import (
	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
)

var log = logger.Log("tracer")

// Tracer is the component that stores the recent position history
type Tracer struct {
	maxPoints        uint32
	nextPosition     uint32
	points           []types.Point
	fullyInitialized bool

	// Channels
	inputChan    chan interface{}
	shutdownChan chan interface{}
	panicChan    chan interface{}
}

// New creates a new Tracer component
func New(maxPoints uint32) *Tracer {
	return &Tracer{
		maxPoints:    maxPoints,
		points:       make([]types.Point, maxPoints),
		nextPosition: 0,
		shutdownChan: make(chan interface{})}
}

// SetInputChan sets the channel where panics are sent
func (t *Tracer) SetInputChan(inputChan chan interface{}) {
	t.inputChan = inputChan
}

// SetPanicChan sets the channel where panics are sent
func (t *Tracer) SetPanicChan(panicChan chan interface{}) {
	t.panicChan = panicChan
}

// Shutdown sets all the state to down and notify the handlers
func (t *Tracer) Shutdown() {

	t.shutdownChan <- 1
	<-t.shutdownChan
}

func (t *Tracer) shutdown() {

	close(t.shutdownChan)
}

func (t *Tracer) processPoint(point types.Point) {
	if t.maxPoints > 0 {
		t.points[t.nextPosition] = point
		if t.nextPosition == t.maxPoints-1 {
			t.nextPosition = 0
			t.fullyInitialized = true
		} else {
			t.nextPosition += 1
		}
	}

}

// MkAddPointMessage creates a new action message to add a Point to the trace
func MkAddPointMessage(point types.Point) interface{} {
	return addPointMessage{point: point}
}

type addPointMessage struct {
	point types.Point
}

type getPointsMessage struct {
	c chan []types.Point
}

// GetPoints returns the Point in the trace
func (t *Tracer) GetPoints() []types.Point {
	c := make(chan []types.Point)
	t.inputChan <- getPointsMessage{c: c}
	return <-c
}

func (t *Tracer) getPoints() []types.Point {
	if t.maxPoints == 0 {
		return []types.Point{}
	}

	if t.fullyInitialized {
		output := make([]types.Point, t.maxPoints)
		copy(output, t.points)
		return output
	} else {
		output := make([]types.Point, t.nextPosition)
		copy(output, t.points)
		return output
	}
}

// Start the event loop of the Tracer component
func (t *Tracer) Start() {

	go func() {

		defer func() {
			if r := recover(); r != nil {
				t.panicChan <- r
			}
		}()

		for {
			select {
			case m := <-t.inputChan:
				switch m := m.(type) {
				case getPointsMessage:
					m.c <- t.getPoints()
				case addPointMessage:
					t.processPoint(m.point)
				}
			case <-t.shutdownChan:
				t.shutdown()

				return
			}

		}
	}()

}
