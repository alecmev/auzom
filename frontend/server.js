const webpack = require('webpack');
const WebpackDevServer = require('webpack-dev-server');

const config = require('./webpack.config.babel');

new WebpackDevServer(webpack(config), {
  historyApiFallback: true,
  hot: true,
  stats: {
    chunks: false,
    colors: true,
    hash: false,
    timings: false,
    version: false,
  },
}).listen(80, '0.0.0.0', (err) => {
  if (err) console.log(err);
});
