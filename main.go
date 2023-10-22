package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/zakisk/docker-clone/utils"
)

var (
	all bool
	digests bool
	filters []string
	noTrunc bool
	quiet bool
)

func main() {
	flag.BoolVar(&all, "all", false, "Show all images (default hides intermediate images)")

	flag.Parse()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Unabel to create docker client, please make sure that docker is installed\n%s", err.Error())
		os.Exit(1)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: all})
	if err != nil {
		log.Fatal("Unabel to get images, please make sure that docker daemon is up and running")
		os.Exit(1)
	}

	rows := [][]string{}
	for _, image := range images {
		repository := "<none>"
		tag := "<none>"
		if len(image.RepoTags) > 0 {
			splitted := strings.Split(image.RepoTags[0], ":")
			repository = splitted[0]
			tag = splitted[1]	
		} else if len(image.RepoDigests) > 0 {
			repository = strings.Split(image.RepoDigests[0], "@")[0]
		}
		duration := utils.HumanDuration(image.Created)
		size := utils.HumanSize(image.Size)
		rows = append(rows, []string{repository, tag, image.ID[7:19], duration, size})
	}
	header := []string{"Repository", "Tag", "Image Id", "Created", "Size"}
	utils.WriteToTable(header, rows)
}
