package uiStructural

func sidebarIsInlineAttachmentExpressionBuilder(
	attachmentModeTwoWayStatePath string,
) string {
	if attachmentModeTwoWayStatePath != "" {
		return attachmentModeTwoWayStatePath + " === '" + SidebarAttachmentModeInline + "'"
	}
	return "true"
}

func sidebarDynamicClassesBuilder(
	componentSettings SidebarSettings, attachmentMode string,
) string {
	dynamicClasses := "{"
	if componentSettings.IsVisibleTwoWayStatePath != "" {
		dynamicClasses += "'invisible': !" +
			componentSettings.IsVisibleTwoWayStatePath + ","
	}
	if componentSettings.IsCollapsedTwoWayStatePath != "" {
		dynamicClasses += "'!w-16': " + componentSettings.IsCollapsedTwoWayStatePath + ","
	}
	if componentSettings.AttachmentModeTwoWayStatePath != "" {
		dynamicClasses += "'!fixed top-0': " +
			componentSettings.AttachmentModeTwoWayStatePath + " === '" +
			SidebarAttachmentModeFixed + "',"
	}
	if componentSettings.IsOffCanvasTwoWayStatePath != "" {
		dynamicClasses += "'!absolute top-0 z-50': " +
			componentSettings.IsOffCanvasTwoWayStatePath + ","
	}
	if componentSettings.SideTwoWayStatePath != "" {
		dynamicClasses += "'left-0': " + componentSettings.SideTwoWayStatePath +
			" === '" + SidebarSideLeft + "',"
		dynamicClasses += "'right-0': " + componentSettings.SideTwoWayStatePath +
			" === '" + SidebarSideRight + "',"
	}
	return dynamicClasses + "}"
}

func sidebarWrapperDynamicClassesBuilder(
	componentSettings SidebarSettings, attachmentMode string,
) string {
	dynamicClasses := "{"
	if componentSettings.IsVisibleTwoWayStatePath != "" {
		dynamicClasses += "'invisible': !" +
			componentSettings.IsVisibleTwoWayStatePath + ","
	}
	if componentSettings.IsOffCanvasTwoWayStatePath != "" {
		dynamicClasses += "'!w-0': " + componentSettings.IsOffCanvasTwoWayStatePath + ","
	}
	if componentSettings.IsCollapsedTwoWayStatePath != "" {
		dynamicClasses += "'!w-16': " + componentSettings.IsCollapsedTwoWayStatePath + ","
	}
	if componentSettings.SideTwoWayStatePath != "" && attachmentMode == SidebarAttachmentModeInline {
		isInlineAttachmentExpression := sidebarIsInlineAttachmentExpressionBuilder(
			componentSettings.AttachmentModeTwoWayStatePath,
		)
		dynamicClasses += "'float-left': " + componentSettings.SideTwoWayStatePath +
			" === '" + SidebarSideLeft + "' && " + isInlineAttachmentExpression + ","
		dynamicClasses += "'float-right': " + componentSettings.SideTwoWayStatePath +
			" === '" + SidebarSideRight + "' && " + isInlineAttachmentExpression + ","
	}
	return dynamicClasses + "}"
}
