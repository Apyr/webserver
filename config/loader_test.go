package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

var basicConfig = `
certsDir: certs
redirectToHttps: false
httpPort: 80
httpsPort: 443
services:
  default: service.yml
`

var basicService = `
enabled: true
endpoints:
  - host: localhost
    path: /
    static: 
      dir: .
      page404: page404.html
  - host: localhost
    path: /api
    reverseProxy:
      host: example.com
      port: 80
  - host: example.com
    redirect: https://example.com
`

var expectedConfig = Config{
	HttpPort:        80,
	HttpsPort:       443,
	RedirectToHttps: false,
	CertsDir:        "certs",
	Services: []Service{
		Service{
			Name:    "default",
			Enabled: true,
			Endpoints: []Endpoint{
				Endpoint{
					Host: "localhost",
					Path: "/",
					Action: ServeStatic{
						Dir:     ".",
						Page404: "page404.html",
					},
				},
				Endpoint{
					Host: "localhost",
					Path: "/api",
					Action: ReverseProxy{
						Host: "example.com",
						Port: 80,
					},
				},
				Endpoint{
					Host: "example.com",
					Path: "/",
					Action: Redirect{
						To: "https://example.com",
					},
				},
			},
		},
	},
}

func TestBasicConfig(t *testing.T) {
	dir := t.TempDir()

	err := ioutil.WriteFile(filepath.Join(dir, "config.yml"), []byte(basicConfig), 0666)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "service.yml"), []byte(basicService), 0666)
	if err != nil {
		t.Error(err)
	}

	cfg, err := LoadConfig(filepath.Join(dir, "config.yml"))
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(cfg, expectedConfig) {
		t.Errorf("Got config: %#v", cfg)
	}
}
