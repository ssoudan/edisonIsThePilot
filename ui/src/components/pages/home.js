/**
 * Created by ssoudan on 9/5/15.
 */

/**
 * Created by ssoudan on 9/5/15.
 */
import React from 'react';

let { Styles } = require('material-ui');
let { Spacing } = Styles;
let SimpleMap = require('../SimpleMap');
let Control = require('../Control');

var mui = require('material-ui');
var DesktopGutter = mui.Styles.Spacing.desktopGutter;

export default class Home extends React.Component {

    constructor(props) {
        super(props);
        this.handleResize = this.handleResize.bind(this);
        this.state = {windowWidth: window.innerWidth, windowHeight: window.innerHeight};
    }

    handleResize(e) {
        this.setState({windowWidth: window.innerWidth, windowHeight: window.innerHeight});
    }

    componentDidMount() {
        window.addEventListener('resize', this.handleResize);
    }

    componentWillUnmount() {
        window.removeEventListener('resize', this.handleResize);
    }

    getStyles(height, width) {
        return {
            root: {
                paddingTop: Spacing.desktopKeylineIncrement,
                height: height - 166 + DesktopGutter + 'px',
            },
            group: {
                paddingTop: '10px',
                margin: '10px',
            },
            map: {
                height: height - 402 + 'px',
                width: '100%',
            }
        };
    }

    render() {

        let styles = this.getStyles(this.state.windowHeight, this.state.windowWidth);
        if (typeof google !== 'undefined') {
            console.log("google is defined")
            return (
                <div style={styles.root}>
                        <Control/>
                        <div style={styles.group}>
                            <div style={styles.map}>
                                <SimpleMap/>
                            </div>
                        </div>
                </div>);
        } else {
            return (
                <div style={styles.root}>
                        <Control/>
                        <div style={styles.group}>
                            Offline mode &mdash; check your can reach the Internet to get the map
                        </div>
                </div>);
        }
        
    }
}
