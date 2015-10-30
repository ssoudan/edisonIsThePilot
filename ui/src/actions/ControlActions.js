/* 
* @Author: Sebastien Soudan
* @Date:   2015-10-14 16:20:18
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-14 18:51:51
*/

'use strict';

import ControlConstants from '../constants/ControlConstants';

import {dispatch} from '../dispatcher/AppDispatcher'

class ControlActions {

    queryControlState() {
        dispatch({
            actionType: ControlConstants.CONTROL_QUERY_AUTOPILOT,
        });
    }

    queryDashboardState() {
        dispatch({
            actionType: ControlConstants.CONTROL_QUERY_DASHBOARD,
        });
    }

    /**
     * @param  {object} data
     */
    updateAutopilot(data) {
        dispatch({
            actionType: ControlConstants.CONTROL_UPDATE_AUTOPILOT_DATA,
            data: data.data,
            status: data.status,
        });
    }

	updateDashboard(data) {
        dispatch({
            actionType: ControlConstants.CONTROL_UPDATE_DASHBOARD_DATA,
            data: data.data,
            status: data.status,
        });
    }

    /**
     * @param  {object} data
     */
    changeAutopilot(data) {
    	console.log("Going to change autopilot status: ", data)
    	// We disable the control and queue a query 
    	// to the server to get the new state
        dispatch({
            actionType: ControlConstants.CONTROL_CHANGE,
            data: data,
        });
        this.queryControlState();
    }

}
const instance = new ControlActions();
export default instance;