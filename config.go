package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	// The page title
	Title string `yaml:"title"`

	// Text to show above the name
	Intro string `yaml:"intro"`
}

func RegisterConfigFlags() {
	pflag.StringP(`config`, `c`, `./thing-namer.yaml`, `configuration file`)
	pflag.StringP(`title`, `t`, `Thing Namer`, `title of the served page`)
	pflag.StringP(`intro`, `i`, `Your Thing is Now Named`, `introductory text to show on served page`)
}

func (c *Config) Parse(data io.Reader) error {
	return yaml.NewDecoder(data).Decode(c)
}

func configSetValueIfChanged(field *string, flag *pflag.Flag) {
	if flag == nil {
		return
	}
	if *field != `` && !flag.Changed {
		return
	}
	*field = flag.Value.String()
}

func (c *Config) CommandOverrides() {
	configSetValueIfChanged(&c.Title, pflag.Lookup(`title`))
	configSetValueIfChanged(&c.Intro, pflag.Lookup(`intro`))
}

func GetConfig() (*Config, error) {
	c := new(Config)
	RegisterConfigFlags()
	pflag.Parse()
	cflag := pflag.Lookup(`config`)
	if cflag != nil {
		cfilename := cflag.Value.String()
		if cfilename != `` {
			ofile, err := os.Open(cfilename)
			if err != nil {
				return nil, fmt.Errorf(`could not read config file %q: %w`, cfilename, err)
			}
			defer ofile.Close()
			if err := c.Parse(ofile); err != nil {
				return nil, fmt.Errorf(`could not understand config file %q: %w`, cfilename, err)
			}
		}
	}
	c.CommandOverrides()
	return c, nil
}
