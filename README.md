Simple binary and github action for checking VictoriaMetrics operator alert rules (VMRules).

# GitHub action
## inputs
Available inputs:
* `path` - path with resource templates (default: `./`)
* `rule_type` - kubernetes rule object kind (`VMRule` or `PrometheusRule`)
* `subcommand` - subcommand with parser type (possible types: `kustomize`)
