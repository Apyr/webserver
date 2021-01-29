package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
      url: http://example.com
  - host:
      - example.com
      - www.example.com
    redirect: https://example2.com
`

func TestBasicConfig(t *testing.T) {
	dir := t.TempDir()

	expectedConfig := Config{
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
						Action: Static{
							Dir:     filepath.Join(dir, "."),
							Page404: "page404.html",
						},
					},
					Endpoint{
						Host: "localhost",
						Path: "/api",
						Action: ReverseProxy{
							Url: "http://example.com",
						},
					},
					Endpoint{
						Host: "example.com",
						Path: "/",
						Action: Redirect{
							To: "https://example2.com",
						},
					},
					Endpoint{
						Host: "www.example.com",
						Path: "/",
						Action: Redirect{
							To: "https://example2.com",
						},
					},
				},
			},
		},
	}

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

	cfgStr := fmt.Sprintf("%#v", cfg)
	expectedStr := fmt.Sprintf("%#v", expectedConfig)

	if cfgStr != expectedStr {
		t.Errorf("Got config:\n%s\nExpected:\n%s", cfgStr, expectedStr)
	}
}
