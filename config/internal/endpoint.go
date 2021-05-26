package internal

import "fmt"

type ReverseProxy struct {
	Url     string
	Replace *string
}

type Static struct {
	Dir     string
	Page404 *string
}

func (s *Static) check(dir string) error {
	var err error
	s.Dir, err = absPath(s.Dir, dir)
	if err != nil {
		return err
	}
	page := "404.html"
	if s.Page404 != nil {
		page = *s.Page404
	}
	s.Page404 = &page
	return err
}

type Deploy struct {
	Token   string
	Command []string
	Dir     string
}

func (d *Deploy) check(dir string) error {
	if d.Command == nil || len(d.Command) == 0 {
		return fmt.Errorf("No command for deploy")
	}
	if d.Dir == "" {
		d.Dir = dir
	}
	return nil
}

type Endpoint struct {
	Host         interface{}
	Path         *string
	ReverseProxy *ReverseProxy `yaml:"reverseProxy"`
	Static       *Static
	Redirect     *string
	Deploy       *Deploy
}

func (e *Endpoint) check(dir string, index int) error {
	count := 0
	if e.ReverseProxy != nil {
		count++
	}
	if e.Static != nil {
		if err := e.Static.check(dir); err != nil {
			return err
		}
		count++
	}
	if e.Redirect != nil {
		count++
	}
	if e.Deploy != nil {
		if err := e.Deploy.check(dir); err != nil {
			return err
		}
		count++
	}
	if count > 1 {
		return fmt.Errorf("Too many actions in the endpoint with index %s", index)
	}
	if count == 0 {
		return fmt.Errorf("No actions in the endpoint with index %s", index)
	}
	hosts, err := e.getHosts()
	if err != nil {
		return err
	}
	e.Host = hosts
	if e.Path == nil {
		s := "/"
		e.Path = &s
	}
	return nil
}

func (e Endpoint) getHosts() ([]string, error) {
	switch r := e.Host.(type) {
	case string:
		return []string{r}, nil
	case []interface{}:
		hosts := []string{}
		for _, h := range r {
			host, ok := h.(string)
			if !ok {
				return nil, fmt.Errorf("Host must be string or array of strings")
			}
			hosts = append(hosts, host)
		}
		return hosts, nil
	default:
		return nil, fmt.Errorf("Host must be string or array of strings")
	}
}

func (e Endpoint) splitByHosts() []Endpoint {
	hosts := e.Host.([]string)
	result := []Endpoint{}
	for _, host := range hosts {
		e.Host = host
		result = append(result, e)
	}
	return result
}
