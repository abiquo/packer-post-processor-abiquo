package main

import (
	"fmt"
)

const BuilderId = "packer.post-processor.abiquo"

type Artifact struct {
	Url string
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (*Artifact) Id() string {
	return ""
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) String() string {
	return fmt.Sprintf("URL of the template : %s", a.Url)
}

func (*Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	// TODO
	// Delete the template from Abiquo
	return nil
}
