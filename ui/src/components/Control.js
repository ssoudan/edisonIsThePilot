/* 
* @Author: Sebastien Soudan
* @Date:   2015-10-14 16:23:10
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-22 15:06:20
*/

'use strict';


/**
 * Created by ssoudan on 9/5/15.
 */

/**
 * Created by ssoudan on 9/5/15.
 */
import React from 'react';

let { Paper, Styles, Toggle, Slider } = require('material-ui');

const MyRawTheme = require('../theme');
const ThemeManager = require('material-ui/lib/styles/theme-manager');

import DirectionWidget from "./DirectionWidget"
import ControlActions from "../actions/ControlActions";
import ControlStore from "../stores/ControlStore";
import ReactGridLayout from "react-grid-layout";
var FontAwesome = require('react-fontawesome');

// TODO(?) this require a bit of refactoring too many things mixed together here: notification list should go to a separate sub-component at least 

var icons = new Map();
icons.set("NoGPSFix",                   "location-arrow")
icons.set("InvalidGPSData",             "compass")
icons.set("SpeedTooLow",                "dashboard")
icons.set("HeadingErrorOutOfBounds",    "exclamation-triangle")
icons.set("CorrectionAtLimit",          "exclamation")

var text = new Map();
text.set("NoGPSFix",                "No GPS fix")
text.set("InvalidGPSData",          "Invalid GPS data")
text.set("SpeedTooLow",             "Speed too low")
text.set("HeadingErrorOutOfBounds", "Error too large")
text.set("CorrectionAtLimit",       "Correction at limits")

function makeNotificationItem(i, id, notificationsStyle) {
        
        return ([<div key={i+1} _grid={{x: 0, y: i/2, w: 1, h: 1}}>
                    <FontAwesome name={icons.has(id)?icons.get(id):"rocket"}
                                 size='2x'
                                 style={notificationsStyle}
                                 />
                    </div>,
                <div key={i+2} _grid={{x: 1, y: i/2, w: 3, h: 1}}>
                <div> {text.has(id)?text.get(id):"unknown: " + id}</div></div>]);
    }

export default class Control extends React.Component {

    static contextTypes = {
        muiTheme: React.PropTypes.object
    };

    constructor(props) {
        super(props);
        this._onChange = this._onChange.bind(this);
        this._handleSliderChange = this._handleSliderChange.bind(this);
        this._handleToggleChange = this._handleToggleChange.bind(this);
        this.state = {};
        this.controlStoreSubscription = null;
    }

    getChildContext() {
        return {
            muiTheme: ThemeManager.getMuiTheme(MyRawTheme),
        };
    }

    componentDidMount() {
        // window.addEventListener('resize', this.handleResize);
        this.controlStoreSubscription = ControlStore.addListener(this._onChange);
        ControlActions.queryControlState();
        ControlActions.queryDashboardState();
        // TODO(?) we want a timer we re-arm after the action associated to the previous tick are done so we don't build up a queue when something is slow - like the connection to the server.
        // TODO(?) ideally we would like to have a websocket here and the server to send us the updates...
        this.timer = setInterval(this.tick, 2000);
    }

    componentWillUnmount() {
        // window.removeEventListener('resize', this.handleResize);
        if (this.controlStoreSubscription) {
            this.controlStoreSubscription.remove();
            this.controlStoreSubscription = null;
        }
        clearInterval(this.timer);
    }

    tick(){
        console.log("tick")
        // This function is called every 2000 ms. 
        ControlActions.queryControlState();
        ControlActions.queryDashboardState();
    }

    _onChange() {
        this.setState(this.state);
    }

    getControl() {
        return ControlStore.getData();
    }

    getStyles() {
        return {
            group: {
                // paddingTop: '10px',
                margin: '10px',
            },
            containerCentered: {
                height: '100%',
                width: '100%',
                margin: '0 auto',
            },
            unavailableNotification: {
                color: this.getChildContext().muiTheme.rawTheme.palette.accent1Color,
                textAlign: 'center',
                textShadow: '0 1px 0 rgba(0, 0, 0, 0.1)',
            },
            warningNotification: {
                color: this.getChildContext().muiTheme.rawTheme.palette.accent1Color,
                textAlign: 'center',
                textShadow: '0 1px 0 rgba(0, 0, 0, 0.1)',
            },
            normalNotification: {
                color: this.getChildContext().muiTheme.rawTheme.palette.primary3Color,
                textAlign: 'center',
                textShadow: '0 1px 0 rgba(0, 0, 0, 0.1)',
            }
        };
    }

    _handleToggleChange(e, value) {
        ControlActions.changeAutopilot({
             enabled: value, 
             headingOffset: this.refs.headingSlider.getValue(),
        })
    }

    _handleSliderChange(e, value) {
        ControlActions.changeAutopilot({
             enabled: this.refs.enabledToggle.isToggled(),
             headingOffset: this.refs.headingSlider.getValue(),
        })
    }

    render() {
        var offsetLimit = 20; // in degree

        var styles = this.getStyles();
        var controlData = this.getControl()

        var autopilotDataPresent = false
        var dashboardDataPresent = false
        if (Object.keys(controlData).length) {
            if (Object.keys(controlData.autopilot).length) {
                autopilotDataPresent = true
            }
            
            if (Object.keys(controlData.dashboard).length) {
                dashboardDataPresent = true
            }
        }

        
        if (dashboardDataPresent) {
            var column = [];
            var count = 0;
            for (var key of Object.keys(controlData.dashboard)) {
                var raised = controlData.dashboard[key]
                column = column.concat(makeNotificationItem(2*count, key, raised?styles.warningNotification:styles.normalNotification))
                count += 2; // cause we add elements two by two in makeNotificationItem
            }

            // Style
            var setPointColor = this.getChildContext().muiTheme.rawTheme.palette.accent3Color
            var courseColor = this.getChildContext().muiTheme.rawTheme.palette.primary1Color
            var textColor = this.getChildContext().muiTheme.rawTheme.palette.textColor
            var circleColor = this.getChildContext().muiTheme.rawTheme.palette.borderColor
            var radius = Math.min((window.innerWidth - 20)/10, 100) // TODO(?) need a better way to figure out the size of the direction widget: the width and height must fit vertically and horizontally in the control panel for any screen size (mobile included)

            // Data
            var setPoint = controlData.autopilot.setPoint + controlData.autopilot.headingOffset
            var course = controlData.autopilot.course 
            var enabled = controlData.autopilot.enabled

            // TODO(?) improve the layout - both on the visual aspect and on the code organization sides.
            return (
                 <Paper zDepth={1}>
                    <ReactGridLayout className="layout" cols={5} rowHeight={220} isResizable={false} isDraggable={false}>
                    <div key={1} _grid={{x: 0, y: 0, w: 2, h: 1}}>
                        <div style={styles.group}>
                            <div style={styles.containerCentered}>
                                <Slider 
                                    ref="headingSlider"
                                    name="headingOffset" 
                                    description="Course offset"
                                    defaultValue={0} 
                                    value={autopilotDataPresent?controlData.autopilot.headingOffset:0}
                                    onDragStop={this._handleSliderChange}
                                    disabled={!autopilotDataPresent}
                                    max={20} 
                                    min={-20} />
                                <Toggle
                                    ref="enabledToggle"
                                    label="Hold course"
                                    onToggle={this._handleToggleChange}
                                    defaultToggled={autopilotDataPresent?controlData.autopilot.enabled:false}
                                    disabled={!autopilotDataPresent}/>
                            </div>
                        </div>
                    </div>
                    <div key={2} _grid={{x: 2, y: 0, w: 1, h: 1}}>
                        <div style={styles.containerCentered}>
                            <DirectionWidget 
                                enabled={enabled} 
                                setPoint={setPoint} 
                                course={course} 
                                setPointColor={setPointColor} 
                                courseColor={courseColor} 
                                textColor={textColor}
                                circleColor={circleColor}
                                fontSize={16} 
                                radius={radius}
                                margin={5}/>
                        </div>
                    </div>
                    <div key={3} _grid={{x: 3, y: 0, w: 2, h: 1}}>
                        <div style={styles.group}>
                            <div style={styles.containerCentered}>
                                <ReactGridLayout className="layoutWF" cols={4} rowHeight={30} isResizable={false} isDraggable={false}>
                                    {
                                        column
                                    }
                                </ReactGridLayout>
                            </div>
                        </div>
                    </div>
                    </ReactGridLayout> 
                </Paper>
            );
        } else
               return (
                 <Paper zDepth={1}>
                    <ReactGridLayout className="layout" cols={5} rowHeight={110} isResizable={false} isDraggable={false}>
                    <div key={1} _grid={{x: 0, y: 0, w: 2, h: 2}}>
                        <div style={styles.group}>
                            <div style={styles.containerCentered}>
                                <Slider 
                                    ref="headingSlider"
                                    name="headingOffset" 
                                    description="Heading offset"
                                    defaultValue={0} 
                                    value={autopilotDataPresent?controlData.autopilot.headingOffset:0}
                                    onDragStop={this._handleSliderChange}
                                    disabled={!autopilotDataPresent}
                                    max={offsetLimit} 
                                    min={-offsetLimit} />
                                <Toggle
                                    ref="enabledToggle"
                                    label="Is Enabled?"
                                    onToggle={this._handleToggleChange}
                                    defaultToggled={autopilotDataPresent?controlData.autopilot.enabled:false}
                                    disabled={!autopilotDataPresent}/>
                            </div>
                        </div>
                    </div>
                    <div key={2} _grid={{x: 2, y: 0, w: 3, h: 2}}>
                        <div style={styles.group}>
                            <div style={styles.containerCentered}>
                                <ReactGridLayout className="layoutW" cols={7} rowHeight={110} isResizable={false} isDraggable={false}>
                                <div key={1} _grid={{x: 0, y: 0, w: 1, h: 1}}>
                                <FontAwesome name='exclamation-triangle'
                                             style={styles.unavailableNotification}
                                             size='2x'/>
                                </div>
                                <div key={2} _grid={{x: 1, y: 0, w: 3, h: 1}}>
                                <div style={styles.unavailableNotification}>Missing informations -- check connection to the autopilot</div>
                                </div>
                                </ReactGridLayout>
                            </div>
                        </div>
                    </div>
                    </ReactGridLayout> 
                </Paper>
            );
    }
}

  
Control.childContextTypes = {
    muiTheme: React.PropTypes.object
};
