
Apex uses a Node.js shim for non-native language support. This is a very small program which executes in a child process, and feeds Lambda input through STDIN, and program output through STDOUT. Because of this STDERR must be used for logging, not STDOUT.
