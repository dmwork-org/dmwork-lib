package pool

// JobFunc is the function signature for job execution
type JobFunc func(id int64, data interface{})

// Job represents a unit of work to be processed by the pool
type Job struct {
	Data    interface{}
	JobFunc JobFunc
}
