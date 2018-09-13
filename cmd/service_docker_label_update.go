package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/jpignata/fargate/console"
	ECS "github.com/jpignata/fargate/ecs"
	"github.com/spf13/cobra"
)

type ServiceDockerLabelUpdateOperation struct {
	ServiceName  string
	DockerLabels map[string]*string
}

func (o *ServiceDockerLabelUpdateOperation) Validate() {
	if len(o.DockerLabels) == 0 {
		console.IssueExit("No docker labels specified")
	}
}

func (o *ServiceDockerLabelUpdateOperation) SetDockerLabels() {
	fileData, err := ioutil.ReadFile("./Dockerfile")
	if err != nil {
		console.IssueExit("Failed to read Dockerfile, %v", err)
		return
	}

	dockerfile := string(fileData)

	// Parse Docker labels
	lines := strings.Split(dockerfile, "\n")
	if len(lines) == 0 {
		return
	}

	labels := map[string]*string{}

	for _, line := range lines {
		if strings.HasPrefix(line, "LABEL ") {
			trimmed := strings.Trim(line, "LABEL ")
			split := strings.SplitN(trimmed, "=", 2)
			if len(split) == 2 {
				value := strings.TrimPrefix(split[1], "\"")
				value = strings.TrimSuffix(value, "\"")

				labels[split[0]] = &value
			}
		}
	}

	o.DockerLabels = labels
}

var serviceDockerLabelsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update docker labels",
	Long: `Update docker labels

This command extracts the docker labels from the Dockerfile and attempts to 
update ECS with these labels.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceDockerLabelUpdateOperation{
			ServiceName: args[0],
		}

		operation.SetDockerLabels()
		operation.Validate()
		serviceDockerLabelsUpdate(operation)
	},
}

func init() {
	serviceDockerLabelCmd.AddCommand(serviceDockerLabelsUpdateCmd)
}

func serviceDockerLabelsUpdate(operation *ServiceDockerLabelUpdateOperation) {
	ecs := ECS.New(sess, clusterName)
	service := ecs.DescribeService(operation.ServiceName)
	taskDefinitionArn := ecs.AddDockerLabelsToECSTaskDefinition(service.TaskDefinitionArn, operation.DockerLabels)

	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)

	console.Info("Updated docker labels:", operation.ServiceName)

	for label, value := range operation.DockerLabels {
		console.Info("- %s=%s", label, *value)
	}

}
