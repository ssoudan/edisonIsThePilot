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
* @Date:   2015-09-27 22:18:56
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-28 23:35:30
 */

package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/utils"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/webserver"
)

var log = logger.Log("systemCalibration")

var Version = "unknown"

type parameters map[string]string

type step struct {
	step_number uint32
	parameters  parameters
}

type point struct {
	timestamp      time.Time
	course         float64
	speed          float64
	delta_steering float64
	latitude       float64
	longitude      float64
	validity       bool
	step_number    uint32
}

const (
	DONE    = "DONE"
	RUNNING = "RUNNING"
)

type experiment struct {
	state        string
	test_type    string
	parameters   parameters
	plot_command string
	steps        []step
	points       []point
}

var world struct {
	mu    sync.RWMutex
	infos []experiment // protected by mu
}

func info(w http.ResponseWriter, r *http.Request) {
	world.mu.RLock()
	defer world.mu.RUnlock()
	json.NewEncoder(w).Encode(world.infos)
}

func main() {

	// TODO(ssoudan) parse inputs

	// TODO(ssoudan) define scenario input format
	// we want scenario with different speed
	// we want scenario with positive and negative steering changes
	// we want different values of those changes
	// we want different initial steering (do we?)

	log.Info("Starting -- version %s", Version)

	panicChan := make(chan interface{})
	defer func() {
		if r := recover(); r != nil {
			panicChan <- r
		}
	}()

	go func() {
		select {
		case m := <-panicChan:

			// TODO(ssoudan) do what has to be done here

			log.Fatalf("Version %v -- Received a panic error -- exiting: %v", Version, m)
		}
	}()

	// - Expose data via a rest api that can be used by matlab as JSON
	// {
	//  state: "RUNNING", // , "DONE"
	// 	settings: {
	// 		test_type: "test type",
	// 		parameters: {
	// 					x: y,
	// 					xx: yy,
	// 					start_time: ...,
	// 					}
	// 		steps: [{
	// 			step_number: xx,
	// 			parameters: {
	// 					x: y,
	// 					xx: yy,
	// 					start_time: ...,
	// 					}
	// 		}]
	// 	},
	// 	points: [{timestamp: t,
	// 			course: c,
	// 			speed: s,
	// 			delta_steering: ds,
	// 			latitude: lat,
	// 			longitude: lng,
	// 			validity: v,
	// 			step_number: xxx,
	// 			}]
	// }

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}
		}()

		http.HandleFunc("/", webserver.VersionEndpoint(Version))
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Panic("Already running! or something is living on port 8000 - exiting")
		}
	}()

	// TODO(ssoudan) populate the scenario

	// TODO(ssoudan) init the steering and gps modules
	// TODO(ssoudan) defer Shutdown()

	// TODO(ssoudan) show the pilot what is going to happen

	// TODO(ssoudan) display pilot instructions (like speed)

	// TODO(ssoudan) wait for start signal

	// TODO(ssoudan) run the scenario

	// TODO(ssoudan) write logs to file

	// FUTURE(ssoudan) Can we imagine to generate the input from matlab? :)

	// TODO(ssoudan) write to a file as well

	// TODO(ssoudan) tel the pilot the calibration test is over and data can be collected

	// FUTURE(ssoudan) support periodic response testing

	// Wait until we receive a signal
	utils.WaitForInterrupt(func() {
		log.Info("Interrupted - exiting")
		log.Info("Exiting -- version %v", Version)
	})
}
