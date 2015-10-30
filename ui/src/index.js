import React from 'react';
import Router from 'react-router';
import AppRoutes from './app-routes.js';

require("!style!css!./main.css");
require("babel/register")({
  optional: ['es7.asyncFunctions']
});
require("font-awesome-webpack");

let injectTapEventPlugin = require("react-tap-event-plugin");

//Needed for React Developer Tools
window.React = React;

//Needed for onTouchTap
//Can go away when react 1.0 release
//Check this repo:
//https://github.com/zilverline/react-tap-event-plugin
injectTapEventPlugin();

// TODO(?) need an favicon

Router
// Runs the router, similiar to the Router.run method. You can think of it as an
// initializer/constructor method.
    .create({
        routes: AppRoutes,
        scrollBehavior: Router.ScrollToTopBehavior
    })
    // This is our callback function, whenever the url changes it will be called again.
    // Handler: The ReactComponent class that will be rendered
    .run(function (Handler) {
        // google.maps is ready
                React.render(<Handler/>, document.body);
    });


