/* 
 * @Author: Sebastien Soudan
 * @Date:   2015-09-10 16:27:47
 * @Last Modified by:   Sebastien Soudan
 * @Last Modified time: 2015-09-16 14:16:39
 */

'use strict';

import MapConstants from '../constants/MapConstants';

import {dispatch} from '../dispatcher/AppDispatcher'

class MapActions {

    /**
     * @param  {object} bounds
     */
    changeBounds(definitions, bounds) {
        dispatch({
            actionType: MapConstants.MAP_CHANGE_BOUNDS,
            bounds: bounds,
            definitions: definitions,
        });
    }

    /**
     * @param  {object} data
     */
    updateData(data) {
        dispatch({
            actionType: MapConstants.MAP_UPDATE_DATA,
            data: data.data,
            part: data.part,
            status: data.status,
        });
    }

}
const instance = new MapActions();
export default instance;