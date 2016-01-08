
console.log('start baz')
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({ hello: 'baz' })
}
