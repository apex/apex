
console.log('start bar LOGGLY_TOKEN=%s', process.env.LOGGLY_TOKEN)
exports.handle = function(e, ctx, cb) {
  console.log('processing event: %j', e)
  cb(null, { hello: 'bar' })
}
