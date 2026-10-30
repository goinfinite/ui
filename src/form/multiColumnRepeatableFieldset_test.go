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
	if resolved := multiColumnRepeatableFieldsetAddRowLabelResolver(defaultSettings); resolved != "add row" {
		t.Errorf("addRowLabel default = %s, expected add row", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetAddRowIconResolver(defaultSettings); resolved != "ph-plus" {
		t.Errorf("addRowIcon default = %s, expected ph-plus", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetEmptyStateLabelResolver(defaultSettings); resolved != "No rows yet" {
		t.Errorf("emptyStateLabel default = %s, expected No rows yet", resolved)
	}
	customSettings := MultiColumnRepeatableFieldsetSettings{
		AddRowButton: MultiColumnRepeatableFieldsetAddRowButtonStyle{
			IconLeft: "ph-plus-circle",
		},
		AddRowLabel:     "Append contact",
		EmptyStateLabel: "No contacts",
	}
	if resolved := multiColumnRepeatableFieldsetAddRowLabelResolver(customSettings); resolved != "Append contact" {
		t.Errorf("addRowLabel custom = %s, expected Append contact", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetAddRowIconResolver(customSettings); resolved != "ph-plus-circle" {
		t.Errorf("addRowIcon custom = %s, expected ph-plus-circle", resolved)
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

func TestMultiColumnRepeatableFieldsetBulkEditResolvers(t *testing.T) {
	defaultSettings := MultiColumnRepeatableFieldsetSettings{
		Columns: []MultiColumnRepeatableFieldsetColumn{
			{Id: "name", InputType: MultiColumnRepeatableFieldsetColumnTypeText},
			{Id: "value", InputType: MultiColumnRepeatableFieldsetColumnTypeText},
		},
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditDelimiterResolver(defaultSettings); resolved != "=" {
		t.Errorf("bulkEditDelimiter default = %s, expected =", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditButtonLabelResolver(defaultSettings); resolved != "bulk edit" {
		t.Errorf("bulkEditButtonLabel default = %s, expected bulk edit", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditTableButtonLabelResolver(defaultSettings); resolved != "table view" {
		t.Errorf("bulkEditTableButtonLabel default = %s, expected table view", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditTextAreaLabelResolver(defaultSettings); resolved != "Bulk edit" {
		t.Errorf("bulkEditTextAreaLabel default = %s, expected Bulk edit", resolved)
	}
	expectedHint := `One row per line. Separate the columns with "=". Blank lines are ignored.`
	if resolved := multiColumnRepeatableFieldsetBulkEditHintValueResolver(defaultSettings); resolved != expectedHint {
		t.Errorf("bulkEditHintValue default = %s, expected %s", resolved, expectedHint)
	}
	singleColumnSettings := MultiColumnRepeatableFieldsetSettings{
		Columns: []MultiColumnRepeatableFieldsetColumn{
			{Id: "name", InputType: MultiColumnRepeatableFieldsetColumnTypeText},
		},
	}
	expectedSingleHint := "One row per line. Blank lines are ignored."
	if resolved := multiColumnRepeatableFieldsetBulkEditHintValueResolver(singleColumnSettings); resolved != expectedSingleHint {
		t.Errorf("bulkEditHintValue single column = %s, expected %s", resolved, expectedSingleHint)
	}
	customSettings := MultiColumnRepeatableFieldsetSettings{
		Columns: defaultSettings.Columns,
		BulkEdit: MultiColumnRepeatableFieldsetBulkEditSettings{
			Delimiter:        "|",
			ButtonIcon:       "ph-clipboard-text",
			ButtonLabel:      "Paste rows",
			TableButtonIcon:  "ph-table",
			TableButtonLabel: "Show table",
			TextAreaLabel:    "Contact rows",
			HintValue:        "One contact per line.",
		},
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditDelimiterResolver(customSettings); resolved != "|" {
		t.Errorf("bulkEditDelimiter custom = %s, expected |", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditButtonLabelResolver(customSettings); resolved != "Paste rows" {
		t.Errorf("bulkEditButtonLabel custom = %s, expected Paste rows", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditTableButtonLabelResolver(customSettings); resolved != "Show table" {
		t.Errorf("bulkEditTableButtonLabel custom = %s, expected Show table", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditTextAreaLabelResolver(customSettings); resolved != "Contact rows" {
		t.Errorf("bulkEditTextAreaLabel custom = %s, expected Contact rows", resolved)
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditHintValueResolver(customSettings); resolved != "One contact per line." {
		t.Errorf("bulkEditHintValue custom = %s, expected One contact per line.", resolved)
	}
	defaultMetadataJson := multiColumnRepeatableFieldsetBulkEditMetadataJsonBuilder(defaultSettings)
	expectedMetadataFragments := []string{
		`"delimiter":"="`,
		`"buttonLabel":"bulk edit"`,
		`"tableButtonLabel":"table view"`,
		`"buttonIcon":""`,
		`"tableButtonIcon":""`,
	}
	for _, expectedFragment := range expectedMetadataFragments {
		if !strings.Contains(defaultMetadataJson, expectedFragment) {
			t.Errorf("default bulk edit metadata missing fragment %s: %s", expectedFragment, defaultMetadataJson)
		}
	}
	customMetadataJson := multiColumnRepeatableFieldsetBulkEditMetadataJsonBuilder(customSettings)
	expectedCustomMetadataFragments := []string{
		`"delimiter":"|"`,
		`"buttonLabel":"Paste rows"`,
		`"tableButtonLabel":"Show table"`,
		`"buttonIcon":"ph-clipboard-text"`,
		`"tableButtonIcon":"ph-table"`,
	}
	for _, expectedFragment := range expectedCustomMetadataFragments {
		if !strings.Contains(customMetadataJson, expectedFragment) {
			t.Errorf("custom bulk edit metadata missing fragment %s: %s", expectedFragment, customMetadataJson)
		}
	}
	multiColumnSettings := MultiColumnRepeatableFieldsetSettings{
		Columns: []MultiColumnRepeatableFieldsetColumn{
			{Id: "name", InputType: MultiColumnRepeatableFieldsetColumnTypeText},
			{Id: "role", InputType: MultiColumnRepeatableFieldsetColumnTypeSelect},
			{Id: "enabled", InputType: MultiColumnRepeatableFieldsetColumnTypeCheckbox},
		},
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditDelimiterResolver(multiColumnSettings); resolved != "|" {
		t.Errorf("bulkEditDelimiter multi-column default = %s, expected |", resolved)
	}
	expectedMultiColumnHint := `One row per line. Separate the columns with "|". Blank lines are ignored.`
	if resolved := multiColumnRepeatableFieldsetBulkEditHintValueResolver(multiColumnSettings); resolved != expectedMultiColumnHint {
		t.Errorf("bulkEditHintValue multi-column default = %s, expected %s", resolved, expectedMultiColumnHint)
	}
	multiColumnMetadataJson := multiColumnRepeatableFieldsetBulkEditMetadataJsonBuilder(multiColumnSettings)
	if !strings.Contains(multiColumnMetadataJson, `"delimiter":"|"`) {
		t.Errorf("multi-column bulk edit metadata missing delimiter: %s", multiColumnMetadataJson)
	}
	explicitDelimiterSettings := MultiColumnRepeatableFieldsetSettings{
		Columns: multiColumnSettings.Columns,
		BulkEdit: MultiColumnRepeatableFieldsetBulkEditSettings{
			Delimiter: "=",
		},
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditDelimiterResolver(explicitDelimiterSettings); resolved != "=" {
		t.Errorf("bulkEditDelimiter explicit multi-column = %s, expected =", resolved)
	}
	unicodeDelimiterSettings := MultiColumnRepeatableFieldsetSettings{
		Columns: multiColumnSettings.Columns,
		BulkEdit: MultiColumnRepeatableFieldsetBulkEditSettings{
			Delimiter: "§",
		},
	}
	if resolved := multiColumnRepeatableFieldsetBulkEditDelimiterResolver(unicodeDelimiterSettings); resolved != "§" {
		t.Errorf("bulkEditDelimiter unicode = %s, expected §", resolved)
	}
	for _, conflictingDelimiter := range []string{"n", "\\", "||", "\n", "a\nb"} {
		conflictingSettings := MultiColumnRepeatableFieldsetSettings{
			Columns: multiColumnSettings.Columns,
			BulkEdit: MultiColumnRepeatableFieldsetBulkEditSettings{
				Delimiter: conflictingDelimiter,
			},
		}
		if resolved := multiColumnRepeatableFieldsetBulkEditDelimiterResolver(conflictingSettings); resolved != "|" {
			t.Errorf("bulkEditDelimiter conflicting %q = %s, expected |", conflictingDelimiter, resolved)
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
		`x-for="(row, rowIndex) in contactRows" :key="rowKeyResolver(row)"`,
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
		`add row`,
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
	if strings.Contains(renderedHtml, `toggleBulkEdit(noteRows)`) {
		t.Error("rendered html contains the bulk edit toggle without IsEnabled")
	}
	if strings.Contains(renderedHtml, `x-model:value="bulkEditText"`) {
		t.Error("rendered html contains the bulk edit textarea without IsEnabled")
	}
}

func TestMultiColumnRepeatableFieldsetRenderBulkEdit(t *testing.T) {
	component := MultiColumnRepeatableFieldset(MultiColumnRepeatableFieldsetSettings{
		Columns:         multiColumnRepeatableFieldsetTestColumns(),
		TwoWayStatePath: "contactRows",
		BulkEdit: MultiColumnRepeatableFieldsetBulkEditSettings{
			IsEnabled:        true,
			Delimiter:        "|",
			ButtonIcon:       "ph-clipboard-text",
			ButtonLabel:      "Paste rows",
			TableButtonIcon:  "ph-table",
			TableButtonLabel: "Show table",
			TextAreaLabel:    "Contact rows",
			HintValue:        "One contact per line.",
		},
	})
	var buffer bytes.Buffer
	if err := component.Render(context.Background(), &buffer); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`x-show="!isBulkEditing"`,
		`x-show="isBulkEditing"`,
		`toggleBulkEdit(contactRows)`,
		`replaceRowsFromBulkText(contactRows, $event.target.value)`,
		`bulkEditToggleLabelResolver`,
		`bulkEditToggleIconResolver`,
		`x-model:value="bulkEditText"`,
		`Contact rows`,
		`One contact per line.`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
	expectedMetadataFragments := []string{
		`&#34;delimiter&#34;:&#34;|&#34;`,
		`&#34;buttonLabel&#34;:&#34;Paste rows&#34;`,
		`&#34;tableButtonLabel&#34;:&#34;Show table&#34;`,
		`&#34;buttonIcon&#34;:&#34;ph-clipboard-text&#34;`,
		`&#34;tableButtonIcon&#34;:&#34;ph-table&#34;`,
	}
	for _, expectedFragment := range expectedMetadataFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing metadata fragment %s", expectedFragment)
		}
	}
}

func TestMultiColumnRepeatableFieldsetRendersColumnLabelStatePath(testingHandle *testing.T) {
	component := MultiColumnRepeatableFieldset(MultiColumnRepeatableFieldsetSettings{
		Columns: []MultiColumnRepeatableFieldsetColumn{
			{
				Id: "name", Label: "Name", LabelOneWayStatePath: "nameColumnLabel",
				InputType: MultiColumnRepeatableFieldsetColumnTypeText,
			},
			{
				Id: "role", Label: "Role", LabelOneWayStatePath: "roleColumnLabel",
				InputType: MultiColumnRepeatableFieldsetColumnTypeSelect,
			},
			{
				Id: "active", Label: "Active", LabelOneWayStatePath: "activeColumnLabel",
				InputType: MultiColumnRepeatableFieldsetColumnTypeCheckbox,
			},
		},
		TwoWayStatePath: "memberRows",
	})
	var buffer bytes.Buffer
	renderErr := component.Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("Render failed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-text="nameColumnLabel"`,
		`x-text="roleColumnLabel"`,
		`x-text="activeColumnLabel"`,
		`:aria-label="(activeColumnLabel) || &#34;Active&#34;"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
	for _, expectedCount := range []struct {
		fragment string
		count    int
	}{
		{`x-text="nameColumnLabel"`, 2},
		{`x-text="roleColumnLabel"`, 2},
		{`x-text="activeColumnLabel"`, 1},
	} {
		actualCount := strings.Count(renderedHtml, expectedCount.fragment)
		if actualCount < expectedCount.count {
			testingHandle.Errorf(
				"rendered html fragment %s count = %d, want at least %d",
				expectedCount.fragment, actualCount, expectedCount.count,
			)
		}
	}
}
