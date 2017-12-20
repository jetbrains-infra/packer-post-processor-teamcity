package main

import (
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/packer/plugin"
	"os"
	"strings"
	"fmt"
)

const TeamcityVersionEnvVar = "TEAMCITY_VERSION"
var AmazonBuilderIds = []string {
	"mitchellh.amazonebs",
	"mitchellh.amazon.ebssurrogate",
	"mitchellh.amazon.instance",
	"mitchellh.amazon.chroot",
}

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(PostProcessor))
	server.Serve()
}

type PostProcessor struct {
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	if os.Getenv(TeamcityVersionEnvVar) != "" {
		if contains(AmazonBuilderIds, artifact.BuilderId())  {
			s := strings.Split(artifact.Id(), ":")
			region, ami := s[0], s[1] // TODO: several AMIs
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.aws.region' value='%v']", region))
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.aws.ami' value='%v']", ami))
		} else {
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.id' value='%v']", artifact.Id()))
		}
	}
	return artifact, true, nil
}

func contains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}
