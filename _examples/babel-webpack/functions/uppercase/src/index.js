
export default function(e, ctx) {
  console.log('processing event: %j', e)
  ctx.succeed(e.value.toUpperCase())
}
