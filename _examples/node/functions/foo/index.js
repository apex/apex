
console.log('start foo')
exports.handle = function(e, ctx) {
  ctx.succeed({ hello: 'foo' })
}
