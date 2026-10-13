UiToolset.RegisterAlpineState(() => {
  const fileUploadDemoUrlPattern = /(^|\/)file-upload-demo(-error)?$/;
  const fileUploadDemoFailureUrlPattern = /-error$/;
  const fileUploadDemoDurationMs = 2000;
  const fileUploadDemoProgressStepCount = 10;
  const multipartOverheadBytes = 512;

  function demoUploadBytesCounter(requestBody) {
    let uploadBytes = multipartOverheadBytes;
    if (!(requestBody instanceof FormData)) {
      return uploadBytes;
    }
    for (const [, fieldValue] of requestBody.entries()) {
      if (fieldValue instanceof File) {
        uploadBytes += fieldValue.size;
      }
    }
    return uploadBytes;
  }

  const originalOpen = XMLHttpRequest.prototype.open;
  const originalSend = XMLHttpRequest.prototype.send;

  XMLHttpRequest.prototype.open = function (method, requestUrl, ...openArgs) {
    this.isFileUploadDemoRequest = fileUploadDemoUrlPattern.test(
      String(requestUrl),
    );
    if (!this.isFileUploadDemoRequest) {
      return originalOpen.call(this, method, requestUrl, ...openArgs);
    }
    let demoResponseUrl = "assets/ui-logo.svg";
    if (fileUploadDemoFailureUrlPattern.test(String(requestUrl))) {
      demoResponseUrl = "assets/file-upload-demo-missing";
    }
    return originalOpen.call(this, "GET", demoResponseUrl, ...openArgs);
  };

  XMLHttpRequest.prototype.send = function (requestBody) {
    if (!this.isFileUploadDemoRequest) {
      return originalSend.call(this, requestBody);
    }
    const uploadTotalBytes = demoUploadBytesCounter(requestBody);
    let progressStepIndex = 0;
    const progressIntervalId = setInterval(() => {
      progressStepIndex += 1;
      const uploadedBytes = Math.round(
        (uploadTotalBytes * progressStepIndex) /
          fileUploadDemoProgressStepCount,
      );
      this.upload.dispatchEvent(
        new ProgressEvent("progress", {
          lengthComputable: true,
          loaded: uploadedBytes,
          total: uploadTotalBytes,
        }),
      );
      if (progressStepIndex < fileUploadDemoProgressStepCount) {
        return;
      }
      clearInterval(progressIntervalId);
      originalSend.call(this);
    }, fileUploadDemoDurationMs / fileUploadDemoProgressStepCount);
  };
});
