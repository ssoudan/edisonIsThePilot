/**
 * Created by ssoudan on 9/5/15.
 */

/**
 * Created by ssoudan on 9/5/15.
 */
import React from 'react';

let { Styles } = require('material-ui');
let FullWidthSection = require('../full-width-section');
let { Spacing } = Styles;

export default class GetStarted extends React.Component {

    static contextTypes = {
        muiTheme: React.PropTypes.object
    }

    constructor(props) {
        super(props);
        this.state = {
            data: []
        };
    }

    getStyles() {
        return {
            root: {
                paddingTop: Spacing.desktopKeylineIncrement
            },
            fullWidthSection: {
                maxWidth: '650px',
                margin: '0 auto'
            }
        };
    }

    render() {

        let styles = this.getStyles();

        // TODO(?) make content for this page

        return (
            <div style={styles.root}>
                <FullWidthSection style={styles.FullWidthSection}>
                    Blah
                </FullWidthSection>
            </div>);
    }
}

