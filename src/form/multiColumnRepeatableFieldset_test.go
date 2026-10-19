package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func multiColumnRepeatableFieldsetTestColumns() []MultiColumnRepeatableFieldsetColumn {
	return []MultiColumnRepeatableFieldsetColumn{
		{Id: "name", Label: "Name", InputType: MultiColumnRepeatableFieldsetColumnTypeText, IsRequired: true},
		{
			Id: "role", Label: "Role", InputType: MultiColumnRepeatableFieldsetColumnTypeSelect,
			Options: []SelectLabelValueOption{
				{Label: "Owner", Value: "owner"},
				{Label: "Member", Value: "member"},
			},
		},
		{Id: "age", Label: "Age", InputType: MultiColumnRepeatableFieldsetColumnTypeNumber, Min: "0", Max: "120", Step: "1"},
		{Id: "active", Label: "Active", InputType: MultiColumnRepeatableFieldsetColumnTypeCheckbox},
	}
}

func TestMultiColumnRepeatableFieldsetColumnMetadataJsonBuilder(t *testing.T) {
	metadataJson := multiColumnRepeatableFieldsetColumnMetadataJsonBuilder(multiColumnRepeatableFieldsetTestColumns())
	expectedFragments := []string{
		`{"id":"name","inputType":"text","isRequired":true}`,
		`{"id":"role","inputType":"select","isRequired":false}`,
		`{"id":"age","inputType":"number","isRequired":false}`,
		`{"id":"active","inputType":"checkbox","isRequired":false}`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(metadataJson, expectedFragment) {
			t.Errorf("metadataJson missing fragment %s: %s", expectedFragment, metadataJson)
		}
	}
}

func TestMultiColumnRepeatableFieldsetLabelResolvers(t *testing.T) {
	defaultSettings := MultiColumnRepeatableFieldsetSettings{}
	if resolved := multiColumnRepeatableFieldsetAddRowLabelResolver(defaultSettings); resolved != "Add row" {
		t.Errorf("addRowLabel default = %s, expected Add row", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetEmptyStateLabelResolver(defaultSettings); resolved != "No rows yet" {
		t.Errorf("emptyStateLabel default = %s, expected No rows yet", resolved)
	}
	customSettings := MultiColumnRepeatableFieldsetSettings{
		AddRowLabel:     "Append contact",
		EmptyStateLabel: "No contacts",
	}
	if resolved := multiColumnRepeatableFieldsetAddRowLabelResolver(customSettings); resolved != "Append contact" {
		t.Errorf("addRowLabel custom = %s, expected Append contact", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetEmptyStateLabelResolver(customSettings); resolved != "No contacts" {
		t.Errorf("emptyStateLabel custom = %s, expected No contacts", resolved)
	}
}

func TestMultiColumnRepeatableFieldsetInputTypeResolver(t *testing.T) {
	tests := []struct {
		inputType string
		expected  string
	}{
		{MultiColumnRepeatableFieldsetColumnTypeText, MultiColumnRepeatableFieldsetColumnTypeText},
		{MultiColumnRepeatableFieldsetColumnTypeNumber, MultiColumnRepeatableFieldsetColumnTypeNumber},
		{MultiColumnRepeatableFieldsetColumnTypeSelect, MultiColumnRepeatableFieldsetColumnTypeText},
		{MultiColumnRepeatableFieldsetColumnTypeCheckbox, MultiColumnRepeatableFieldsetColumnTypeText},
		{"", MultiColumnRepeatableFieldsetColumnTypeText},
	}
	for _, test := range tests {
		column := MultiColumnRepeatableFieldsetColumn{InputType: test.inputType}
		if resolved := multiColumnRepeatableFieldsetInputTypeResolver(column); resolved != test.expected {
			t.Errorf("inputTypeResolver(%s) = %s, expected %s", test.inputType, resolved, test.expected)
		}
	}
}

func TestMultiColumnRepeatableFieldsetRender(t *testing.T) {
	component := MultiColumnRepeatableFieldset(MultiColumnRepeatableFieldsetSettings{
		Columns:         multiColumnRepeatableFieldsetTestColumns(),
		TwoWayStatePath: "contactRows",
		InputName:       "contacts",
	})
	var buffer bytes.Buffer
	if err := component.Render(context.Background(), &buffer); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`x-for="(row, rowIndex) in contactRows" :key="rowIndex"`,
		`:class="{&#39;mt-0&#39;: !row[&#34;name&#34;], &#39;!border-red-500&#39;: columnIsInvalid(row, &#34;name&#34;)}"`,
		`x-model:value="row[&#34;name&#34;]"`,
		`:aria-invalid="columnIsInvalid(row, &#34;name&#34;)"`,
		`placeholder="Name"`,
		`x-model:value="row[&#34;age&#34;]"`,
		`type="number"`,
		`min="0"`,
		`max="120"`,
		`step="1"`,
		`x-model="row[&#34;active&#34;]"`,
		`type="checkbox"`,
		`String(row[&#34;role&#34;]) === $el.value`,
		`{"label":"Owner","value":"owner"}`,
		`{"label":"Member","value":"member"}`,
		`<span>--</span>`,
		`appendRow(contactRows)`,
		`duplicateRow(contactRows, rowIndex)`,
		`moveRow(contactRows, rowIndex, -1)`,
		`moveRow(contactRows, rowIndex, 1)`,
		`removeRow(contactRows, rowIndex)`,
		`x-show="contactRows.length === 0"`,
		`colspan="5"`,
		`name="contacts"`,
		`:value="JSON.stringify(contactRows)"`,
		`Add row`,
		`No rows yet`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
}

func TestMultiColumnRepeatableFieldsetRenderNewRowFunc(t *testing.T) {
	component := MultiColumnRepeatableFieldset(MultiColumnRepeatableFieldsetSettings{
		Columns:         multiColumnRepeatableFieldsetTestColumns(),
		TwoWayStatePath: "contactRows",
		NewRowFunc:      "demoRandomRow()",
	})
	var buffer bytes.Buffer
	if err := component.Render(context.Background(), &buffer); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	renderedHtml := buffer.String()
	xDataStartIndex := strings.Index(renderedHtml, `x-data="`) + len(`x-data="`)
	xDataEndIndex := strings.Index(renderedHtml[xDataStartIndex:], `"`)
	renderedXData := renderedHtml[xDataStartIndex : xDataStartIndex+xDataEndIndex]
	if !strings.Contains(renderedXData, `demoRandomRow()`) {
		t.Errorf("rendered html missing the NewRowFunc expression: %s", renderedXData)
	}
}

func TestMultiColumnRepeatableFieldsetRenderStyleGroups(t *testing.T) {
	component := MultiColumnRepeatableFieldset(MultiColumnRepeatableFieldsetSettings{
		Columns:         multiColumnRepeatableFieldsetTestColumns(),
		TwoWayStatePath: "contactRows",
		Checkbox: MultiColumnRepeatableFieldsetCheckboxStyle{
			CheckedColor:   "emerald-500",
			FocusRingColor: "emerald-400/50",
			Shape:          uiToolset.ShapeCircular,
			UncheckedColor: "emerald-900",
		},
		Dropdown: MultiColumnRepeatableFieldsetDropdownStyle{
			BackgroundColor: "neutral-900",
			BorderColor:     "neutral-50/20",
		},
		Select: MultiColumnRepeatableFieldsetSelectStyle{
			SelectedIconColor: "emerald-400",
			Trigger: SelectInputTriggerStyle{
				BackgroundColor: "neutral-50/5",
				BorderColor:     "neutral-50/10",
			},
		},
	})
	var buffer bytes.Buffer
	if err := component.Render(context.Background(), &buffer); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`peer-checked:border-emerald-500`,
		`peer-checked:text-emerald-500`,
		`peer-focus-visible:ring-emerald-400/50`,
		`rounded-full`,
		`border-emerald-900`,
		`bg-neutral-900`,
		`border-neutral-50/20`,
		`text-emerald-400`,
		`bg-neutral-50/5`,
		`border-neutral-50/10`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
}

func TestMultiColumnRepeatableFieldsetRenderDefaults(t *testing.T) {
	component := MultiColumnRepeatableFieldset(MultiColumnRepeatableFieldsetSettings{
		Columns: []MultiColumnRepeatableFieldsetColumn{
			{Id: "note", Label: "Note", InputType: MultiColumnRepeatableFieldsetColumnTypeText},
		},
		TwoWayStatePath: "noteRows",
		AddRowLabel:     "Append note",
		EmptyStateLabel: "No notes",
	})
	var buffer bytes.Buffer
	if err := component.Render(context.Background(), &buffer); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`Append note`,
		`No notes`,
		`colspan="2"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, `type="hidden"`) {
		t.Error("rendered html contains hidden input without InputName")
	}
}
