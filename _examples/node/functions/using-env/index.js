
console.log('start using-env LOGGLY_TOKEN=%s', process.env.LOGGLY_TOKEN)
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({
    hello: 'bar',
    token_used: process.env.LOGGLY_TOKEN
  })
}
