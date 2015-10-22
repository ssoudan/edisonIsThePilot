/** MyAwesomeReactComponent.jsx */
import React from 'react';
import Router from 'react-router';

import mui from 'material-ui';

var {MenuItem, LeftNav} = mui;
const {Colors, Spacing, Typography} = mui.Styles;

var menuItems = [
    {route: 'get-started', text: 'Get Started'},
    //{route: 'customization', text: 'Customization'},
    //{route: 'components', text: 'Components'},
    {type: MenuItem.Types.SUBHEADER, text: 'Resources'},
    {
        type: MenuItem.Types.LINK,
        payload: 'https://github.com/ssoudan/edisonIsThePilot',
        text: 'GitHub'
    },
    //{
    //    text: 'Disabled',
    //    disabled: true
    //},
    //{
    //    type: MenuItem.Types.LINK,
    //    payload: 'https://www.google.com',
    //    text: 'Disabled Link',
    //    disabled: true
    //}
];

export default class AppLeftNav extends React.Component {

    static contextTypes = {
        router: React.PropTypes.func
    };

    constructor() {
        super();

        this.toggle = this.toggle.bind(this);
        this._getSelectedIndex = this._getSelectedIndex.bind(this);
        this._onLeftNavChange = this._onLeftNavChange.bind(this);
        this._onHeaderClick = this._onHeaderClick.bind(this);
    }

    getStyles() {
        return {
            cursor: 'pointer',
            fontSize: '24px',
            color: Typography.textFullWhite,
            lineHeight: Spacing.desktopKeylineIncrement + 'px',
            fontWeight: Typography.fontWeightLight,
            backgroundColor: Colors.blueGrey500,
            paddingLeft: Spacing.desktopGutter,
            paddingTop: '0px',
            marginBottom: '8px'
        };
    }

    render() {

        var header = (
            <div style={this.getStyles()} onClick={this._onHeaderClick}>
                edisonIsThePilot
            </div>
        );

        return (
            <LeftNav
                ref="leftNav"
                docked={false}
                isInitiallyOpen={false}
                header={header}
                menuItems={menuItems}
                selectedIndex={this._getSelectedIndex()}
                onChange={this._onLeftNavChange}/>
        );
    }

    toggle() {
        this.refs.leftNav.toggle();
    }

    _getSelectedIndex() {
        var currentItem;
        for (var i = menuItems.length - 1; i >= 0; i--) {
            currentItem = menuItems[i];
            if (currentItem.route && this.context.router.isActive(currentItem.route)) return i;
        }
    }

    _onLeftNavChange(e, key, payload) {
        console.log("Going to " + payload.route);
        this.context.router.transitionTo(payload.route);
    }

    _onHeaderClick() {
        console.log("Going to root");
        this.context.router.transitionTo('root');
        this.refs.leftNav.close();
    }
};

AppLeftNav.contextTypes = {
  router: React.PropTypes.func
};

