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
* @Last Modified time: 2015-09-23 13:36:01
 */
package logger

import (
	"fmt"
	"os"
	"sync"
)

const (
	maxWriteCountWithoutCheck = 512
)

type RotateWriter struct {
	lock       sync.Mutex
	filename   string // should be set to the actual filename
	fp         *os.File
	writeCount int32
	maxSize    int64
}

// New returns a new RotateWriter. Return nil if error occurs during setup.
func New(filename string, maxSize int64) *RotateWriter {
	w := &RotateWriter{filename: filename, maxSize: maxSize, writeCount: maxWriteCountWithoutCheck}
	err := w.Rotate()
	if err != nil {
		return nil
	}
	return w
}

// Write satisfies the io.Writer interface and does the rotation when the file get too big - size check is done every maxWriteCountWithoutCheck writes.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fp == nil {
		return 0, fmt.Errorf("no open log file")
	}

	w.writeCount = w.writeCount - 1
	if w.writeCount == 0 {
		w.writeCount = maxWriteCountWithoutCheck
		stat, err := w.fp.Stat()
		if err != nil {
			if stat.Size() > w.maxSize {
				w.Rotate()
			}
		}
	}

	return w.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		_, err = os.Stat(w.filename + ".old")
		if err == nil {
			os.Remove(w.filename + ".old")
		}
		err = os.Rename(w.filename, w.filename+".old")
		if err != nil {
			return
		}
	}

	// Create a file.
	w.fp, err = os.Create(w.filename)
	return
}
