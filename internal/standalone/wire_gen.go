// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package standalone

import (
	"context"
	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/cache"
	"github.com/aquasecurity/fanal/extractor/docker"
	"github.com/aquasecurity/trivy-db/pkg/db"
	"github.com/aquasecurity/trivy/internal/operation"
	"github.com/aquasecurity/trivy/pkg/detector/library"
	"github.com/aquasecurity/trivy/pkg/detector/ospkg"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/scanner/local"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/aquasecurity/trivy/pkg/vulnerability"
)

// Injectors from inject.go:

func initializeCacheClient(cacheDir string) (operation.Cache, error) {
	fsCache, err := cache.NewFSCache(cacheDir)
	if err != nil {
		return operation.Cache{}, err
	}
	operationCache := operation.NewCache(fsCache)
	return operationCache, nil
}

func initializeScanner(ctx context.Context, imageName string, layerCache cache.LayerCache, localLayerCache cache.LocalLayerCache) (scanner.Scanner, error) {
	applier := analyzer.NewApplier(localLayerCache)
	detector := ospkg.Detector{}
	driverFactory := library.DriverFactory{}
	libraryDetector := library.NewDetector(driverFactory)
	localScanner := local.NewScanner(applier, detector, libraryDetector)
	dockerOption, err := types.GetDockerOption()
	if err != nil {
		return scanner.Scanner{}, err
	}
	extractor, err := docker.NewDockerExtractor(ctx, imageName, dockerOption)
	if err != nil {
		return scanner.Scanner{}, err
	}
	config := analyzer.New(extractor, layerCache)
	scannerScanner := scanner.NewScanner(localScanner, config, layerCache)
	return scannerScanner, nil
}

func initializeVulnerabilityClient() vulnerability.Client {
	config := db.Config{}
	client := vulnerability.NewClient(config)
	return client
}
