/* 
* @Author: Sebastien Soudan
* @Date:   2015-10-15 12:04:10
* @Last Modified by:   Sebastien Soudan
* @Last Modified time: 2015-10-16 17:32:01
*/

'use strict';

// TODO(?) convert to ES6+: see SimpleMap.js for an example

var React = require('react');
import {default as ReactART, Circle} from 'react-art';
var Group = ReactART.Group;
var Shape = ReactART.Shape;
import Path from 'paths-js/path';
var Surface = ReactART.Surface;
var Text = ReactART.Text;

var DirectionWidget = React.createClass({
  /**
   * Initialize state members.
   */
  getInitialState: function() {
    return {};
  },

  componentDidMount: function() {
    // Nothing special
  },

  componentWillUnmount: function() {
    // Nothing special
  },

  render: function() {
    var radius = this.props.radius
    var margin = this.props.margin

    return (
      <Surface
            width={2*radius + 2*margin}
            height={2*radius + 2*margin}>
                {this.renderGraphic(margin, 
                                    radius, 
                                    this.props.setPoint, 
                                    this.props.course)}
      </Surface>
    );
  },

  
  makeArrow: function(xCenter, yCenter, radius, angle, color, strokeWidth) {
    var angleInRadian = angle * Math.PI / 180
    var path = Path()
        .moveto(xCenter + 0.1*(radius * Math.sin(angleInRadian)), 
                yCenter - 0.1*(radius * Math.cos(angle * Math.PI/180)))
        .lineto(xCenter + 0.8*(radius * Math.sin(angleInRadian)), 
                yCenter - 0.8*(radius * Math.cos(angle * Math.PI/180)));
    return (
        <Shape 
            d={path.print()} 
            strokeWidth={strokeWidth} 
            stroke={color} />
    ); 

  },

  renderGraphic: function(margin, radius, setPoint, course) {
    // Border
    var circleColor = this.props.circleColor
    var circlePath = "M-" + radius + ",0a" + radius + "," + 
                     radius + " 0 1,0 " + 2*radius + ",0a" + 
                     radius + "," + radius + " 0 1,0 -" + 2*radius + ",0";
    var circle = <Shape
                        d={circlePath}
                        key="a"
                        strokeWidth={4}
                        stroke={circleColor}
                        width={2*radius} height={2*radius}
                        x={radius} y={radius}
                        opacity={0.3}/>

    // Arrows
    var setPointColor = this.props.setPointColor
    var courseColor = this.props.courseColor

    var sp = {}
    if (this.props.enabled) {
        sp = this.makeArrow(radius, radius, radius, setPoint, setPointColor, 2)
    }
    var course = this.makeArrow(radius, radius, radius, course, courseColor, 4)
    
    // Text
    var textColor = this.props.textColor
    var fontSize = this.props.fontSize
    var font = 'bold '+fontSize+'px "Arial"'

    return (
        <Group x={margin} y={margin}>
            {circle}
            {sp}
            {course}
            <Text x={radius} y={0} alignment="middle" fill={textColor} font={font}>N</Text>
            <Text x={radius} y={2*radius-fontSize} alignment="middle" fill={textColor} font={font}>S</Text>
            <Text x={2*radius-fontSize/2} y={radius-fontSize/2} alignment="middle" fill={textColor} font={font}>E</Text>
            <Text x={0+fontSize/2} y={radius-fontSize/2} alignment="middle" fill={textColor} font={font}>W</Text>
        </Group>);
  }
});



module.exports = DirectionWidget;