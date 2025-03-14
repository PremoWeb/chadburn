package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PremoWeb/Chadburn/core"
	"github.com/mitchellh/mapstructure"
)

const (
	labelPrefix    = "chadburn"
	LabelNamespace = "chadburn"

	requiredLabel       = labelPrefix + ".enabled"
	requiredLabelFilter = requiredLabel + "=true"
	serviceLabel        = labelPrefix + ".service"
)

func (c *Config) buildFromDockerLabels(labels map[string]map[string]string) error {
	execJobs := make(map[string]map[string]interface{})
	localJobs := make(map[string]map[string]interface{})
	runJobs := make(map[string]map[string]interface{})
	serviceJobs := make(map[string]map[string]interface{})
	lifecycleJobs := make(map[string]map[string]interface{})
	globalConfigs := make(map[string]interface{})

	for c, l := range labels {
		isServiceContainer := func() bool {
			for k, v := range l {
				if k == serviceLabel {
					return v == "true"
				}
			}
			return false
		}()

		for k, v := range l {
			parts := strings.Split(k, ".")
			if len(parts) < 4 {
				if isServiceContainer {
					globalConfigs[parts[1]] = v
				}

				continue
			}

			jobType, jobName, jopParam := parts[1], parts[2], parts[3]
			switch {
			case jobType == jobExec: // only job exec can be provided on the non-service container
				if _, ok := execJobs[jobName]; !ok {
					execJobs[jobName] = make(map[string]interface{})
					execJobs[jobName]["fromDockerLabel"] = true
				}

				setJobParam(execJobs[jobName], jopParam, v)
				// since this label was placed not on the service container
				// this means we need to `exec` command in this container
				if !isServiceContainer {
					execJobs[jobName]["container"] = c
				}
			case jobType == jobLocal && isServiceContainer:
				if _, ok := localJobs[jobName]; !ok {
					localJobs[jobName] = make(map[string]interface{})
					localJobs[jobName]["fromDockerLabel"] = false
				}
				setJobParam(localJobs[jobName], jopParam, v)
			case jobType == jobServiceRun && isServiceContainer:
				if _, ok := serviceJobs[jobName]; !ok {
					serviceJobs[jobName] = make(map[string]interface{})
				}
				setJobParam(serviceJobs[jobName], jopParam, v)
			case jobType == jobRun:
				if _, ok := runJobs[jobName]; !ok {
					runJobs[jobName] = make(map[string]interface{})
				}
				setJobParam(runJobs[jobName], jopParam, v)
				// If the label is on a non-service container, set the container name
				if !isServiceContainer {
					runJobs[jobName]["container"] = c
				}
			case jobType == jobLifecycle:
				if _, ok := lifecycleJobs[jobName]; !ok {
					lifecycleJobs[jobName] = make(map[string]interface{})
					lifecycleJobs[jobName]["fromDockerLabel"] = true
				}
				setJobParam(lifecycleJobs[jobName], jopParam, v)
				// Set the container name for the lifecycle job
				lifecycleJobs[jobName]["container"] = c
			default:
				// TODO: warn about unknown parameter
			}
		}
	}

	if len(globalConfigs) > 0 {
		if err := mapstructure.WeakDecode(globalConfigs, &c.Global); err != nil {
			return err
		}
	}

	if len(execJobs) > 0 {
		if err := mapstructure.WeakDecode(execJobs, &c.ExecJobs); err != nil {
			return err
		}
	}

	if len(localJobs) > 0 {
		if err := mapstructure.WeakDecode(localJobs, &c.LocalJobs); err != nil {
			return err
		}
	}

	if len(serviceJobs) > 0 {
		if err := mapstructure.WeakDecode(serviceJobs, &c.ServiceJobs); err != nil {
			return err
		}
	}

	if len(runJobs) > 0 {
		if err := mapstructure.WeakDecode(runJobs, &c.RunJobs); err != nil {
			return err
		}
	}

	if len(lifecycleJobs) > 0 {
		if err := mapstructure.WeakDecode(lifecycleJobs, &c.LifecycleJobs); err != nil {
			return err
		}
	}

	return nil
}

func setJobParam(params map[string]interface{}, paramName, paramVal string) {
	switch paramName {
	case "volume":
		arr := []string{} // allow providing JSON arr of volume mounts
		if err := json.Unmarshal([]byte(paramVal), &arr); err == nil {
			params[paramName] = arr
			return
		}
	}

	params[paramName] = paramVal
}

func buildFromDockerLabels(client core.DockerClient, labels map[string]string, container string) (map[string]core.Job, error) {
	jobs := make(map[string]core.Job, 0)
	prefix := fmt.Sprintf("%s.", LabelNamespace)
	prefixLen := len(prefix)

	containerId := container
	containerName := container
	if strings.Contains(container, ":") {
		parts := strings.Split(container, ":")
		containerId = parts[0]
		containerName = parts[1]
	}

	for key, _ := range labels {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		jobType, jobName, field := extractFromKey(key[prefixLen:])
		if field == "schedule" {
			job, err := buildJob(client, jobType, jobName, containerId, containerName)
			if err != nil {
				return nil, err
			}

			job.SetCronJobID(0)
			jobs[jobName] = job
		}
	}

	for key, value := range labels {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		_, jobName, field := extractFromKey(key[prefixLen:])
		if field == "schedule" {
			continue
		}

		if _, ok := jobs[jobName]; !ok {
			continue
		}

		if err := updateJobField(jobs[jobName], field, value); err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func buildJob(client core.DockerClient, jobType, jobName, containerId, containerName string) (core.Job, error) {
	switch jobType {
	case "job-exec":
		j := client.CreateExecJob()
		j.Container = containerId
		return j, nil
	case "job-run":
		j := client.CreateRunJob()
		// RunJob doesn't have Container field, but uses Image
		// We'll set it during job configuration
		return j, nil
	case "job-local":
		j := core.NewLocalJob()
		j.ContainerName = containerName
		j.ContainerID = containerId
		return j, nil
	case "service":
		j := client.CreateServiceJob()
		j.Container = containerId
		return j, nil
	case "job-lifecycle":
		j := client.CreateLifecycleJob()
		j.Container = containerId
		return j, nil
	default:
		return nil, fmt.Errorf("unknown job type: %s", jobType)
	}
}

func extractFromKey(key string) (string, string, string) {
	parts := strings.Split(key, ".")
	if len(parts) < 3 {
		return "", "", ""
	}

	return parts[0], parts[1], parts[2]
}

func updateJobField(job core.Job, field, value string) error {
	switch field {
	case "volume":
		arr := []string{}
		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return err
		}
		job.SetVolumeMounts(arr)
	case "event-type":
		// Handle event-type field for lifecycle jobs
		if lifecycleJob, ok := job.(*core.LifecycleJob); ok {
			switch value {
			case "start":
				lifecycleJob.EventType = core.ContainerStart
			case "stop":
				lifecycleJob.EventType = core.ContainerStop
			default:
				return fmt.Errorf("unknown event type: %s", value)
			}
		} else if officialLifecycleJob, ok := job.(*core.OfficialLifecycleJob); ok {
			switch value {
			case "start":
				officialLifecycleJob.EventType = core.ContainerStart
			case "stop":
				officialLifecycleJob.EventType = core.ContainerStop
			default:
				return fmt.Errorf("unknown event type: %s", value)
			}
		}
	case "container":
		// Handle container field for run jobs
		if runJob, ok := job.(*core.OfficialRunJob); ok {
			// For run jobs, we use the container value as the image
			runJob.Image = value
		}
	default:
		// TODO: handle other fields
	}

	return nil
}
