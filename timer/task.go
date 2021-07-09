package timer

type Task struct {
	fun func(...interface{})
	args []interface{}
}

func (t *Task) Call()  {
	t.fun(t.args...)
}

func NewTask(fun func(...interface{}), args []interface{}) *Task{
	return &Task{
		fun:  fun,
		args: args,
	}
}