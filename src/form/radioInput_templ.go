// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package uiForm

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

const (
	RadioInputSizeXs string = "xs"
	RadioInputSizeSm string = "sm"
	RadioInputSizeMd string = "md"
	RadioInputSizeLg string = "lg"
	RadioInputSizeXl string = "xl"
)

type RadioInputSettings struct {
	Label           string
	StateValue      string
	TwoWayStatePath string

	// OptionalFields
	InputId   string
	InputName string
	Size      string
}

func RadioInput(componentSettings RadioInputSettings) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- RadioInput HTML --><div class=\"hover:text-secondary-500 inline-flex items-center text-neutral-100\" @click=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(componentSettings.TwoWayStatePath + ` = "` + componentSettings.StateValue + `"`)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 26, Col: 90}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		labelTextSize := "text-sm"
		switch componentSettings.Size {
		case RadioInputSizeXs:
			labelTextSize = "text-xs"
		case RadioInputSizeSm:
			labelTextSize = "text-sm"
		case RadioInputSizeMd:
			labelTextSize = "text-base"
		case RadioInputSizeLg:
			labelTextSize = "text-lg"
		case RadioInputSizeXl:
			labelTextSize = "text-xl"
		}
		var templ_7745c5c3_Var3 = []any{"flex flex-row gap-1 cursor-pointer items-center " + labelTextSize}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var3...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<label class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var3).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\" :class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(`String(` + componentSettings.TwoWayStatePath + `) == "` + componentSettings.StateValue + `" && "text-secondary-500 font-bold"`)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 43, Col: 139}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if componentSettings.InputId != "" && componentSettings.InputName != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "<input type=\"radio\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if componentSettings.InputId != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, " id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(componentSettings.InputId)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 49, Col: 36}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if componentSettings.InputName != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, " name=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(componentSettings.InputName)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 52, Col: 40}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if componentSettings.StateValue != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, " value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(componentSettings.StateValue)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 55, Col: 42}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if componentSettings.TwoWayStatePath != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, " x-model=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(componentSettings.TwoWayStatePath)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 58, Col: 49}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, " class=\"hidden\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		iconSize := "text-xl"
		switch componentSettings.Size {
		case RadioInputSizeXs:
			iconSize = "text-lg"
		case RadioInputSizeSm:
			iconSize = "text-xl"
		case RadioInputSizeMd:
			iconSize = "text-2xl"
		case RadioInputSizeLg:
			iconSize = "text-3xl"
		case RadioInputSizeXl:
			iconSize = "text-4xl"
		}
		var templ_7745c5c3_Var10 = []any{"ph-duotone " + iconSize}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var10...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "<i class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var11 string
		templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var10).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "\" :class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var12 string
		templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(`String(` + componentSettings.TwoWayStatePath + `) == "` + componentSettings.StateValue + `" ? "ph-check-circle" : "ph-circle"`)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 78, Col: 140}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "\"></i> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var13 string
		templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(componentSettings.Label)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/form/radioInput.templ`, Line: 80, Col: 28}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "</label></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
