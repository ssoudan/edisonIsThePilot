/**
 * Created by ssoudan on 9/7/15.
 */
import {default as React, addons} from "react/addons";

import {GoogleMap} from "react-google-maps";
import {Marker} from "react-google-maps";
import {Heatmap} from "react-google-maps";
import MapActions from "../actions/MapActions";
import PointStore from "../stores/PointStore";

const MyRawTheme = require('../theme');
const ThemeManager = require('material-ui/lib/styles/theme-manager');

// TODO(?) remove the ThemeManager, gradient, styling stuffs from here: This component should be independant of the material-ui framework, move these stuffs as props and use the theme manager at an upper level: page level or group of components?

var gradient = [
    'rgba(  33,  150, 243, 0)',
    'rgba(  33,  150, 243, 1)',
    'rgba(  65,  140, 220, 1)',
    'rgba(  96,  131, 197, 1)',
    'rgba(  128, 121, 174, 1)',
    'rgba(  160, 111, 151, 1)',
    'rgba(  192, 101, 128, 1)',
    'rgba(  223, 92,  105, 1)',
    'rgba(  255, 82,  82,  1)',
]

function makeMapStyle(color) {
    return [
        {
            "featureType": "all",
            "elementType": "all",
            "stylers": [
                {
                    "saturation": -100
                },
                {
                    "gamma": 0.2
                }
            ]
        },
        {
            "featureType": "all",
            "elementType": "labels.text.fill",
            "stylers": [
                {
                    "saturation": 36
                },
                {
                    "color": color,
                },
                {
                    "lightness": 20
                }
            ]
        },
        {
            "featureType": "all",
            "elementType": "labels.text.stroke",
            "stylers": [
                {
                    "visibility": "off"
                },
                {
                    "color": "#000000"
                },
                {
                    "lightness": 16
                }
            ]
        },
        {
            "featureType": "administrative",
            "elementType": "geometry.stroke",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 17
                },
                {
                    "weight": 1.2
                }
            ]
        },
        {
            "featureType": "road.highway",
            "elementType": "geometry.fill",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 17
                }
            ]
        },
        {
            "featureType": "road.highway",
            "elementType": "geometry.stroke",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 29
                },
                {
                    "weight": 0.2
                }
            ]
        },
        {
            "featureType": "road.arterial",
            "elementType": "geometry",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 18
                }
            ]
        },
        {
            "featureType": "road.local",
            "elementType": "geometry",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 16
                }
            ]
        },
        {
            "featureType": "road.arterial",
            "elementType": "labels.icon",
            "stylers": [
                {
                    "visibility": "off"
                },
                {
                    "lightness": 30
                },
                {
                    "gamma": 1.00
                }
            ]
        },
        {
            "featureType": "road.local",
            "elementType": "labels.icon",
            "stylers": [
                {
                    "visibility": "on"
                },
                {
                    "lightness": 30
                },
                {
                    "gamma": 1.00
                }
            ]
        },
        {
            "featureType": "transit",
            "elementType": "geometry",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 19
                }
            ]
        },
        {
            "featureType": "water",
            "elementType": "geometry",
            "stylers": [
                {
                    "color": "#000000"
                },
                {
                    "lightness": 17
                }
            ]
        }
    ];
};

export default class SimpleMap extends React.Component {

    getChildContext() {
        return {
            muiTheme: ThemeManager.getMuiTheme(MyRawTheme),
        };
    }

    constructor() {
        super();
        this._handle_map_bounds_changed = this._handle_map_bounds_changed.bind(this);
        this._onChange = this._onChange.bind(this);
        this.state = {
            windowWidth: window.innerWidth,
            windowHeight: window.innerHeight,
        };
        this.pointStoreSubscription = null;
    }


    componentWillUnmount() {
        window.removeEventListener('resize', this._onChange);
        if (this.pointStoreSubscription) {
            this.pointStoreSubscription.remove();
            this.pointStoreSubscription = null;
        }
    }

    getStyles(height, width) {
        return {
            root: {
                height: '50%' // height - 228 + 'px',
            }
        };
    }

    getMapState() {
        return {
            windowWidth: window.innerWidth,
            windowHeight: window.innerHeight,
        }
    }

    /**
     * Event handler for 'change' events coming from the PointStore
     */
    _onChange() {
        this.setState(this.getMapState());
    }

    componentDidMount() {
        window.addEventListener('resize', this._onChange);
        this.pointStoreSubscription = PointStore.addListener(this._onChange);
    }


    _handle_map_bounds_changed() {
        if (this.refs.map) {
            console.log("Bounds changed! ");
            console.log(this.state);
            var bounds = this.refs.map.getBounds();
            MapActions.changeBounds({
                vertDef: this.state.windowHeight / 2,
                horizDef: this.state.windowWidth / 2,
            }, bounds);
        }
    }

    getPoints() {
        return PointStore.getData();
    }

    render() {

        let styles = this.getStyles(this.state.windowHeight, this.state.windowWidth);

        console.log("SimpleMap.render()");

        var points = this.getPoints();
        var mapStyle = makeMapStyle(this.getChildContext().muiTheme.rawTheme.palette.accent2Color);

        return (
            <GoogleMap 
                       ref="map"
                       containerProps={{
                            style: {
                                height: "100%",
                                },
                       }}
                       mapTypeId={google.maps.MapTypeId.TERRAIN}
                       defaultZoom={4}
                       defaultCenter={{lat: 45., lng:5.}}
                       defaultOptions={{styles: mapStyle, }}
                       onBoundsChanged={this._handle_map_bounds_changed}>
                <Heatmap radius={8}
                         dissipating={true}
                         gradient={gradient}
                         data={points.map((m) => {
                                return {
                                    location: new google.maps.LatLng(m.lat, m.lng),
                                    weight: 1
                                 }})}/>
            </GoogleMap>);
    };
}

SimpleMap.childContextTypes = {
    muiTheme: React.PropTypes.object
};