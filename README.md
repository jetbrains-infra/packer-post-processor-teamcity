# Packer Post-Processor for TeamСity

HashiCorp Packer generates image IDs, like `ami-387dc380`.
When creating a build chain with separate configurations for building, testing, and deploying images, we need to reference their IDs.

This is a Packer plugin, which saves the generated image ID as a parameter in the build history:
![packer.artifact.id](docs/parameters.png)

For AMIs in AWS, there is a special format:
![packer.artifact.aws.ami and packer.artifact.aws.region](docs/parameters-aws.png)

Now dependent build configurations can reference these parametes and resolve the IDs dynamically:
![reference](docs/reference.png)

## Usage

1. Download the binaries from the [Releases](https://github.com/JetBrains/packer-post-processor-teamcity/releases) page.
2. [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin on build agents.
3. Add the TeamCity post-processor to Packer configurations:
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
