var webpack = require('webpack');
var WebpackDevServer = require('webpack-dev-server');
var config = require('./webpack.config');

require("babel/register")({
    stage: 0,
    plugins: ["typecheck"]
});

new WebpackDevServer(webpack(config), {
    devtool: 'eval',
    publicPath: config.output.publicPath,
    hot: true,
    historyApiFallback: true,
    proxy: {
        "*": "http://localhost:8000/"
        //"*": "http://104.154.60.24:8080/"
    }
}).listen(3000, 'localhost', function (err, result) {
    if (err) {
        console.log(err);
    }
    console.log('Listening at localhost:3000');
});
