package uiToolset

import (
	"github.com/a-h/templ"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
)

const (
	MinifierContentTypeJavaScript string = "application/javascript"
	MinifierContentTypeCss        string = "text/css"
)

type MinifierSettings struct {
	ContentType       string
	UnminifiedContent *string
	MinifyWhitespace  bool
	MinifyIdentifiers bool
	MinifySyntax      bool
}

func Minifier(componentSettings MinifierSettings) *string {
	var minifierResult esbuildApi.TransformResult

	transformOptions := esbuildApi.TransformOptions{
		MinifyWhitespace:  componentSettings.MinifyWhitespace,
		MinifyIdentifiers: componentSettings.MinifyIdentifiers,
		MinifySyntax:      componentSettings.MinifySyntax,
	}

	switch componentSettings.ContentType {
	case MinifierContentTypeJavaScript:
		transformOptions.Loader = esbuildApi.LoaderJS
		minifierResult = esbuildApi.Transform(*componentSettings.UnminifiedContent, transformOptions)
	case MinifierContentTypeCss:
		transformOptions.Loader = esbuildApi.LoaderCSS
		minifierResult = esbuildApi.Transform(*componentSettings.UnminifiedContent, transformOptions)
	default:
		return componentSettings.UnminifiedContent
	}

	if len(minifierResult.Errors) > 0 {
		return componentSettings.UnminifiedContent
	}

	minifiedContent := string(minifierResult.Code)
	return &minifiedContent
}

func MinifierJs(unminifiedContent *string) *string {
	return Minifier(MinifierSettings{
		ContentType:       MinifierContentTypeJavaScript,
		UnminifiedContent: unminifiedContent,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})
}

func MinifierCss(unminifiedContent *string) *string {
	return Minifier(MinifierSettings{
		ContentType:       MinifierContentTypeCss,
		UnminifiedContent: unminifiedContent,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})
}

func MinifierTemplate(componentSettings MinifierSettings) templ.Component {
	return templ.Raw(*Minifier(componentSettings))
}

func MinifierTemplateJs(unminifiedContent *string) templ.Component {
	return templ.Raw(`<script type="text/javascript" charset="UTF-8">` + *MinifierJs(unminifiedContent) + `</script>`)
}

func MinifierTemplateCss(unminifiedContent *string) templ.Component {
	return templ.Raw(`<style>` + *MinifierCss(unminifiedContent) + `</style>`)
}
