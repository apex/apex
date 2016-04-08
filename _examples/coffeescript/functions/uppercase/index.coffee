
console.log 'start bar'
exports.handle = (e, ctx, cb) ->
  console.log 'processing event: %j', e
  cb null, value: e.value.toUpperCase()
