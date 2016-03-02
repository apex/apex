
console.log('starting function')
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({ hello: 'world' })
}