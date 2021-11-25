package prometheus

import (
	"fmt"

	"github.com/ghodss/yaml"
	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/model/rulefmt"

	"github.com/sputnik-systems/alertrules-checker/internal/utils"
)

func Validate(b []byte) interface{} {
	var rule promoperator.PrometheusRule

	if err := yaml.Unmarshal(b, &rule); err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	b, err := yaml.Marshal(rule.Spec)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	_, perrs := rulefmt.Parse(b)
	var errs utils.ErrorGroup
	for _, err := range perrs {
		errs.Add(err)
	}

	return errs
}
