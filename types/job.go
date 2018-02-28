package types

// JobStatus represents the status of a job
type JobStatus string

// Category represents the category of a media file
type Category int

const (
	// IMAGE represents a image job
	IMAGE Category = iota
	// VIDEO represents a video job
	VIDEO
)

// These constants are used on the status field of Job type
const (
	JobCreated     = JobStatus("created")
	JobDownloading = JobStatus("downloading")
	JobProcessing  = JobStatus("processing")
	JobUploading   = JobStatus("uploading")
	JobFinished    = JobStatus("finished")
	JobError       = JobStatus("error")
)

// Job is the set of parameters of a given job
type Job struct {
	ID               string    `json:"id"`
	Source           string    `json:"source"`
	Destination      string    `json:"destination"`
	Media            MediaType `json:"mediaType"`
	Status           JobStatus `json:"status"`
	Details          string    `json:"details"`
	LocalSource      string    `json:"-"`
	LocalDestination string    `json:"-"`
}

// MediaType defines the type of output file
type MediaType struct {
	Cate      Category `json:"cate"`
	Name      string   `json:"name"`
	Container string   `json:"container"`
}

// JobInput stores the information passed from the
// user when creating a job.
type JobInput struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Cate        Category `json:"cate"`
}
