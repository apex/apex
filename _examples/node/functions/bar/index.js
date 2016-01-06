
console.log('start bar')
exports.handle = function(e, ctx) {
  ctx.succeed({ hello: 'bar' })
}
