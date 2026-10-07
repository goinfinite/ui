package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	uiStructural "github.com/goinfinite/ui/src/structural"
)

//go:embed dataTableDemoRouting.js
var dataTableDemoRoutingScript string

const (
	dataTableDemoFragmentsAssetPath       = "docs/assets/dataTableDemoRefreshFragments.json"
	carouselDemoFragmentsAssetPath        = "docs/assets/carouselDemoRefreshFragments.json"
	carouselTooltipDemoFragmentsAssetPath = "docs/assets/carouselTooltipDemoRefreshFragments.json"
	carouselDemoItemsPerPage              = 6
)

type demoGenerator struct{}

func (generator demoGenerator) reconcileOutput(
	outputPath string, outputBody []byte,
) error {
	existingBody, readErr := os.ReadFile(outputPath)
	if readErr == nil && bytes.Equal(existingBody, outputBody) {
		return nil
	}
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		return errors.New(
			"ReadDemoOutputFailed (" + outputPath + "): " + readErr.Error(),
		)
	}
	writeErr := os.WriteFile(outputPath, outputBody, 0o644)
	if writeErr != nil {
		return errors.New(
			"WriteDemoOutputFailed (" + outputPath + "): " + writeErr.Error(),
		)
	}
	return nil
}

func (generator demoGenerator) renderIndex() error {
	indexBody := &bytes.Buffer{}
	renderErr := DemoIndex().Render(context.Background(), indexBody)
	if renderErr != nil {
		return errors.New(
			"RenderHtmlFailed (index.html): " + renderErr.Error(),
		)
	}
	return generator.reconcileOutput("docs/index.html", indexBody.Bytes())
}

func (generator demoGenerator) formatFragmentKey(
	pageNumber uint, itemsPerPage uiStructural.ItemsPerPage,
) string {
	return fmt.Sprintf("%d-%d", pageNumber, itemsPerPage)
}

func (generator demoGenerator) renderDataTableFragment(
	pageNumber uint, itemsPerPage uiStructural.ItemsPerPage,
) (string, error) {
	fragmentBody := &bytes.Buffer{}
	renderErr := DataTableRefreshFragment(pageNumber, itemsPerPage).
		Render(context.Background(), fragmentBody)
	if renderErr != nil {
		return "", errors.New(
			"RenderDataTableFragmentFailed: " + renderErr.Error(),
		)
	}
	return fragmentBody.String(), nil
}

func (generator demoGenerator) renderCarouselFragment(
	pageNumber uint, itemsPerPage uiStructural.ItemsPerPage,
) (string, error) {
	fragmentBody := &bytes.Buffer{}
	renderErr := CarouselRefreshFragment(pageNumber, itemsPerPage).
		Render(context.Background(), fragmentBody)
	if renderErr != nil {
		return "", errors.New(
			"RenderCarouselFragmentFailed: " + renderErr.Error(),
		)
	}
	return fragmentBody.String(), nil
}

func (generator demoGenerator) buildFragments(
	pagesTotal uint, itemsPerPage, allRecordsItemsPerPage uiStructural.ItemsPerPage,
	renderFragment func(
		pageNumber uint, itemsPerPage uiStructural.ItemsPerPage,
	) (string, error),
) (map[string]string, error) {
	fragments := map[string]string{}
	for pageNumber := range pagesTotal {
		fragmentBody, err := renderFragment(pageNumber, itemsPerPage)
		if err != nil {
			return nil, err
		}
		fragments[generator.formatFragmentKey(pageNumber, itemsPerPage)] = fragmentBody
	}
	allRecordsBody, err := renderFragment(0, allRecordsItemsPerPage)
	if err != nil {
		return nil, err
	}
	fragments[generator.formatFragmentKey(0, allRecordsItemsPerPage)] = allRecordsBody
	return fragments, nil
}

func (generator demoGenerator) buildDataTableFragments() (map[string]string, error) {
	return generator.buildFragments(
		dataTableDemoPagesTotal, dataTableDemoItemsPerPage,
		uiStructural.ItemsPerPage(len(dataTableDemoRecords)),
		generator.renderDataTableFragment,
	)
}

func (generator demoGenerator) buildCarouselFragments() (map[string]string, error) {
	carouselItemsPerPage := uiStructural.ItemsPerPage(carouselDemoItemsPerPage)
	carouselDemoPagesTotal := demoTablePageCountResolver(
		uint(len(dataTableDemoRecords)), carouselItemsPerPage,
	)
	return generator.buildFragments(
		carouselDemoPagesTotal, carouselItemsPerPage,
		uiStructural.ItemsPerPage(len(dataTableDemoRecords)),
		generator.renderCarouselFragment,
	)
}

func (generator demoGenerator) renderCarouselTooltipFragment(
	pageNumber uint, itemsPerPage uiStructural.ItemsPerPage,
) (string, error) {
	fragmentBody := &bytes.Buffer{}
	renderErr := CarouselTooltipRefreshFragment(pageNumber, itemsPerPage).
		Render(context.Background(), fragmentBody)
	if renderErr != nil {
		return "", errors.New(
			"RenderCarouselTooltipFragmentFailed: " + renderErr.Error(),
		)
	}
	return fragmentBody.String(), nil
}

func (generator demoGenerator) buildCarouselTooltipFragments() (map[string]string, error) {
	carouselItemsPerPage := uiStructural.ItemsPerPage(carouselDemoItemsPerPage)
	carouselDemoPagesTotal := demoTablePageCountResolver(
		uint(len(dataTableDemoRecords)), carouselItemsPerPage,
	)
	return generator.buildFragments(
		carouselDemoPagesTotal, carouselItemsPerPage,
		uiStructural.ItemsPerPage(len(dataTableDemoRecords)),
		generator.renderCarouselTooltipFragment,
	)
}

func (generator demoGenerator) writeFragmentsAsset(
	assetPath string, buildFragments func() (map[string]string, error),
) error {
	fragments, err := buildFragments()
	if err != nil {
		return err
	}
	fragmentsBody := &bytes.Buffer{}
	fragmentsEncoder := json.NewEncoder(fragmentsBody)
	fragmentsEncoder.SetEscapeHTML(false)
	encodeErr := fragmentsEncoder.Encode(fragments)
	if encodeErr != nil {
		return errors.New(
			"EncodeDemoFragmentsFailed (" + assetPath + "): " +
				encodeErr.Error(),
		)
	}
	return generator.reconcileOutput(assetPath, fragmentsBody.Bytes())
}

func main() {
	generator := demoGenerator{}
	renderErr := generator.renderIndex()
	if renderErr != nil {
		slog.Error("RenderDemoIndexFailed", slog.String("err", renderErr.Error()))
		os.Exit(1)
	}
	fragmentAssets := []struct {
		assetPath      string
		buildFragments func() (map[string]string, error)
		failureLogKey  string
	}{
		{
			assetPath:      dataTableDemoFragmentsAssetPath,
			buildFragments: generator.buildDataTableFragments,
			failureLogKey:  "WriteDataTableDemoFragmentsFailed",
		},
		{
			assetPath:      carouselDemoFragmentsAssetPath,
			buildFragments: generator.buildCarouselFragments,
			failureLogKey:  "WriteCarouselDemoFragmentsFailed",
		},
		{
			assetPath:      carouselTooltipDemoFragmentsAssetPath,
			buildFragments: generator.buildCarouselTooltipFragments,
			failureLogKey:  "WriteCarouselTooltipDemoFragmentsFailed",
		},
	}
	for _, fragmentAsset := range fragmentAssets {
		writeErr := generator.writeFragmentsAsset(
			fragmentAsset.assetPath, fragmentAsset.buildFragments,
		)
		if writeErr != nil {
			slog.Error(
				fragmentAsset.failureLogKey, slog.String("err", writeErr.Error()),
			)
			os.Exit(1)
		}
	}
}
