
console.log 'start bar'
exports.handle = (e, ctx) ->
  console.log 'processing event: %j', e
  ctx.succeed value: e.value.toUpperCase()
