
console.log('start bar')
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({ hello: 'bar' })
}
