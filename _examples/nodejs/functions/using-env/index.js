
console.log('start using-env LOGGLY_TOKEN=%s', process.env.LOGGLY_TOKEN)
exports.handle = function(e, ctx, cb) {
  console.log('processing event: %j', e)
  cb(null, {
    hello: 'bar',
    token_used: process.env.LOGGLY_TOKEN
  })
}
