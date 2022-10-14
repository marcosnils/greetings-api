package tasks

import (
	"context"
	"fmt"

	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger/api"
)

func Push(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context) error {
		core := api.New(ctx.Client)

		builder, err := goBuilder(
			core,
			ctx,
			[]string{"go", "build"},
		)
		if err != nil {
			return err
		}

		// Get built binary
		greetingsBin, err := builder.File("/src/greetings-api").ID(ctx)
		if err != nil {
			return err
		}

		// Get base image for publishing
		base := core.Container().From(baseImage)
		// Add built binary to /bin
		base = base.WithMountedFile("/tmp/greetings-api", greetingsBin)
		// Copy mounted file to rootfs
		base = base.Exec(api.ContainerExecOpts{
			Args: []string{"cp", "/tmp/greetings-api", "/bin/greetings-api"},
		})
		// Set entrypoint
		base = base.WithEntrypoint([]string{"/bin/greetings-api"})
		// Publish image
		addr, err := base.Publish(ctx, publishAddress)
		if err != nil {
			return err
		}

		fmt.Println(addr)

		// Create ECS task deployment
		err = deployGreetingsService()
		if err != nil {
			return err
		}
		fmt.Println("Deployed ECS task")

		return nil
	}); err != nil {
		panic(err)
	}
}