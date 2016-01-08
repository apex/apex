
console.log('start foo')
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({ hello: 'foo' })
}
