# My HTTP Server

Simple http and https server (static and reverse proxy) written in Go. 

### Features
* minimal configuring in yaml file
* logging
* serving static files
* auto letsencrypt
* http to https redirect
* reverse proxy
* redirect
* protected api routes which can execute arbitrary commands (useful for deploy)

### Config

```yaml
logFile: /path/to/log/file.log
logLevel: debug #one of debug|info|warn|error|fatal
certsFile: /path/to/certs.json
ports:
  http: 80
  https: 443

endpoints:

  - url: proxy.example.com/myapp
    https: true #default: true
    redirectToHttps: true #default: true
    disabled: true #default: false
    proxy: 
        url: https://example.com
        removePrefix: /myapp

  - url: static.example.com
    static: 
      dir: /path/to/dir/with/static/files
      index: index.html
      notFound: 404.html

  - url: redirect.example.com
    redirect: https://google.com

  - url: command.example.com
    runCommand:
      command: ["/bin/bash", "/root/deploy.sh"]
      token: "my-very-secret-token"

```

### License
MIT
