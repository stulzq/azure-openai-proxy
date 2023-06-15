package azure

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"text/template"
)

type DeploymentConfig struct {
	DeploymentName string   `yaml:"deployment_name" json:"deployment_name" mapstructure:"deployment_name"` // azure openai deployment name
	ModelName      string   `yaml:"model_name" json:"model_name" mapstructure:"model_name"`                // corresponding model name in openai
	Endpoint       string   `yaml:"endpoint" json:"endpoint" mapstructure:"endpoint"`                      // deployment endpoint
	ApiKey         string   `yaml:"api_key" json:"api_key" mapstructure:"api_key"`                         // secrect key1 or 2
	ApiVersion     string   `yaml:"api_version" json:"api_version" mapstructure:"api_version"`             // deployment version, not required
	EndpointUrl    *url.URL // url.URL form deployment endpoint
}

type Config struct {
	ApiBase          string             `yaml:"api_base" mapstructure:"api_base"`                   // if you use openai、langchain as sdk, it will be useful
	DeploymentConfig []DeploymentConfig `yaml:"deployment_config" mapstructure:"deployment_config"` // deployment config
}

type RequestConverter interface {
	Name() string
	Convert(req *http.Request, config *DeploymentConfig) (*http.Request, error)
}

type StripPrefixConverter struct {
	Prefix string
}

func (c *StripPrefixConverter) Name() string {
	return "StripPrefix"
}
func (c *StripPrefixConverter) Convert(req *http.Request, config *DeploymentConfig) (*http.Request, error) {
	req.Host = config.EndpointUrl.Host
	req.URL.Scheme = config.EndpointUrl.Scheme
	req.URL.Host = config.EndpointUrl.Host
	req.URL.Path = path.Join(fmt.Sprintf("/openai/deployments/%s", config.DeploymentName), strings.Replace(req.URL.Path, c.Prefix+"/", "/", 1))
	req.URL.RawPath = req.URL.EscapedPath()

	query := req.URL.Query()
	query.Add("api-version", config.ApiVersion)
	req.URL.RawQuery = query.Encode()
	return req, nil
}
func NewStripPrefixConverter(prefix string) *StripPrefixConverter {
	return &StripPrefixConverter{
		Prefix: prefix,
	}
}

type TemplateConverter struct {
	Tpl      string
	Tempalte *template.Template
}

func (c *TemplateConverter) Name() string {
	return "Template"
}
func (c *TemplateConverter) Convert(req *http.Request, config *DeploymentConfig) (*http.Request, error) {
	data := map[string]interface{}{
		"DeploymentName": config.DeploymentName,
		"ModelName":      config.ModelName,
		"Endpoint":       config.Endpoint,
		"ApiKey":         config.ApiKey,
		"ApiVersion":     config.ApiVersion,
	}
	buff := new(bytes.Buffer)
	if err := c.Tempalte.Execute(buff, data); err != nil {
		return req, errors.Wrap(err, "template execute error")
	}

	req.Host = config.EndpointUrl.Host
	req.URL.Scheme = config.EndpointUrl.Scheme
	req.URL.Host = config.EndpointUrl.Host
	req.URL.Path = buff.String()
	req.URL.RawPath = req.URL.EscapedPath()

	query := req.URL.Query()
	query.Add("api-version", config.ApiVersion)
	req.URL.RawQuery = query.Encode()
	return req, nil
}
func NewTemplateConverter(tpl string) *TemplateConverter {
	_template, err := template.New("template").Parse(tpl)
	if err != nil {
		log.Fatalf("template parse error: %s", err.Error())
	}
	return &TemplateConverter{
		Tpl:      tpl,
		Tempalte: _template,
	}
}
