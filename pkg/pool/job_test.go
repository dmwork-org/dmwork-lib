package pool

import (
	"testing"
)

func TestJobType(t *testing.T) {
	executed := false
	job := &Job{
		Data: "test data",
		JobFunc: func(id int64, data interface{}) {
			executed = true
			if data != "test data" {
				t.Errorf("expected 'test data', got %v", data)
			}
		},
	}

	// Execute the job
	job.JobFunc(1, job.Data)

	if !executed {
		t.Error("job function was not executed")
	}
}

func TestJobWithNilData(t *testing.T) {
	job := &Job{
		Data: nil,
		JobFunc: func(id int64, data interface{}) {
			if data != nil {
				t.Errorf("expected nil, got %v", data)
			}
		},
	}

	job.JobFunc(1, job.Data)
}
