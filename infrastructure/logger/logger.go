/*
* @Author: Sebastien Soudan
* @Date:   2015-03-31 21:34:39
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-04-16 13:55:49
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
