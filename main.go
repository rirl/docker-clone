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
	"github.com/docker/docker/api/types/filters"
)

var (
	all     bool
	digests bool
	filter stringSlice
	noTrunc bool
	quiet   bool
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	flag.BoolVar(&all, "all", false, "Show all images (default hides intermediate images)")
	flag.BoolVar(&digests, "digests", false, "Show digests")
	flag.Var(&filter, "filter", "Filter output based on conditions provided")

	flag.Parse()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Unabel to create docker client, please make sure that docker is installed\n%s", err.Error())
		os.Exit(1)
	}
	options := types.ImageListOptions{All: all}
	filterArgs := filters.NewArgs()
	for _, filt := range filter {
		key, value, err := utils.FormatFilter(filt)
		if err != nil {
			log.Fatal(err)
		}
		filterArgs.Add(key, value)
	}
	options.Filters = filterArgs

	images, err := cli.ImageList(context.Background(), options)
	if err != nil {
		log.Fatal("Unabel to get images, please make sure that docker daemon is up and running")
	}

	rows := [][]string{}
	for index, image := range images {
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

		if digests {
			digest := "<none>"
			if len(image.RepoDigests) > 0 {
				digest = strings.Split(image.RepoDigests[0], "@")[1]
			}
			rows[index] = append(rows[index][:3], rows[index][2:]...)
			rows[index][2] = digest
		}

	}
	header := []string{"Repository", "Tag", "Image Id", "Created", "Size"}

	// if digests flag is passed add digests column in table
	if digests {
		header = append(header[:3], header[2:]...)
		header[2] = "Digests"
	}
	utils.WriteToTable(header, rows)
}
