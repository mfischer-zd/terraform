package docker

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_HOST", "unix:///var/run/docker.sock"),
				Description: "The Docker daemon address",
			},

			"ca_material": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("DOCKER_CA_MATERIAL", ""),
				ConflictsWith: []string{"cert_path"},
				Description:   "PEM-encoded content of Docker host CA certificate",
			},
			"cert_material": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("DOCKER_CERT_MATERIAL", ""),
				ConflictsWith: []string{"cert_path"},
				Description:   "PEM-encoded content of Docker client certificate",
			},
			"key_material": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("DOCKER_KEY_MATERIAL", ""),
				ConflictsWith: []string{"cert_path"},
				Description:   "PEM-encoded content of Docker client private key",
			},

			"cert_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_CERT_PATH", ""),
				Description: "Path to directory with Docker TLS config",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"docker_container": resourceDockerContainer(),
			"docker_image":     resourceDockerImage(),
			"docker_network":   resourceDockerNetwork(),
			"docker_volume":    resourceDockerVolume(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"docker_registry_image": dataSourceDockerRegistryImage(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:     d.Get("host").(string),
		Ca:       d.Get("ca_material").(string),
		Cert:     d.Get("cert_material").(string),
		Key:      d.Get("key_material").(string),
		CertPath: d.Get("cert_path").(string),
	}

	client, err := config.NewClient()
	if err != nil {
		return nil, fmt.Errorf("Error initializing Docker client: %s", err)
	}

	err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging Docker server: %s", err)
	}

	return client, nil
}
