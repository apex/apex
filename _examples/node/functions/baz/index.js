
console.log('start baz')
exports.handle = function(e, ctx) {
  ctx.succeed({ hello: 'baz' })
}
