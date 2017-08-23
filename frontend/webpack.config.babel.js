const CleanWebpackPlugin = require('clean-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

const isDev = process.env.NODE_ENV === 'development';
const isHot = isDev; // process.env.IS_HOT === 'yes'
const dist = `${__dirname}/dist`;

module.exports = {
  entry: (isHot ? [ // https://goo.gl/Zo2pVJ
    'react-hot-loader/patch',
    'webpack-dev-server/client',
    'webpack/hot/only-dev-server',
  ] : []).concat([
    './src/index',
  ]),
  output: {
    path: dist,
    publicPath: '/',
    filename: '[hash].js',
  },
  resolve: {
    extensions: ['*', '.js', '.jsx'], // https://goo.gl/cJKkwl
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: `babel-loader${isDev ? '!eslint-loader' : ''}`,
      }, {
        test: /\.scss$/,
        loaders: isDev ? [
          'style-loader?sourceMap',
          'css-loader?sourceMap&importLoaders=2',
          'postcss-loader?sourceMap',
          'sass-loader?sourceMap&sourceMapContents&outputStyle=expanded',
        ] : undefined,
        loader: isDev ? undefined : ExtractTextPlugin.extract({
          fallbackLoader: 'style',
          loader: 'css-loader?importLoaders=1!postcss-loader!sass-loader',
        }),
      }, {
        test: /\.css$/,
        loader: isDev ?
          'style-loader?sourceMap!css-loader?sourceMap' :
          'style-loader!css-loader',
      }, {
        test: /\.(png|ico|svg)$/,
        loader: 'file-loader',
      }, { // needed for react-markdown
        test: /\.json$/,
        loader: 'json-loader',
      },
    ],
  },
  plugins: [
    new CleanWebpackPlugin([dist]),
    new HtmlWebpackPlugin({
      template: 'index.html',
      title: 'auzom',
      favicon: 'favicon.ico',
      inject: true,
    }),
    new webpack.DefinePlugin({
      __DEV__: isDev,
      'process.env': {
        NODE_ENV: `"${process.env.NODE_ENV}"`,
      },
    }),
  ].concat(isHot ? [
    new webpack.HotModuleReplacementPlugin(),
  ] : []).concat(isDev ? [] : [
    new ExtractTextPlugin({
      filename: '[hash].css',
      allChunks: true,
    }),
    new webpack.LoaderOptionsPlugin({ minimize: true, debug: false }),
    // mangle is needed for CSS; remove when switching to CSS modules
    new webpack.optimize.UglifyJsPlugin({ mangle: false }),
  ]),
  // https://goo.gl/MgYrON
  devtool: isDev ? '#eval-cheap-module-source-map' : undefined,
  node: { fs: 'empty' },
};
