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
* @Date:   2015-09-28 22:13:28
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 16:39:00
 */

package webserver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/ssoudan/edisonIsThePilot/infrastructure/logger"
	"github.com/ssoudan/edisonIsThePilot/infrastructure/types"
	"github.com/ssoudan/edisonIsThePilot/pilot"
)

var log = logger.Log("webserver")

// VersionEndpoint is the endpoint providing the version of this piece of software
func VersionEndpoint(version string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("Edison is the pilot - %s", version))
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "static/", 301)
}

// Control is the serializable structure used to change the autopilot state
type Control struct {
	Enabled       bool    `json:"enabled"`
	HeadingOffset float64 `json:"headingOffset"`
}

// Autopilot is the serializable structure used to get the autopilot state
type Autopilot struct {
	Enabled       bool    `json:"enabled"`
	HeadingOffset float64 `json:"headingOffset"`
	SetPoint      float64 `json:"setPoint"`
	Course        float64 `json:"course"`
	Speed         float64 `json:"speed"`
}

// Webserver is a web server component exposing both static files (static/) and the api (api/)
type Webserver struct {
	pilot     pilotable
	dashboard queryable
	tracer    tracer
	version   string

	panicChan chan interface{}
}

type pilotable interface {
	GetInfoAction() pilot.Info
	Enable() error
	Disable() error
	SetOffset(headingOffset float64) error
}

type queryable interface {
	GetDashboardInfoAction() map[string]bool
}

type tracer interface {
	GetPoints() []types.Point
}

// New creates a new Webserver
func New(version string) *Webserver {
	return &Webserver{version: version}
}

// SetPanicChan sets the channel where panics are sent
func (ws *Webserver) SetPanicChan(panicChan chan interface{}) {
	ws.panicChan = panicChan
}

// SetPilot sets the Pilot used by the Webserver
func (ws *Webserver) SetPilot(p pilotable) {
	ws.pilot = p
}

// SetTracer sets the Tracer used by the Webserver
func (ws *Webserver) SetTracer(p tracer) {
	ws.tracer = p
}

// SetDashboard sets the Dashboard where the Webserver will get the LED state from
func (ws *Webserver) SetDashboard(dashboard queryable) {
	ws.dashboard = dashboard
}

// Start the webserver (non-blocking)
func (ws *Webserver) Start() {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ws.panicChan <- r
			}
		}()

		api := rest.NewApi()
		api.Use(rest.DefaultDevStack...)

		router, err := rest.MakeRouter(
			rest.Get("/points", func(w rest.ResponseWriter, req *rest.Request) {
				if _, ok := ws.pilot.(pilotable); !ok {
					log.Error("WS is not initialized")
					rest.Error(w, "WS is not initialized", http.StatusInternalServerError)
					return
				}
				w.WriteJson(ws.tracer.GetPoints())
			}),
			rest.Get("/autopilot", func(w rest.ResponseWriter, req *rest.Request) {
				if _, ok := ws.pilot.(pilotable); !ok {
					log.Error("WS is not initialized")
					rest.Error(w, "WS is not initialized", http.StatusInternalServerError)
					return
				}

				pi := ws.pilot.GetInfoAction()
				w.WriteJson(Autopilot{
					Enabled:       pi.Enabled,
					HeadingOffset: pi.HeadingOffset,
					SetPoint:      pi.SetPoint,
					Course:        pi.Course,
					Speed:         pi.Speed,
				})
				return
			}),
			rest.Get("/dashboard", func(w rest.ResponseWriter, req *rest.Request) {
				if _, ok := ws.dashboard.(queryable); !ok {
					log.Error("WS is not initialized")
					rest.Error(w, "WS is not initialized", http.StatusInternalServerError)
					return
				}
				data := ws.dashboard.GetDashboardInfoAction()
				w.WriteJson(data)
			}),
			rest.Put("/autopilot", func(w rest.ResponseWriter, r *rest.Request) {
				if _, ok := ws.pilot.(pilotable); !ok {
					log.Error("WS is not initialized")
					rest.Error(w, "WS is not initialized", http.StatusInternalServerError)
					return
				}

				autopilot := Control{}
				err := r.DecodeJsonPayload(&autopilot)

				if err != nil {
					log.Error("Failed to parse json:", err)
					rest.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				log.Info("Got %v", autopilot)
				ws.pilot.SetOffset(autopilot.HeadingOffset)
				if autopilot.Enabled {
					err = ws.pilot.Enable()
					if err != nil {
						log.Error("Failed to enable autopilot:", err)
						rest.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					err = ws.pilot.Disable()
					if err != nil {
						log.Error("Failed to disable autopilot:", err)
						rest.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				w.WriteJson(map[string]string{"status": "OK"})
			}),
		)
		if err != nil {
			log.Panic(err)
		}
		api.SetApp(router)

		http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

		http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("."))))

		http.HandleFunc("/version", VersionEndpoint(ws.version))
		http.HandleFunc("/", redirect)

		err = http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Panic("Already running! or there is something else alive on port 8000 - exiting")
		}
	}()
}
