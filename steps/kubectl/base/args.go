package base

import (
	"encoding/json"
	"fmt"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
)

// You can specify ignore args if sub command don't want to add certain arguments to the kubectl execution
type IgnoredArgs int8

const (
	IgnoreFormat IgnoredArgs = iota
	IgnoreNamespace
	IgnoreFieldSelector
)

type Args struct {
	KubeConfigContent       string            `env:"KUBECONFIG_CONTENT"`
	LocalToken              bool              `env:"-"`
	GcloudAuth              string            `env:"GCLOUD_AUTH_CODE_B64"`
	Namespace               string            `env:"NAMESPACE" envDefault:"default"`
	Format                  string            `env:"FORMAT" envDefault:"json"`
	AllNamespaces           bool              `env:"ALL_NAMESPACES" envDefault:"false"`
	FieldSelector           []string          `env:"FIELD_SELECTOR"`
	Debug                   bool              `env:"DEBUG" envDefault:"false"`
	Pretty                  bool              `env:"PRETTY" envDefault:"false"`
	OrderBy                 string            `env:"ORDER_BY"`
	OrderByDescending       bool              `env:"ORDER_BY_DESC" envDefault:"false"`
	FilterContains          string            `env:"FILTER_CONTAINS"`
	FilterEquals            string            `env:"FILTER_EQUALS"`
	FilterNotEquals         string            `env:"FILTER_NOT_EQUALS"`
	FilterNotContains       string            `env:"FILTER_NOT_CONTAINS"`
	FilterContainsParsed    map[string]string `env:"-"`
	FilterEqualsParsed      map[string]string `env:"-"`
	FilterNotEqualsParsed   map[string]string `env:"-"`
	FilterNotContainsParsed map[string]string `env:"-"`
}

func (a *Args) BaseArgs() *Args {
	return a
}

type BaseArgs interface {
	BaseArgs() *Args
}

func parseFilter(filter string) (map[string]string, error) {
	filterMap := make(map[string]string)
	if filter == "" {
		return filterMap, nil
	}
	if err := json.Unmarshal([]byte(filter), &filterMap); err != nil {
		return nil, err
	}
	return filterMap, nil
}

func Parse(args BaseArgs) error {
	if err := env.Parse(args); err != nil {
		return err
	}
	baseArgs := args.BaseArgs()

	if baseArgs.KubeConfigContent == "" {
		log.Debugln("No integration specified, using kubernetes local token")
		baseArgs.LocalToken = true
	}

	filters, err := parseFilter(baseArgs.FilterContains)
	if err != nil {
		return fmt.Errorf("can't parse FILTER_CONTAINS as json map: %w", err)
	}
	baseArgs.FilterContainsParsed = filters

	filters, err = parseFilter(baseArgs.FilterEquals)
	if err != nil {
		return fmt.Errorf("can't parse FILTER_EQUALS as json map: %w", err)
	}
	baseArgs.FilterEqualsParsed = filters

	filters, err = parseFilter(baseArgs.FilterNotContains)
	if err != nil {
		return fmt.Errorf("can't parse FILTER_NOT_CONTAINS as json map: %w", err)
	}
	baseArgs.FilterNotContainsParsed = filters

	filters, err = parseFilter(baseArgs.FilterNotEquals)
	if err != nil {
		return fmt.Errorf("can't parse FILTER_NOT_EQUALS as json map: %w", err)
	}
	baseArgs.FilterNotEqualsParsed = filters

	return nil
}
