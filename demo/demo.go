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

const dataTableDemoFragmentsAssetPath = "assets/dataTableDemoRefreshFragments.json"

func renderDemoIndex() error {
	indexFile, err := os.Create("index.html")
	if err != nil {
		return errors.New("CreateHtmlFileFailed (index.html): " + err.Error())
	}
	renderErr := DemoIndex().Render(context.Background(), indexFile)
	closeErr := indexFile.Close()
	if renderErr != nil {
		return errors.New(
			"WriteHtmlFileFailed (index.html): " + renderErr.Error(),
		)
	}
	if closeErr != nil {
		return errors.New(
			"CloseHtmlFileFailed (index.html): " + closeErr.Error(),
		)
	}
	return nil
}

func formatDataTableDemoFragmentKey(
	pageNumber uint, itemsPerPage uiStructural.DataTablePageSize,
) string {
	return fmt.Sprintf("%d-%d", pageNumber, itemsPerPage)
}

func renderDataTableDemoFragment(
	pageNumber uint, itemsPerPage uiStructural.DataTablePageSize,
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

func buildDataTableDemoFragments() (map[string]string, error) {
	fragments := map[string]string{}
	for pageNumber := uint(1); pageNumber <= dataTableDemoPagesTotal; pageNumber++ {
		fragmentBody, err := renderDataTableDemoFragment(
			pageNumber, dataTableDemoItemsPerPage,
		)
		if err != nil {
			return nil, err
		}
		fragmentKey := formatDataTableDemoFragmentKey(
			pageNumber, dataTableDemoItemsPerPage,
		)
		fragments[fragmentKey] = fragmentBody
	}
	allRecordsPageSize := uiStructural.DataTablePageSize(len(dataTableDemoRecords))
	allRecordsBody, err := renderDataTableDemoFragment(1, allRecordsPageSize)
	if err != nil {
		return nil, err
	}
	allRecordsKey := formatDataTableDemoFragmentKey(1, allRecordsPageSize)
	fragments[allRecordsKey] = allRecordsBody
	return fragments, nil
}

func writeDataTableDemoFragmentsAsset() error {
	fragments, err := buildDataTableDemoFragments()
	if err != nil {
		return err
	}
	fragmentsBody := &bytes.Buffer{}
	fragmentsEncoder := json.NewEncoder(fragmentsBody)
	fragmentsEncoder.SetEscapeHTML(false)
	encodeErr := fragmentsEncoder.Encode(fragments)
	if encodeErr != nil {
		return errors.New(
			"EncodeDataTableDemoFragmentsFailed: " + encodeErr.Error(),
		)
	}
	writeErr := os.WriteFile(
		dataTableDemoFragmentsAssetPath, fragmentsBody.Bytes(), 0o644,
	)
	if writeErr != nil {
		return errors.New(
			"WriteDataTableDemoFragmentsFailed (" +
				dataTableDemoFragmentsAssetPath + "): " + writeErr.Error(),
		)
	}
	return nil
}

func main() {
	err := renderDemoIndex()
	if err != nil {
		slog.Error("RenderDemoIndexFailed", slog.String("err", err.Error()))
		os.Exit(1)
	}
	err = writeDataTableDemoFragmentsAsset()
	if err != nil {
		slog.Error(
			"WriteDataTableDemoFragmentsFailed", slog.String("err", err.Error()),
		)
		os.Exit(1)
	}
}
