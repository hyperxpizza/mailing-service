package job_pool

type JobStatus int64

const (
	Waiting JobStatus = iota
	Running
	Cancelled
	Done
	Failed
	Cooling
)

func (j JobStatus) String() string {
	var status string
	switch j {
	case Waiting:
		status = "Waiting"
	case Running:
		status = "Running"
	case Cancelled:
		status = "Cancelled"
	case Done:
		status = "Done"
	case Failed:
		status = "Failed"
	case Cooling:
		status = "Cooling"
	default:
		status = "Unknown"
	}

	return status
}

type JobType int64

const (
	Recurrent JobType = iota
	OneTime
)

func (j JobType) String() string {
	var jobType string
	switch j {
	case Recurrent:
		jobType = "Recurrent"
	case OneTime:
		jobType = "OneTime"
	default:
		jobType = "Unknown"
	}

	return jobType
}
