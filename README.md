# Packer Post-Processor for TeamСity

HashiCorp Packer generates image IDs, like `ami-387dc380`.
Having a build chain with separate configurations for building, testing, and deploying images, we need to reference those IDs.

This is a Packer plugin, which saves generated image ID as a parameter in a build history:
![packer.artifact.id](docs/parameters.png)

For AMIs in AWS, there is a special format:
![packer.artifact.aws.ami and packer.artifact.aws.region](docs/parameters-aws.png)

Now dependent build configurations can reference these parametes, and resolve IDs dynamically:
![reference](docs/reference.png)

## Usage

1. Download binaries from [Releases](https://github.com/JetBrains/packer-post-processor-teamcity/releases) page.
2. [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin to build agents.
3. Add post-processor to Packer configurations:
```json
{
  "builders": [
    ...
  ],
  "post-processors": [
    "teamcity"
  ]
}
```
