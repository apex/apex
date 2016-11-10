// NOTE: paths are relative to each functions folder
const Webpack = require('webpack');

module.exports = {
  entry: './src/index.js',
  target: 'node',
  output: {
    path: './lib',
    filename: 'index.js',
    libraryTarget: 'commonjs2'
  },
  externals: {
    'aws-sdk': 'aws-sdk'
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel',
        query: {
          presets: [
            'stage-0',
            'latest'
          ],
          plugins: [
            'transform-promise-to-bluebird',
            ['transform-async-to-module-method', {
              module: 'bluebird',
              method: 'coroutine',
            }],
            'transform-runtime',
          ],
          cacheDirectory: true,
        },
        exclude: [/node_modules/]
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      }
    ],
  },
  plugins: [
    new Webpack.LoaderOptionsPlugin({
      minimize: true,
      debug: false,
    }),
    new Webpack.optimize.UglifyJsPlugin({
      compress: { warnings: false },
      output: {
        comments: false,
      },
      mangle: false,
    }),
  ]
}
