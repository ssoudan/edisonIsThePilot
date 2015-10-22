/* 
* @Author: Sebastien Soudan
* @Date:   2015-10-14 16:32:57
* @Last Modified by:   ssoudan
* @Last Modified time: 2015-10-15 17:09:04
*/

/**
 * Created by ssoudan on 9/10/15.
 */

'use strict';

import AppDispatcher from '../dispatcher/AppDispatcher';
import {Store} from 'flux/utils'

import ControlConstants from '../constants/ControlConstants';
import Constants from '../constants/Constants';
import Api from '../api/Api';

class ControlStore extends Store {

  constructor(props) {
    super(props);
    this.state = {dashboard: {}, autopilot: {}};  
  }

  __onDispatch(action) {
    switch (action.actionType) {

      case ControlConstants.CONTROL_QUERY_AUTOPILOT:
        Api.getAutopilot();
        break;

      case ControlConstants.CONTROL_QUERY_DASHBOARD:
        Api.getDashboard();
        break;

      case ControlConstants.CONTROL_CHANGE:
      	// this._changeStateAutopilot({});
        Api.setAutopilot(action.data)
      	// this.__emitChange();
      	break;
      
      case ControlConstants.CONTROL_UPDATE_AUTOPILOT_DATA:
        if (action.status == Constants.OK) {
          this._changeStateAutopilot(action.data);
          this.__emitChange();
        } else {
          this._changeStateAutopilot({});
          this.__emitChange();
        }
        break;
      
      case ControlConstants.CONTROL_UPDATE_DASHBOARD_DATA:
        if (action.status == Constants.OK) {
          this._changeStateDashboard(action.data);
          this.__emitChange();
        } else {
          this._changeStateDashboard({});
          this.__emitChange();
        }
        break;
      default:
        // no op
    }
  }

  getData() {
    return this.state
  }

  _changeStateDashboard(data) {
    this.state.dashboard = data;    
  }

  _changeStateAutopilot(data) {
    this.state.autopilot = data;    
  }

}

const instance = new ControlStore(AppDispatcher);
export default instance;