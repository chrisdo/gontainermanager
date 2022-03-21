package html

import (
	"embed"
	"html/template"
	"io"
)

//go:embed *
var files embed.FS

type ContainerData struct {
	Id     string
	Names  []string
	Image  string
	Status string
}

type ContainerDetails struct {
	Error  string
	Data   ContainerData
	Labels map[string]string
	IP     string
}

type ContainerList struct {
	Containers []ContainerData
}

func ListContainers(w io.Writer, cList ContainerList) error {
	tmpl, err := parse("containers.html") //later move parse to global to not parse it with each call
	if err != nil {
		return err
	}
	return tmpl.Execute(w, cList)
}

func InspectContainer(w io.Writer, details ContainerDetails) error {
	tmpl, err := parse("container.html") //later move parse to global to not parse it with each call
	if err != nil {
		return err
	}
	return tmpl.Execute(w, details)
}

func parse(file string) (*template.Template, error) {
	return template.New("layout.html").ParseFS(files, "layout.html", file)
}
