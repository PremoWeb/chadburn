package cli

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/PremoWeb/Chadburn/core"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DaemonCommand daemon process
type DaemonCommand struct {
	ConfigFile    string `long:"config" description:"configuration file" default:"/etc/chadburn.conf"`
	Metrics       bool   `long:"metrics" description:"Enable Prometheus compatible metrics endpoint"`
	MetricsAddr   string `long:"listen-address" description:"Metrics endpoint listen address." default:":8080"`
	DisableDocker bool   `long:"disable-docker" description:"Disable docker integration. All job kinds except 'job-local' will be ignored"`
	scheduler     *core.Scheduler
	signals       chan os.Signal
	done          chan bool
	Logger        core.Logger
}

// Execute runs the daemon
func (c *DaemonCommand) Execute(args []string) error {
	if err := c.boot(); err != nil {
		return err
	}

	if err := c.start(); err != nil {
		return err
	}

	if err := c.shutdown(); err != nil {
		return err
	}

	return nil
}

func (c *DaemonCommand) boot() (err error) {
	// Always try to read the config file, as there are options such as globals or some tasks that can be specified there and not in docker
	config, err := BuildFromFile(c.ConfigFile, c.Logger)
	if err != nil {
		c.Logger.Debugf("Config file: %v not found", c.ConfigFile)
	}

	err = config.InitializeApp(c.DisableDocker)
	if err != nil {
		c.Logger.Criticalf("Can't start the app: %v", err)
	}
	c.scheduler = config.sh

	return err
}

func startHttpServer(c *DaemonCommand, wg *sync.WaitGroup) *http.Server {
	c.Logger.Debugf("Starting metrics on %s", c.MetricsAddr)
	srv := &http.Server{Addr: c.MetricsAddr}
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		defer wg.Done()

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			c.Logger.Errorf("Metrics serving failed: %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func (c *DaemonCommand) start() error {
	if c.Metrics {
		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := startHttpServer(c, httpServerExitDone)
		c.setSignals(srv)
	} else {
		c.setSignals(nil)
	}

	if err := c.scheduler.Start(); err != nil {
		return err
	}

	return nil
}

func (c *DaemonCommand) setSignals(srv *http.Server) {
	c.signals = make(chan os.Signal, 1)
	c.done = make(chan bool, 1)

	signal.Notify(c.signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c.signals
		c.Logger.Warningf(
			"Signal received: %s, shutting down the process\n", sig,
		)
		if srv != nil {
			if err := srv.Shutdown(context.TODO()); err != nil {
				panic(err) // failure/timeout shutting down the server gracefully
			}
		}
		c.done <- true
	}()
}

func (c *DaemonCommand) shutdown() error {
	<-c.done
	if !c.scheduler.IsRunning() {
		return nil
	}

	c.Logger.Warningf("Waiting running jobs.")
	return c.scheduler.Stop()
}
