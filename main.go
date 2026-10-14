package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"helm.sh/helm/v4/pkg/action"
	"slices"
)

// Local chartmuseum API
const localChartServer = "http://localhost:8879/api/charts"
const charts_dir = "../local_packages/"

func CreatePackages () {
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
		// Format: dependencyName/version
		dependencyName := fmt.Sprintf("%s", dep.Name)
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
}

func PostPackages() {
	packages, err := os.ReadDir(charts_dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, chartPackage := range packages {
		fileName := fmt.Sprintf("%s/%s", charts_dir, chartPackage.Name())
		chart, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(
			localChartServer,
			"application/octet-stream",
			bytes.NewReader(chart),
		)
		if err != nil {
        	log.Fatal(err)
    	}
		fmt.Println(resp.Status)
    	defer resp.Body.Close()
	}
}

func main() {
	// First read the directory for local charts
	packages, err := os.ReadDir(charts_dir)
	if err != nil {
		log.Fatal(err)
	}

	// empties the old dir, so no issues when packages were previously created.
	if len(packages) > 0 {
		err := os.RemoveAll(charts_dir)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Mkdir(charts_dir, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	CreatePackages()
	PostPackages()
}

// bad; drw-r--r--
// good: drwxr-xr-x = 755
