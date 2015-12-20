package python

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
