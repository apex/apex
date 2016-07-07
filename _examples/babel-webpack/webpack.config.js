// NOTE: paths are relative to each functions folder

module.exports = {
  entry: './src/index.js',
  target: 'node',
  output: {
    path: './lib',
    filename: 'index.js',
    libraryTarget: 'commonjs2'
  },
  externals: {
    // aws-sdk does not (currently) build correctly with webpack. However,
    // Lambda includes it in its environment, so omit it from the bundle.
    // See: https://github.com/apex/apex/issues/217#issuecomment-194247472
    'aws-sdk': 'aws-sdk'
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel',
        exclude: [/node_modules/]
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      }
    ]
  }
}
