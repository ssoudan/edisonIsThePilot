/*
 * @Author: Sebastien Soudan
 * @Date:   2015-09-10 16:28:18
 * @Last Modified by:   Sebastien Soudan
 * @Last Modified time: 2015-10-21 12:22:01
 */

var rest = require('rest');
var mime = require('rest/interceptor/mime');
var client = rest.wrap(mime, { mime: 'application/json' });

import MapActions from "../actions/MapActions";
import ControlActions from "../actions/ControlActions";
import Constants from "../constants/Constants";
var SPLITS = 1;

class Api {

    getAutopilot() {
        var query = {
            path: "/api/autopilot",
        };

        client(query)
            .then(res => {
                if (res.status.code == 200)
                       ControlActions.updateAutopilot({
                            data: res.entity,
                            status: Constants.OK
                        });
                else {
                    console.error("failed with status:", res.status.text);
                        ControlActions.updateAutopilot({
                            data: res.entity,
                            status: Constants.KO
                        });
                }
            }).catch(res => {
                console.log(res)
                ControlActions.updateAutopilot({
                        data: res.entity,
                        status: Constants.KO
                    });
            });
    }

     getDashboard() {
        var query = {
            path: "/api/dashboard",
        };

        client(query)
            .then(res => {
                if (res.status.code == 200)
                       ControlActions.updateDashboard({
                            data: res.entity,
                            status: Constants.OK
                        });
                else {
                    console.error("failed with status:", res.status.text);
                        ControlActions.updateDashboard({
                            data: res.entity,
                            status: Constants.KO
                        });
                }
            }).catch(res => {
                console.log(res)
                ControlActions.updateDashboard({
                            data: res.entity,
                            status: Constants.KO
                        });
            });
    }

    setAutopilot(data) {
        console.log("Data", data)
         var query = {
            path: "/api/autopilot",
            method: 'PUT',
            entity: data, 
        };

        client(query)
            .then(res => {
                if (res.status.code == 200)
                       ControlActions.updateAutopilot({
                            data: res.entity,
                            status: Constants.OK
                        });
                else {
                    console.error("failed with status:", res.status.text);
                }
            }).catch(res => { 
                console.log(res)
                ControlActions.updateAutopilot({
                        status: Constants.KO
                });
            });
    }

    fetchPoints(q) {
        var bounds = q.bounds;
        var definitions = q.definitions;
        var ne = bounds.getNorthEast();
        var sw = bounds.getSouthWest();

        var minLat_ = sw.lat();
        var maxLat_ = ne.lat();
        var maxLng_ = ne.lng();
        var minLng_ = sw.lng();

        var w_ = Math.abs(maxLng_ - minLng_) / SPLITS;
        var h_ = Math.abs(maxLat_ - minLat_) / SPLITS;
        console.log("w_=", w_);
        console.log("h_=", h_);

        console.log(SPLITS);

        for (var i = 0; i < SPLITS; i++) {
            for (var j = 0; j < SPLITS; j++) {
                this.do_query(definitions, i, j, h_, w_, minLat_, minLng_);
            }
        }
    }

    do_query(definitions, i, j, h_, w_, minLat_, minLng_) {
        var part = (i * SPLITS + j);
        var query = {
            path: "/api/points",
            params: {
                vertDef: definitions.vertDef / SPLITS,
                horizDef: definitions.horizDef / SPLITS,
                minLat: minLat_ + i * h_,
                maxLat: minLat_ + (i + 1) * h_,
                minLong: minLng_ + j * w_,
                maxLong: minLng_ + (j + 1) * w_,
            }
        };

        console.log("about to fetch data (part=" + part + "): " + query.path + "?minLat=" + query.params.minLat + "&maxLat=" + query.params.maxLat + "&minLong=" + query.params.minLong + "&maxLong=" + query.params.maxLong);
        console.log(query);

        client(query)
            .then(res => {
            if (res.status.code == 200)
                return res.entity;
            else
            {
                console.error("failed with status:", res.status.text);
                return [];
                        }
                    })
            .then(res => {
                return res.map(pt => {
                    var w = !pt.weigth ? pt.weight : 1;
                    return {
                        lat: pt.latitude,
                        lng: pt.longitude,
                        w: w,
                        key: pt.latitude + ":" + pt.longitude + ":" + w,
                    }
                    })
                },
                console.error, console.error)
            .then(newPoints => {
                console.log("fetched " + newPoints.length + " points for part" + part);

                if (newPoints.length != 0) {
                    MapActions.updateData({
                        data: newPoints,
                        part: part,
                        status: Constants.OK
                    });
                } else {
                    MapActions.updateData({
                        data: [],
                        part: part,
                        status: Constants.KO
                    });
                }
                console.log("pushed");
            }).catch(res => console.log(res));
    };

}


const instance = new Api();
export default instance;