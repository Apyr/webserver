package config

import (
	"webserver/config/internal"
)

func toEndpoint(ep internal.Endpoint) Endpoint {
	endpoint := Endpoint{
		Host: ep.Host.(string),
		Path: *ep.Path,
	}
	if ep.ReverseProxy != nil {
		endpoint.Action = ReverseProxy{ep.ReverseProxy.Url, ep.ReverseProxy.Replace}
	} else if ep.Static != nil {
		endpoint.Action = Static{ep.Static.Dir, *ep.Static.Page404}
	} else if ep.Redirect != nil {
		endpoint.Action = Redirect{*ep.Redirect}
	} else if ep.Deploy != nil {
		endpoint.Action = Deploy{
			Token:   ep.Deploy.Token,
			Command: ep.Deploy.Command,
			Dir:     ep.Deploy.Dir,
		}
	}
	return endpoint
}

func LoadConfig(fileName string) (Config, error) {
	cfg, err := internal.Load(fileName)
	if err != nil {
		return Config{}, err
	}
	files := []string{fileName}
	services := []Service{}
	for _, srv := range cfg.ExtraServices {
		service := Service{
			Name:      srv.Name,
			Enabled:   *srv.Enabled,
			Endpoints: []Endpoint{},
		}
		for _, ep := range srv.Endpoints {
			service.Endpoints = append(service.Endpoints, toEndpoint(ep))
		}
		services = append(services, service)
		files = append(files, srv.FileName)
	}
	return Config{
		HttpPort:        cfg.HttpPort,
		HttpsPort:       cfg.HttpsPort,
		RedirectToHttps: *cfg.RedirectToHttps,
		CertsDir:        cfg.CertsDir,
		Logging:         cfg.Logging,
		Services:        services,
		ConfigFiles:     files,
	}, nil
}
