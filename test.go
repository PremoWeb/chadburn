package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*
	{
	Status:destroy
	ID:f1a09c48aa893fb68388f1ebeda0ea12dc4df7bcba9a98abab3abb9cee06e9e5
	From:nginx
	Type:container
	Action:destroy
	Actor: {
		ID:f1a09c48aa893fb68388f1ebeda0ea12dc4df7bcba9a98abab3abb9cee06e9e5
		Attributes:map[
			chadburn.enabled:true
			chadburn.job-exec.test-exec-job.command:uname -a
			chadburn.job-exec.test-exec-job.schedule:@every 5s
			image:nginx maintainer:NGINX Docker Maintainers <docker-maint@nginx.com>
			name:festive_tharp
		]
		}
	Scope:local
	Time:1641000210
	TimeNano:1641000210626965000
*/

func main() {

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	messages, errs := cli.Events(context.Background(), types.EventsOptions{})

EventLoop:
	for {
		select {
		case err := <-errs:
			if err != nil && err != io.EOF {
				panic(err)
			}
			if err != nil && err == io.EOF {
				fmt.Printf("lost connection with Docker. exiting\n")
				break EventLoop
			}
		case message := <-messages:

			if message.Action == "start" || message.Action == "die" {

				container_name := message.Actor.Attributes["name"]
				fmt.Println("Registered:", container_name)

				//	fmt.Printf("%+v\n\n", message)

				// Get service labels from the event
				enabled := message.Actor.Attributes["chadburn.enabled"]

				if enabled == "true" {

					exec_type := ""
					exec_cmd := ""

					// Get attributes that start with chadburn.job-*
					for k, v := range message.Actor.Attributes {
						if strings.HasPrefix(k, "chadburn.job-") {
							exec_type = k
							exec_cmd = v
							fmt.Println(exec_type, "", exec_cmd)
						}
					}

					fmt.Printf("%+v\n\n", message)

					if message.Actor.Attributes["chadburn.job-exec"] == "true" {
						// This job is executed inside of a running container.
					} else if message.Actor.Attributes["chadburn.job-run"] == "true" {
						// Runs a command inside of a new container, using a specific image.
					} else if message.Actor.Attributes["chadburn.job-local"] == "true" {
						// Runs the command inside of the host running Chadburn.
					} else if message.Actor.Attributes["chadburn.job-service-run"] == "true" {
						// Runs the command inside a new "run-once" service, for running inside a swarm
					}

				}

			}

			if message.Action == "stop" || message.Action == "die" || message.Action == "destroy" {
				// Dereigster the service

			}

		}
	}

}
