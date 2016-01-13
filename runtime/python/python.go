package python

import (
	"github.com/apex/apex/runtime"
)

func init() {
	runtime.Register("python", new(Runtime))
}

type Runtime struct{}

func (r *Runtime) Name() string {
	return "python2.7"
}

func (r *Runtime) Handler() string {
	return "main.handle"
}

func (r *Runtime) Shimmed() bool {
	return false
}

func (r *Runtime) DefaultFile() string {
	return "main.py"
}
