package tasks

import (
	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger/api"
)

const (
	golangImage    = "golang:latest"
	baseImage      = "alpine:latest"
	publishAddress = "518461225764.dkr.ecr.us-east-1.amazonaws.com/greetings:latest"
)

func goBuilder(core *api.Query, ctx engine.Context, command []string) (*api.Container, error) {
	// Load image
	builder := core.Container().From(golangImage)

	// Set workdir
	src, err := core.Host().Workdir().Read().ID(ctx)
	if err != nil {
		return nil, err
	}
	builder = builder.WithMountedDirectory("/src", src).WithWorkdir("/src")
	builder = builder.WithEnvVariable("CGO_ENABLED", "0")
	builder = builder.WithEnvVariable("GOARCH", "amd64")
	builder = builder.WithEnvVariable("GOOS", "linux")

	// Execute Command
	builder = builder.Exec(api.ContainerExecOpts{
		Args: command,
	})
	return builder, nil
}
