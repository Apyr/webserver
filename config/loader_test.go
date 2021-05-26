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
  - host: localhost
    path: /deploy
    deploy:
      token: UFUFUF
      command: ["bash", "-c", "echo Deployed >> deploy.txt"]
  - host:
      - example.com
      - www.example.com
    redirect: https://example2.com
`

func TestBasicConfig(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.yml")
	srvPath := filepath.Join(dir, "service.yml")

	expectedConfig := Config{
		HttpPort:        80,
		HttpsPort:       443,
		RedirectToHttps: false,
		CertsDir:        filepath.Join(dir, "certs"),
		Logging:         false,
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
						Host: "localhost",
						Path: "/deploy",
						Action: Deploy{
							Token:   "UFUFUF",
							Command: []string{"bash", "-c", "echo Deployed >> deploy.txt"},
							Dir:     dir,
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
		ConfigFiles: []string{cfgPath, srvPath},
	}

	err := ioutil.WriteFile(cfgPath, []byte(basicConfig), 0666)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(srvPath, []byte(basicService), 0666)
	if err != nil {
		t.Error(err)
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Error(err)
	}

	cfgStr := fmt.Sprintf("%#v", cfg)
	expectedStr := fmt.Sprintf("%#v", expectedConfig)

	if cfgStr != expectedStr {
		t.Errorf("Got config:\n%s\nExpected:\n%s", cfgStr, expectedStr)
	}
}
