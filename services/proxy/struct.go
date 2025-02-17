package proxy

import (
	"github.com/ironstar-io/tokaido/constants"
)

// DockerCompose ...
type DockerCompose struct {
	Version  string `yaml:"version"`
	Services struct {
		Unison struct {
			Hostname    string   `yaml:"hostname,omitempty"`
			Entrypoint  []string `yaml:"entrypoint,omitempty"`
			User        string   `yaml:"user,omitempty"`
			Cmd         string   `yaml:"cmd,omitempty"`
			Volumesfrom []string `yaml:"volumes_from,omitempty"`
			Dependson   []string `yaml:"depends_on,omitempty"`
			Image       string   `yaml:"image"`
			Environment []string `yaml:"environment"`
			Ports       []string `yaml:"ports"`
			Volumes     []string `yaml:"volumes"`
			Networks    []string `yaml:"networks"`
		} `yaml:"unison,omitempty"`
		Proxy struct {
			Hostname    string   `yaml:"hostname,omitempty"`
			Entrypoint  []string `yaml:"entrypoint,omitempty"`
			User        string   `yaml:"user,omitempty"`
			Cmd         string   `yaml:"cmd,omitempty"`
			Dependson   []string `yaml:"depends_on,omitempty"`
			Environment []string `yaml:"environment,omitempty"`
			Volumes     []string `yaml:"volumes,omitempty"`
			Volumesfrom []string `yaml:"volumes_from,omitempty"`
			Image       string   `yaml:"image"`
			Ports       []string `yaml:"ports"`
			Networks    []string `yaml:"networks"`
		} `yaml:"proxy"`
	} `yaml:"services"`
}

// ComposeDefaults - Template byte array for proxy `docker-compose.yml`
func ComposeDefaults() []byte {
	return []byte(`version: "2"
services:
  proxy:
    image: tokaido/proxy:` + constants.ProxyStableVersion + `
    ports:
      - "` + constants.ProxyPort + `:` + constants.ProxyPort + `"
    volumes:
      - ./client:/tokaido/proxy/config/client
    networks:
      - default
      - tokaido_proxy

`)
}

// ComposeDefaultsUnison - Template for proxy's docker-compose file used if Unison sync is enabled
func ComposeDefaultsUnison() []byte {
	return []byte(`version: "2"
services:
  unison:
    image: tokaido/unison:2.51.2
    environment:
      - UNISON_DIR=/tokaido/proxy/config/client
      - UNISON_UID=1002
      - UNISON_GID=1001
    ports:
      - "5000"
    volumes:
      - /tokaido/proxy/config/client
    networks:
      - default
      - tokaido_proxy
  proxy:
    image: tokaido/proxy:` + constants.ProxyStableVersion + `
    ports:
      - "` + constants.ProxyPort + `:` + constants.ProxyPort + `"
    volumes_from:
	  - unison
    networks:
      - default
      - tokaido_proxy
`)
}
