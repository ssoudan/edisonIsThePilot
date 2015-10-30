/**
 * Created by ssoudan on 9/10/15.
 */

'use strict';

import AppDispatcher from '../dispatcher/AppDispatcher';
import {Store} from 'flux/utils'

import MapConstants from '../constants/MapConstants';
import Constants from '../constants/Constants';
import Api from '../api/Api';

class PointStore extends Store {

  constructor(props) {
    super(props);
    this.queryInProgress = false;
    this.nextQuery = null;
    this.state = {};  
  }

  __onDispatch(action) {
    switch (action.actionType) {

      case MapConstants.MAP_CHANGE_BOUNDS:
        if (this.queryInProgress) {
          this.nextQuery = {
            bounds: action.bounds,
            definitions: action.definitions,
          }
        } else {
          this.queryInProgress = true;
          Api.fetchPoints({
            bounds: action.bounds,
            definitions: action.definitions,
          });
        }        
        break;

      case MapConstants.MAP_UPDATE_DATA:
        if (action.status == Constants.OK) {
          this._changeState(action.part, action.data);
         this.__emitChange();
        }

        if (this.nextQuery != null) {
          Api.fetchPoints(this.nextQuery);
          this.nextQuery = null;
        } else {
          this.queryInProgress = false;
        }
        
        break;

      default:
        // no op
    }
  }

  getData() {
    let arr = [];
    for (var e in this.state) {
      arr = arr.concat(this.state[e]);
    }
    return arr
  }

  _changeState(part, data) {
    this.state[part] = data;    
  }

}

const instance = new PointStore(AppDispatcher);
export default instance;