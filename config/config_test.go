package config

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	var configYaml = `
logFile: /home/log.log
certsFile: /home/certs.json
ports:
  http: 8080
  https: 8443
endpoints:
  - url: localhost
    static: 
      dir: /home/web/static
`
	endpointURL, err := url.Parse("http://localhost")
	assert.NoError(t, err)
	assert.Equal(t, endpointURL.Hostname(), "localhost")

	u, err := url.Parse("http://localhost")
	assert.NoError(t, err)
	expected := Config{
		LogFile:   "/home/log.log",
		CertsFile: "/home/certs.json",
		Ports:     Ports{HTTP: 8080, HTTPS: 8443},
		LogLevel:  "debug",
		Endpoints: []Endpoint{{
			URL:             URL{u},
			HTTPS:           "",
			RedirectToHTTPS: false,
			Static: &Static{
				Dir:      "/home/web/static",
				Index:    "index.html",
				NotFound: "404.html",
			},
		}},
	}

	var config Config
	err = config.LoadFromYAML([]byte(configYaml))
	assert.NoError(t, err)

	assert.Equal(t, expected, config)
}

func TestConfigHosts(t *testing.T) {
	var configYaml = `
endpoints:
  - url: localhost1
    static:
      dir: /static
  - url: localhost2
    static:
      dir: /static
  - url: localhost3/abc
    static:
      dir: /static
  - url: localhost3/def
    static:
      dir: /static
`

	var config Config
	err := config.LoadFromYAML([]byte(configYaml))
	assert.NoError(t, err)

	hosts := config.AllHosts()
	assert.Equal(t, 3, len(hosts))
	assert.Contains(t, hosts, "localhost1")
	assert.Contains(t, hosts, "localhost2")
	assert.Contains(t, hosts, "localhost3")
}
