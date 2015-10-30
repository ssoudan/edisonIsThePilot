var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    // devtool: 'source-map',
    entry: [
        './src/index'
    ],
    output: {
        path: path.join(__dirname, 'dist'),
        filename: 'bundle.js',
        publicPath: 'http://localhost:3000/'
    },
    // TODO(?) make a webpack config for prod: better packing, embedded roboto font,...
    plugins: [
        // new webpack.HotModuleReplacementPlugin(),
        new webpack.NoErrorsPlugin(),
        new HtmlWebpackPlugin({
            title: 'EdisonIsThePilot!',
            template: path.join(__dirname, 'assets/index-template.html'),
            "files": {
                "css": ["assets/main.css"]
            }
        })
    ],
    resolve: {
        extensions: ['', '.js', '.jsx']
    },
    module: {
        loaders: [
        {include: /\.js$/, loaders: ["babel-loader?stage=0&optional=runtime&plugins=typecheck"], exclude: /node_modules/},
        {
            test: /\.css$/,
            loader: "style!css"
        },
        { test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/, loader: "url-loader?limit=10000&minetype=application/font-woff" },
        { test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/, loader: "file-loader" }
        ]
    }
};
