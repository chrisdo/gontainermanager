package main

import (
	"context"
	"fmt"
	"gontainermanager/html"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type containerManager struct {
	dockerClient *client.Client
	ctx          context.Context
}

func main() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	manager := &containerManager{cli, ctx}

	mux := http.NewServeMux()
	mux.HandleFunc("/containers", listContainers(manager))
	mux.HandleFunc("/containers/inspect", inspectContainer(manager))

	log.Println("ready")
	http.ListenAndServe(":8080", mux)

}

//listContainers handles /containers and will list all running containers
func listContainers(mgr *containerManager) func(w http.ResponseWriter, r *http.Request) {

	containers, err := mgr.dockerClient.ContainerList(mgr.ctx, types.ContainerListOptions{Latest: true})
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			fmt.Fprint(w, err)
		}
		containerList := make([]html.ContainerData, len(containers))
		for i, c := range containers {
			containerList[i] = html.ContainerData{Id: shortenId(c.ID), Names: c.Names, Image: c.Image, Status: c.Status}
		}
		log.Printf("%v", containerList)
		html.ListContainers(w, html.ContainerList{Containers: containerList})
	}
}

func shortenId(id string) string {
	return id[:12]
}

//inspectContainer is the request handler for /containers/inspect?cid=<containerId>
//if cid is missing or no ontainer with specified ID can be found, an error is returned
func inspectContainer(mgr *containerManager) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		cid := r.URL.Query().Get("cid")
		if len(cid) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Empty ID not allowed"))
			return
		}
		log.Printf("Container ID: %s", cid)
		container, err := mgr.dockerClient.ContainerInspect(mgr.ctx, cid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			html.InspectContainer(w, html.ContainerDetails{Error: "No Container Data found for id " + cid})
			return
		}
		names := make([]string, 1)
		names[0] = container.Name
		html.InspectContainer(w, html.ContainerDetails{Data: html.ContainerData{Id: container.ID, Names: names, Image: container.Image, Status: container.State.Status},
			Labels: container.Config.Labels, IP: container.NetworkSettings.IPAddress})
	}
}
