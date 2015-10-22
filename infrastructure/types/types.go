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
* @Date:   2015-10-19 15:35:41
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-21 16:47:46
 */

package types

import (
	"fmt"
	"time"
)

// Enablable is something that can be Enable() or Disable()
type Enablable interface {
	Enable() error
	Disable() error
}

// Readable is something that can be read
type Readable interface {
	Value() (bool, error)
}

// JSONTime is an alias for time.Time with a custom serialization
type JSONTime time.Time

// MarshalJSON does the JSON serialization for JSONTime
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(stamp), nil
}

// JSONDuration is an alias for time.Duration with a custom serialization
type JSONDuration time.Duration

// MarshalJSON does the JSON serialization for JSONDuration
func (d JSONDuration) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%f", time.Duration(d).Seconds())
	return []byte(stamp), nil
}

// Point is a position on earth at a specified time instant
type Point struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Time      JSONTime `json:"time"`
}
