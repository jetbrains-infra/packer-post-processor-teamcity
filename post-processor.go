//go:generate mapstructure-to-hcl2 -type Config

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/packer/plugin"
	"net/http"
	"os"
	"strings"
)

var AmazonBuilderIds = []string{
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
	config Config
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	TeamCityUrl string `mapstructure:"teamcity_url"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	ProjectId   string `mapstructure:"project_id"`
	CloudImage  string `mapstructure:"cloud_image"`
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, nil, raws...)
	if err != nil {
		return err
	}

	errs := new(packer.MultiError)
	if p.config.TeamCityUrl != "" {
		if p.config.Username == "" {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("username is required"))
		}
		if p.config.Password == "" {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("password is required"))
		}
		if p.config.ProjectId == "" {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("project_id is required"))
		}
		if p.config.CloudImage == "" {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("cloud_image is required"))
		}
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (a packer.Artifact, keep bool, forceOverride bool, err error) {
	isAmazonArtifact := contains(AmazonBuilderIds, artifact.BuilderId())
	var image string
	if isAmazonArtifact {
		s := strings.Split(artifact.Id(), ":")
		image = s[1]
	} else {
		image = artifact.Id()
	}

	if os.Getenv("TEAMCITY_VERSION") != "" {
		ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.%v.id' value='%v']", p.config.PackerBuildName, image))

		if isAmazonArtifact {
			s := strings.Split(artifact.Id(), ":")
			region, ami := s[0], s[1]
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.%v.aws.region' value='%v']", p.config.PackerBuildName, region))
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.%v.aws.ami' value='%v']", p.config.PackerBuildName, ami))
		} else {
			ui.Message(fmt.Sprintf("##teamcity[setParameter name='packer.artifact.last.id' value='%v']", image))
		}
	}

	if p.config.TeamCityUrl != "" {
		var url string
		if isAmazonArtifact {
			url = fmt.Sprintf(
				"%v/httpAuth/app/rest/projects/id:%v/projectFeatures/type:CloudImage,property(name:image-name-prefix,value:%v)/properties/amazon-id",
				strings.TrimRight(p.config.TeamCityUrl, "/"),
				p.config.ProjectId,
				p.config.CloudImage,
			)
		} else {
			url = fmt.Sprintf(
				"%v/httpAuth/app/rest/projects/id:%v/projectFeatures/type:CloudImage,property(name:source-id,value:%v)/properties/sourceVmName",
				strings.TrimRight(p.config.TeamCityUrl, "/"),
				p.config.ProjectId,
				p.config.CloudImage,
			)
		}

		body := bytes.NewBufferString(image)

		c := &http.Client{}
		req, err := http.NewRequest("PUT", url, body)
		if err != nil {
			return artifact, true, false, err
		}
		req.Header.Add("Content-Type", "text/plain")
		req.SetBasicAuth(p.config.Username, p.config.Password)

		resp, err := c.Do(req)
		if err != nil {
			return artifact, true, false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return artifact, true, false, errors.New(fmt.Sprintf("Error updating a cloud profile: %v", resp.Status))
		}

		ui.Message(fmt.Sprintf("Cloud agent image '%v' is switched to image '%v'", p.config.CloudImage, image))
	}

	return artifact, true, false, nil
}

func contains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}
