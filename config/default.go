package config

import (
	"io/ioutil"
	"os"
)

var defaultConfig = `certsDir: ./certs
redirectToHttps: true
httpPort: 8080
httpsPort: 0
services:
  default: ./default/service.yml
`

var defaultService = `enabled: true
endpoints:
  - host: localhost:8080
    path: /
    static: 
      dir: ./static
      page404: 404.html
`

var defaultIndex = `<html>
<head><title>Default page</title></head>
<body>
    <h3>Default page</h3>
</body>
</html>
`

var default404 = `<html>
<head><title>404 Not found</title></head>
<body>
    <h3>404 Not found</h3>
</body>
</html>
`

func SaveDefault() {
	os.Mkdir("default", os.ModePerm|os.ModeDir)
	os.Mkdir("default/static", os.ModePerm|os.ModeDir)
	ioutil.WriteFile("config.yml", []byte(defaultConfig), os.ModePerm)
	ioutil.WriteFile("default/service.yml", []byte(defaultService), os.ModePerm)
	ioutil.WriteFile("default/static/index.html", []byte(defaultIndex), os.ModePerm)
	ioutil.WriteFile("default/static/404.html", []byte(default404), os.ModePerm)
}
