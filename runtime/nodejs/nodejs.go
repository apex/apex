package nodejs

type Runtime struct{}

func (r *Runtime) Name() string {
	return "nodejs"
}

func (r *Runtime) Shimmed() bool {
	return false
}
