package server

import "net/http"

// Route defines all job routes
type Route int

const (
	// CreateJob is the route for creating
	CreateJob Route = iota
	// StartJob is the route for starting
	StartJob
	// ListJobs is the route for listing
	ListJobs
	// GetJobDetails is the route for getting
	GetJobDetails
	// DeleteJob is the route for deleting
	DeleteJob
)

// Routes defines a map between the routes and the routes' argument
var Routes = map[Route]RouterArguments{
	CreateJob:     RouterArguments{Path: "/jobs", Method: http.MethodPost},
	StartJob:      RouterArguments{Path: "/jobs/{jobID}/start", Method: http.MethodPost},
	ListJobs:      RouterArguments{Path: "/jobs", Method: http.MethodGet},
	GetJobDetails: RouterArguments{Path: "/jobs/{jobID}", Method: http.MethodGet},
	DeleteJob:     RouterArguments{Path: "/jobs/{jobID}", Method: http.MethodDelete},
}
