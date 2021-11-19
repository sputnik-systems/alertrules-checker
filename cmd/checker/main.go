package main

import (
	"log"

	vm "github.com/VictoriaMetrics/operator/api/v1beta1"
	"github.com/ghodss/yaml"
	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
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

	var warnCount uint
	for _, path := range args {
		m, err := k.Run(filesys.MakeFsOnDisk(), path)
		if err != nil {
			log.Printf("::warning title=%s ::failed to generate templates from given path: %s", path, err)

			warnCount++
		}

		for _, resource := range m.Resources() {
			name := resource.GetName()

			b, err := resource.AsYAML()
			if err != nil {
				log.Printf("::warning title=%s ::failed to get resource as yaml: %s", name, err)

				warnCount++
			}

			var rule vm.VMRule
			if err := yaml.Unmarshal(b, &rule); err != nil {
				log.Printf("::warning title=%s ::failed to unmarshal resource: %s", name, err)

				warnCount++
			}

			b, err = yaml.Marshal(rule.Spec)
			if err != nil {
				log.Printf("::warning title=%s ::failed to marshal resource into rule file: %s", name, err)

				warnCount++
			}

			_, errs := rulefmt.Parse(b)
			if len(errs) != 0 {
				log.Printf("::warning title=%s ::failed to parse resource with %d errors", name, len(errs))

				warnCount++

				for index, err := range errs {
					log.Printf("::warning title=%s ::resource error #%d: %s", name, index+1, err)
				}
			}
		}
	}

	if warnCount > 0 {
		log.Fatalf("::error ::checks failed with %d warnings", warnCount)
	}
}
