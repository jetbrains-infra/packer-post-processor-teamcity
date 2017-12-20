package main

import (
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
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
	config common.PackerConfig
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, nil, raws...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	if os.Getenv(TeamcityVersionEnvVar) != "" {
		if contains(AmazonBuilderIds, artifact.BuilderId())  {
			s := strings.Split(artifact.Id(), ":")
			region, ami := s[0], s[1] // TODO: several AMIs
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.%v.aws.region' value='%v']", p.config.PackerBuildName, region))
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.%v.aws.ami' value='%v']", p.config.PackerBuildName, ami))
		} else {
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.%v.id' value='%v']", p.config.PackerBuildName, artifact.Id()))
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
