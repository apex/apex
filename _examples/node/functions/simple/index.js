
console.log('start simple')
exports.handle = function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed({ hello: e.hello })
}
