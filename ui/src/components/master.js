/**
 * Created by ssoudan on 9/5/15.
 */

import React from 'react';
import Router from 'react-router';
var RouteHandler = Router.RouteHandler;

import mui from 'material-ui';
import AppLeftNav from './app-left-nav.js';
import FullWidthSection from './full-width-section.js';

const MyRawTheme = require('../theme');
const ThemeManager = require('material-ui/lib/styles/theme-manager');

const Colors = mui.Styles.Colors;

var { AppBar, AppCanvas} = mui;

class Master extends React.Component {

    constructor() {
        super();
        this._onLeftIconButtonTouchTap = this._onLeftIconButtonTouchTap.bind(this);
    }

    getChildContext() {
        return {
            muiTheme: ThemeManager.getMuiTheme(MyRawTheme),
        };
    }

    getStyles() {
        var darkWhite = Colors.darkWhite;
        return {
            footer: {
                backgroundColor: this.getChildContext().muiTheme.rawTheme.palette.primary3Color,
                textAlign: 'center',
                marginTop: '10px',
            },
            a: {
                color: darkWhite
            },
            p: {
                margin: '0 auto',
                padding: '0',
                color: this.getChildContext().muiTheme.rawTheme.palette.alternateTextColor,
                maxWidth: '335px'
            },
            iconButton: {
                color: darkWhite
            }, 
            appBar: {
                backgroundColor: this.getChildContext().muiTheme.rawTheme.palette.primary1Color,
                margin: '0 auto'
            }
        };
    }

    render() {
        var styles = this.getStyles();
        var title = this.context.router.isActive('get-started') ? 'Get Started' :
            this.context.router.isActive('home') ? 'Home' : 'edisonIsThePilot';

        return (
            <AppCanvas predefinedLayout={1}>

                <AppBar
                    onLeftIconButtonTouchTap={this._onLeftIconButtonTouchTap}
                    title={title}
                    zDepth={0}
                    style={styles.appBar}/>

                <AppLeftNav ref="leftNav"/>

                <RouteHandler/>

                <FullWidthSection style={styles.footer}>
                    <p style={styles.p}>
                        © 2015 ssoudan and co.
                    </p>
                </FullWidthSection>

            </AppCanvas>
        );
    }

    _onLeftIconButtonTouchTap() {
        this.refs.leftNav.toggle();
    }
}

Master.contextTypes = {
    router: React.PropTypes.func
};

  
Master.childContextTypes = {
    muiTheme: React.PropTypes.object
};

module.exports = Master;
