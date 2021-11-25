package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/sputnik-systems/alertrules-checker/internal/github"
	prom "github.com/sputnik-systems/alertrules-checker/internal/prometheus"
	"github.com/sputnik-systems/alertrules-checker/internal/utils"
	vm "github.com/sputnik-systems/alertrules-checker/internal/victoriametrics"
)

var (
	ErrMarshal   = errors.New("failed to marshal object")
	ErrUnmarshal = errors.New("failed to unmarshal object")

	ruleType string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "checker",
		Short: "Checking alert rule templates",
	}

	kustomizeCmd := &cobra.Command{
		Use:   "kustomize",
		Short: "Checking alert rules given like kustomize templates",
		Run:   kustomize,
	}

	rootCmd.AddCommand(kustomizeCmd)

	rootCmd.PersistentFlags().StringVar(&ruleType, "rule-type", "VMRule", "Alert rules definition kubernetes object type (VMRule or PrometheusRule)")

	log.SetFlags(0)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func kustomize(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatalln("::error ::at least one path should be given")
	}

	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

	var events []*github.Event
	for _, path := range args {
		m, err := k.Run(filesys.MakeFsOnDisk(), path)
		if err != nil {
			event := github.NewEvent("warning",
				fmt.Sprintf("failed to generate templates from given path %s: %s", path, err))
			events = append(events, event)
		}

		for _, resource := range m.Resources() {
			name := resource.GetName()

			b, err := resource.AsYAML()
			if err != nil {
				event := github.NewEvent("warning",
					fmt.Sprintf("failed to get resource %s as yaml: %s", name, err))
				events = append(events, event)
			}

			var perrs interface{}
			switch ruleType {
			case "VMRule":
				perrs = vm.Validate(b)
			case "PrometheusRule":
				perrs = prom.Validate(b)
			}

			if err, ok := perrs.(error); ok {
				event := github.NewEvent("warning", err.Error())
				events = append(events, event)
			} else if errgr, ok := perrs.(utils.ErrorGroup); ok {
				if errgr.Count() > 0 {
					event := github.NewEvent("warning",
						fmt.Sprintf("failed to parse resource with %d errors", errgr.Count())).WithTitle(name)
					events = append(events, event)

					for index, err := range errgr.List() {
						event := github.NewEvent("warning",
							fmt.Sprintf("resource error #%d: %s", index+1, err)).WithTitle(name)
						events = append(events, event)
					}
				}
			}

		}
	}

	for _, event := range events {
		log.Println(event)
	}

	if len(events) > 0 {
		log.Fatalf("::error ::checks failed with %d warnings", len(events))
	}
}
