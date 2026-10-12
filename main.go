package main

import (
	//"bytes"
	"fmt"
	"log"
	//"net/http"
	"os"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"helm.sh/helm/v4/pkg/action"
	"slices"
)

// Local chartmuseum API
const localChartServer = "http://localhost:8879/api/charts/"

func main() {
	// First read the directory for local charts
	charts_dir := "../local_packages/"
	packages, err := os.ReadDir(charts_dir)
	if err != nil {
		log.Fatal(err)
	}
	// If it is empty, run helm package
	if len(packages) == 0 {
		// List dependencies
		chartpath := "/path/to/chart/"
		// from https://github.com/helm/helm/blob/v4.3.0/pkg/action/dependency.go#L30
		c, err := loader.Load(chartpath)
		if err != nil {
			log.Fatal(err)
		}
		
		if c.Metadata.Dependencies == nil {
			log.Fatal("No dependencies found.")
		}

		var dependencyList = []string{}
		for _, dep := range c.Metadata.Dependencies {
			// Format: dependencyName
			dependencyName := fmt.Sprintf("../%s", dep.Name)
			if !slices.Contains(dependencyList, dependencyName) {
				dependencyList = append(dependencyList, dependencyName)
			}
		}

		// Initialize Package object and configure
		packageHandler := action.NewPackage()
		packageHandler.Destination = charts_dir

		for _, depPackage := range dependencyList {
			// At the time of creating this there was a "_" argument: 
			// https://github.com/helm/helm/blob/v4.3.0/pkg/action/package.go#L76
			cryptidArg := map[string]any{}
			out, err := packageHandler.Run(depPackage, cryptidArg)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Saved to: ", out)
			}
		}
	} else {
		fmt.Println("Empty directory local_packages before running.")
	}

	//for _, chart := range packages {
	//	fmt.Println(chart.Name())
	//}
	//data, err := os.ReadFile()
}
