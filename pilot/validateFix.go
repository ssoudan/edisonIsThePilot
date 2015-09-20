/*
* @Author: Sebastien Soudan
* @Date:   2015-09-20 22:08:20
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-09-20 22:09:10
 */

package pilot

type FixStatus byte

const (
	NOFIX    = 0
	FIX      = 1
	DGPS_FIX = 2
)

func validateFixStatus(fix FixStatus) (alarm Alarm, ledEnabled bool) {
	alarm = Alarm(UNRAISED)
	ledEnabled = false

	switch fix {
	case NOFIX:

		alarm = RAISED
		ledEnabled = true

	case FIX, DGPS_FIX:
		// blah
	}

	return
}
