package victoriametrics

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmalert/config"
	vmutils "github.com/VictoriaMetrics/VictoriaMetrics/app/vmalert/utils"
	vmoperator "github.com/VictoriaMetrics/operator/api/v1beta1"
	"github.com/ghodss/yaml"

	"github.com/sputnik-systems/alertrules-checker/internal/utils"
)

func Validate(b []byte) interface{} {
	var rule vmoperator.VMRule

	if err := yaml.Unmarshal(b, &rule); err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	groups, err := getConfigGroups(rule.Spec.Groups)
	if err != nil {
		return fmt.Errorf("failed to get groups: %w", err)
	}

	var errs utils.ErrorGroup
	for _, group := range groups {
		if err := group.Validate(true, true); err != nil {
			errs.Add(err)
		}
	}

	return errs
}

func getConfigGroups(orgs []vmoperator.RuleGroup) (gs []config.Group, err error) {
	var d time.Duration

	for _, group := range orgs {
		var rs []config.Rule
		rs, err = getConfigRules(group.Rules)
		if err != nil {
			return nil, fmt.Errorf("failed to get rules: %s", err)
		}

		g := config.Group{
			Name:              group.Name,
			Rules:             rs,
			Concurrency:       group.Concurrency,
			Labels:            group.Labels,
			ExtraFilterLabels: group.ExtraFilterLabels,
		}

		if group.Interval != "" {
			d, err = time.ParseDuration(group.Interval)
			if err != nil {
				return nil, fmt.Errorf("failed to parse group \"Interval\" field value: %s", err)
			}

			g.Interval = vmutils.NewPromDuration(d)
		}

		gs = append(gs, g)
	}

	return
}

func getConfigRules(ors []vmoperator.Rule) (rs []config.Rule, err error) {
	var d time.Duration

	for _, rule := range ors {
		r := config.Rule{
			Record:      rule.Record,
			Alert:       rule.Alert,
			Expr:        rule.Expr.String(),
			Labels:      rule.Labels,
			Annotations: rule.Annotations,
		}

		if rule.For != "" {
			d, err = time.ParseDuration(rule.For)
			if err != nil {
				return nil, fmt.Errorf("failed to parse rule \"For\" field value: %s", err)
			}

			r.For = vmutils.NewPromDuration(d)
		}

		rs = append(rs, r)
	}

	return
}
