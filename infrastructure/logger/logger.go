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
* @Date:   2015-03-31 21:34:39
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-18 17:53:46
 */

package logger

import (
	"github.com/op/go-logging"
	"os"
)

func Log(name string) *logging.Logger {
	return logging.MustGetLogger(name)
}

func init() {
	format := logging.MustStringFormatter(
		"%{color:bold}%{time:15:04:05.000} %{level:-6s} [%{module}] %{shortfunc:.10s} ▶ %{id:03x}%{color:reset} %{message}",
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)

	backendFormatter := logging.NewBackendFormatter(backend, format)

	// Only errors and more severe messages should be sent to backend1
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.INFO, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
}
