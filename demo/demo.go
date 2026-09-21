package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/a-h/templ"
	uiStructural "github.com/goinfinite/ui/src/structural"
)

//go:embed dataTableDemoRouting.js
var dataTableDemoRoutingScript string

type demoArtifact struct {
	component templ.Component
	filePath  string
}

func (artifact demoArtifact) fileRender() error {
	file, err := os.Create(artifact.filePath)
	if err != nil {
		return errors.New(
			"CreateHtmlFileFailed (" + artifact.filePath + "): " + err.Error(),
		)
	}

	renderErr := artifact.component.Render(context.Background(), file)
	closeErr := file.Close()
	if renderErr != nil {
		return errors.New(
			"WriteHtmlFileFailed (" + artifact.filePath + "): " + renderErr.Error(),
		)
	}
	if closeErr != nil {
		return errors.New(
			"CloseHtmlFileFailed (" + artifact.filePath + "): " + closeErr.Error(),
		)
	}
	return nil
}

func newDemoArtifacts() []demoArtifact {
	artifacts := []demoArtifact{
		{component: DemoIndex(), filePath: "index.html"},
		{
			component: DataTableRefreshFragment(1, dataTableDemoItemsPerPage),
			filePath:  "assets/dataTableDemoRefresh.html",
		},
	}
	for pageNumber := uint(1); pageNumber <= dataTableDemoPagesTotal; pageNumber++ {
		artifacts = append(artifacts, demoArtifact{
			component: DataTableRefreshFragment(pageNumber, dataTableDemoItemsPerPage),
			filePath:  fmt.Sprintf("assets/dataTableDemoRefreshPage%d.html", pageNumber),
		})
	}
	return append(artifacts, demoArtifact{
		component: DataTableRefreshFragment(1, uiStructural.DataTablePageSize(len(dataTableDemoRecords))),
		filePath:  "assets/dataTableDemoRefreshAll.html",
	})
}

func main() {
	for _, artifact := range newDemoArtifacts() {
		err := artifact.fileRender()
		if err != nil {
			slog.Error("RenderDemoArtifactFailed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}
}
