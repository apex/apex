
var child = require('child_process')
var byline = require('./byline')

/**
 * Callback for the request.
 */
var cb

/**
 * Child process for binary I/O.
 */

var proc = child.spawn('./main', { stdio: ['pipe', 'pipe', process.stderr] })

proc.on('error', function(err){
  console.error('error: %s', err)
  process.exit(1)
})

proc.on('exit', function(code){
  console.error('exit: %s', code)
  process.exit(1)
})

/**
 * Newline-delimited JSON stdout.
 */

var out = byline(proc.stdout)

out.on('data', function(line){
  if (process.env.DEBUG_SHIM) console.log('[shim] parsing: %j', line)
  var msg = JSON.parse(line)
  cb(msg.error, msg.value)
})

/**
 * Handle events.
 */

exports.handle = function(event, context, callback) {
  cb = callback
  context.callbackWaitsForEmptyEventLoop = false

  proc.stdin.write(JSON.stringify({
    "event": event,
    "context": context
  })+'\n');
}
