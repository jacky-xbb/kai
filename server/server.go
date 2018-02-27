package server

import (
	"net"
	"net/http"
	"os"

	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// KaiServer represents the server for processing all job requests
type KaiServer struct {
	net.Listener
	logger        *logging.Logger
	config        types.ServerConfig
	listenAddr    string
	listenNetwork string
	router        *Router
	server        *http.Server
	db            db.Storage
	// jobDownBuff is the buffered channel for job downloading
	jobDownBuff chan types.Job
	// jobDownBuff is the buffered channel for job todo
	jobTodoBuff chan types.Job
	// jobDownBuff is the buffered channel for job done
	jobDoneBuff chan types.Job
}

// New creates a kai server
func New(log *logging.Logger, config types.ServerConfig, listenNetwork string, listenAddr string, db db.Storage) *KaiServer {
	ks := &KaiServer{
		logger:        log,
		listenAddr:    listenAddr,
		listenNetwork: listenNetwork,
		router:        NewRouter(),
		config:        config,
		db:            db,
	}

	ks.logger.Debug("setting-up-routes")
	// Set up routes
	routes := map[Route]RouterArguments{
		CreateJob:     {Path: Routes[CreateJob].Path, Method: Routes[CreateJob].Method, Handler: ks.CreateJob},
		StartJob:      {Path: Routes[StartJob].Path, Method: Routes[StartJob].Method, Handler: ks.StartJob},
		ListJobs:      {Path: Routes[ListJobs].Path, Method: Routes[ListJobs].Method, Handler: ks.ListJobs},
		GetJobDetails: {Path: Routes[GetJobDetails].Path, Method: Routes[GetJobDetails].Method, Handler: ks.GetJobDetails},
		DeleteJob:     {Path: Routes[DeleteJob].Path, Method: Routes[DeleteJob].Method, Handler: ks.DeleteJob},
	}
	for _, route := range routes {
		ks.router.AddHandler(RouterArguments{Path: route.Path, Method: route.Method, Handler: route.Handler})
	}

	bSize := config.System.Workers
	// set twice workers for job buffered channels
	ks.jobDownBuff = make(chan types.Job, 2*bSize)
	ks.jobTodoBuff = make(chan types.Job, bSize)
	ks.jobDoneBuff = make(chan types.Job, bSize)

	ks.server = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ks.router.r.ServeHTTP(w, r)
		}),
	}

	return ks
}

// Handler returns the router of the kai server
func (ks *KaiServer) Handler() http.Handler {
	return ks.router.Handler()
}

// Start starts the kai server
func (ks *KaiServer) Start(keep bool) error {
	var err error

	ks.Listener, err = net.Listen(ks.listenNetwork, ks.listenAddr)
	if err != nil {
		ks.logger.Error("kai-failed-starting-server", err)
		return err
	}

	if keep {
		ks.logger.Info("started")
		ks.server.Serve(ks.Listener)
		return nil
	}

	go ks.server.Serve(ks.Listener)
	ks.logger.Info("started")
	return nil
}

func (ks *KaiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ks.router.Handler()
}

// Stop stops the kai server
func (ks *KaiServer) Stop() error {
	ks.logger.Info("stop-server")
	defer ks.logger.Info("stop")

	if ks.listenNetwork == "unix" {
		if err := os.Remove(ks.listenAddr); err != nil {
			ks.logger.Infof("failed-to-stop-server listenAddr: %s", ks.listenAddr)
			return err
		}
	}

	return ks.Listener.Close()
}
