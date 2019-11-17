package models

import (
	"testing"
	"context"
	"flag"
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
)

var a App

func TestMain(m *testing.M) {
	var containerID string
	a = App{}

	flag.Parse()

	if !testing.Short() {
		fmt.Println("Not short!")

		err := godotenv.Load("../.env")
		if err != nil {
			fmt.Println("Failed to load the env file!")
		panic(err)
		}

		config := Config{}
		config.DBUser = os.Getenv("db_user")
		config.DBPassword = os.Getenv("db_pass")
		config.DBName = os.Getenv("db_name")
		config.DBHost = os.Getenv("db_host")
		config.DBPort = os.Getenv(("db_port"))

		config.TemplatePath = "../" + os.Getenv("template_path")

		a.SetConfig(&config)

		containerID, err = startDBinDocker(config.DBUser, config.DBPassword, config.DBName, config.DBPort, "")
		if err != nil {
			panic(err)
		}

		a.Initialize()
	}

	exitValue := m.Run()

	if !testing.Short() {
		defer a.GetDB().Close()
		fmt.Println("not short end!")

		err := stopDBinDocker(containerID)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(exitValue)
}

func startDBinDocker(user, password, name, port, volume string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	// reader, err := cli.ImagePull(context.Background(), "postgres/postgres", types.ImagePullOptions{})
	// if err != nil {
	// 	fmt.Println("Unable to pull image")
	// 	panic(err)
	// }
	// io.Copy(os.Stdout, reader)

	hostBinding := nat.PortBinding{HostIP: "0.0.0.0", HostPort: port}
	containerPort, err := nat.NewPort("tcp", port)
	if err != nil {
		panic("Unable to get the port")
	}
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	var binds []string
	if volume != "" {
		binds = []string{fmt.Sprintf("%s:/var/lib/postgresql/data", volume)}
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: "postgres/postgres",
		Env: []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", password)},
		Cmd: []string{"postgres"},
	}, &container.HostConfig{
		PortBindings: portBinding,
		Binds: binds,
	}, nil, "")
	if err != nil {
		fmt.Println("Unable to create docker container")
		panic(err)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID, nil
}

func stopDBinDocker(containerID string) (error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	return cli.ContainerStop(context.Background(), containerID, nil)
}