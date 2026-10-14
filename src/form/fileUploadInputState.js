UiToolset.RegisterAlpineState(() => {
  const fileUploadMessagesByLanguage = {
    en: {
      duplicate: "Already added",
      wrongType: "File type not accepted",
      tooLarge: "Larger than {limit}",
      limitReached: "File limit reached ({limit})",
      uploaded: "Uploaded",
      failed: "Upload failed",
      removeFile: "Remove {name}",
    },
    pt: {
      duplicate: "Já adicionado",
      wrongType: "Tipo de arquivo não aceito",
      tooLarge: "Maior que {limit}",
      limitReached: "Limite de arquivos atingido ({limit})",
      uploaded: "Enviado",
      failed: "Falha no envio",
      removeFile: "Remover {name}",
    },
    es: {
      duplicate: "Ya añadido",
      wrongType: "Tipo de archivo no aceptado",
      tooLarge: "Mayor que {limit}",
      limitReached: "Límite de archivos alcanzado ({limit})",
      uploaded: "Subido",
      failed: "Error al subir",
      removeFile: "Quitar {name}",
    },
    zh: {
      duplicate: "已添加",
      wrongType: "不接受此文件类型",
      tooLarge: "超过 {limit}",
      limitReached: "已达到文件数量上限（{limit}）",
      uploaded: "已上传",
      failed: "上传失败",
      removeFile: "移除 {name}",
    },
  };

  const fileTypeIconRules = [
    {
      icon: "ph-file-image",
      mimePattern: /^image\/[\w.+-]{1,100}$/,
      extensionPattern: /^(png|jpe?g|gif|webp|svg|avif|bmp|ico)$/,
    },
    {
      icon: "ph-file-video",
      mimePattern: /^video\/[\w.+-]{1,100}$/,
      extensionPattern: /^(mp4|mkv|webm|mov|avi)$/,
    },
    {
      icon: "ph-file-audio",
      mimePattern: /^audio\/[\w.+-]{1,100}$/,
      extensionPattern: /^(mp3|wav|ogg|flac|m4a)$/,
    },
    {
      icon: "ph-file-pdf",
      mimePattern: /^application\/pdf$/,
      extensionPattern: /^pdf$/,
    },
    {
      icon: "ph-file-zip",
      mimePattern:
        /^application\/(zip|gzip|x-tar|x-xz|x-7z-compressed|x-bzip2|vnd\.rar|zstd)$/,
      extensionPattern: /^(zip|tar|gz|tgz|xz|br|7z|bz2|rar|zst)$/,
    },
    {
      icon: "ph-file-code",
      mimePattern:
        /^(application\/(json|xml|javascript)|text\/(html|css|javascript))$/,
      extensionPattern: /^(json|xml|js|ts|html|css|go|py|sh|yaml|yml|toml)$/,
    },
    {
      icon: "ph-file-text",
      mimePattern: /^text\/[\w.+-]{1,100}$/,
      extensionPattern: /^(txt|md|csv|log)$/,
    },
  ];
  const fallbackFileTypeIcon = "ph-file";

  function fileExtensionReader(fileName) {
    const extensionSeparatorIndex = fileName.lastIndexOf(".");
    if (extensionSeparatorIndex === -1) {
      return "";
    }
    return fileName.slice(extensionSeparatorIndex + 1).toLowerCase();
  }

  function fileTypeIconResolver(candidateFile) {
    const fileExtension = fileExtensionReader(candidateFile.name);
    const fileMimeType = candidateFile.type.toLowerCase();
    for (const iconRule of fileTypeIconRules) {
      if (iconRule.mimePattern.test(fileMimeType)) {
        return iconRule.icon;
      }
      if (iconRule.extensionPattern.test(fileExtension)) {
        return iconRule.icon;
      }
    }
    return fallbackFileTypeIcon;
  }

  function fileSizeFormatter(sizeBytes) {
    const sizeUnits = ["B", "KiB", "MiB", "GiB", "TiB"];
    let scaledSize = sizeBytes;
    let unitIndex = 0;
    while (scaledSize >= 1024 && unitIndex < sizeUnits.length - 1) {
      scaledSize = scaledSize / 1024;
      unitIndex += 1;
    }
    if (unitIndex === 0) {
      return `${scaledSize} B`;
    }
    return `${scaledSize.toFixed(1)} ${sizeUnits[unitIndex]}`;
  }

  function isFileTypeAccepted(candidateFile, acceptedTypes) {
    const acceptTokens = acceptedTypes
      .split(",")
      .map((acceptToken) => acceptToken.trim().toLowerCase())
      .filter((acceptToken) => acceptToken !== "");
    if (acceptTokens.length === 0) {
      return true;
    }

    const fileName = candidateFile.name.toLowerCase();
    const fileMimeType = candidateFile.type.toLowerCase();
    for (const acceptToken of acceptTokens) {
      if (acceptToken === "*/*" || acceptToken === fileMimeType) {
        return true;
      }
      const isExtensionToken = acceptToken.startsWith(".");
      if (isExtensionToken && fileName.endsWith(acceptToken)) {
        return true;
      }
      const isMimeFamilyToken = acceptToken.endsWith("/*");
      const mimeFamilyPrefix = acceptToken.slice(0, -1);
      if (isMimeFamilyToken && fileMimeType.startsWith(mimeFamilyPrefix)) {
        return true;
      }
    }
    return false;
  }

  function fileUploadMessagesReader() {
    const browserLanguageCode = (navigator.language || "en")
      .toLowerCase()
      .split("-")[0];
    return (
      fileUploadMessagesByLanguage[browserLanguageCode] ??
      fileUploadMessagesByLanguage.en
    );
  }

  function fileEntryBuilder(candidateFile, entryId) {
    return {
      entryId: entryId,
      name: candidateFile.name,
      sizeBytes: candidateFile.size,
      sizeLabel: fileSizeFormatter(candidateFile.size),
      lastModified: candidateFile.lastModified,
      typeIcon: fileTypeIconResolver(candidateFile),
      status: "ready",
      progressPercent: 0,
    };
  }

  Alpine.data("fileUploadInput", (uploadLimits) => {
    const filesByEntryId = new Map();
    const isSingleFileMode = uploadLimits.maxFileCount === 1;

    return {
      fileEntries: [],
      rejectedFiles: [],
      dragDepth: 0,
      lastFileEntryId: 0,
      ownerForm: null,
      formListenersController: new AbortController(),

      fileUploadMessages: fileUploadMessagesReader(),

      get isDragActive() {
        return this.dragDepth > 0;
      },

      removeFileLabelResolver(fileName) {
        return this.fileUploadMessages.removeFile.replace(
          "{name}",
          () => fileName,
        );
      },

      get pendingFileEntries() {
        return this.fileEntries.filter(
          (fileEntry) => fileEntry.status !== "uploaded",
        );
      },

      init() {
        const ownerForm = this.$el.closest("form");
        if (!ownerForm) {
          console.error("FileUploadFormNotFound");
          return;
        }
        this.ownerForm = ownerForm;

        const listenerOptions = {
          signal: this.formListenersController.signal,
        };
        ownerForm.addEventListener(
          "reset",
          () => this.clearFileEntries(),
          listenerOptions,
        );
        ownerForm.addEventListener(
          "htmx:beforeSend",
          (requestEvent) => this.markUploadStarted(requestEvent),
          listenerOptions,
        );
        ownerForm.addEventListener(
          "htmx:xhr:progress",
          (progressEvent) => this.showUploadProgress(progressEvent),
          listenerOptions,
        );
        ownerForm.addEventListener(
          "htmx:afterRequest",
          (requestEvent) => this.showUploadResult(requestEvent),
          listenerOptions,
        );
      },

      destroy() {
        this.formListenersController.abort();
      },

      fileRejectionReasonResolver(candidateFile) {
        const rejectionMessages = this.fileUploadMessages;
        const isDuplicate = this.pendingFileEntries.some(
          (fileEntry) =>
            fileEntry.name === candidateFile.name &&
            fileEntry.sizeBytes === candidateFile.size &&
            fileEntry.lastModified === candidateFile.lastModified,
        );
        if (isDuplicate) {
          return rejectionMessages.duplicate;
        }
        if (!isFileTypeAccepted(candidateFile, uploadLimits.acceptedTypes)) {
          return rejectionMessages.wrongType;
        }
        const maxFileSizeBytes = uploadLimits.maxFileSizeBytes;
        if (maxFileSizeBytes > 0 && candidateFile.size > maxFileSizeBytes) {
          return rejectionMessages.tooLarge.replace(
            "{limit}",
            fileSizeFormatter(maxFileSizeBytes),
          );
        }
        const maxFileCount = uploadLimits.maxFileCount;
        const isFileCountReached =
          maxFileCount > 0 && this.pendingFileEntries.length >= maxFileCount;
        if (isFileCountReached && !isSingleFileMode) {
          return rejectionMessages.limitReached.replace(
            "{limit}",
            String(maxFileCount),
          );
        }
        return "";
      },

      removeReplaceableFileEntries() {
        for (const fileEntry of this.fileEntries) {
          const isReplaceable =
            fileEntry.status === "ready" || fileEntry.status === "failed";
          if (!isReplaceable) {
            continue;
          }
          filesByEntryId.delete(fileEntry.entryId);
        }
        this.fileEntries = this.fileEntries.filter(
          (fileEntry) =>
            fileEntry.status === "uploaded" || fileEntry.status === "uploading",
        );
      },

      syncFileInput() {
        const pendingFiles = new DataTransfer();
        for (const fileEntry of this.pendingFileEntries) {
          pendingFiles.items.add(filesByEntryId.get(fileEntry.entryId));
        }
        this.$refs.fileInput.files = pendingFiles.files;
      },

      addFiles(candidateFiles) {
        this.rejectedFiles = [];
        for (const candidateFile of candidateFiles) {
          const rejectionReason =
            this.fileRejectionReasonResolver(candidateFile);
          if (rejectionReason !== "") {
            this.rejectedFiles.push({
              name: candidateFile.name,
              reason: rejectionReason,
            });
            continue;
          }
          if (isSingleFileMode) {
            this.removeReplaceableFileEntries();
          }
          this.lastFileEntryId += 1;
          filesByEntryId.set(this.lastFileEntryId, candidateFile);
          this.fileEntries.push(
            fileEntryBuilder(candidateFile, this.lastFileEntryId),
          );
        }
        this.syncFileInput();
      },

      addPickedFiles() {
        const pickedFiles = Array.from(this.$refs.fileInput.files);
        this.addFiles(pickedFiles);
      },

      addDroppedFiles(dropEvent) {
        this.dragDepth = 0;
        this.addFiles(Array.from(dropEvent.dataTransfer.files));
      },

      removeFileEntry(entryId) {
        filesByEntryId.delete(entryId);
        this.fileEntries = this.fileEntries.filter(
          (fileEntry) => fileEntry.entryId !== entryId,
        );
        this.syncFileInput();
      },

      clearFileEntries() {
        filesByEntryId.clear();
        this.fileEntries = [];
        this.rejectedFiles = [];
        this.dragDepth = 0;
      },

      markUploadStarted(requestEvent) {
        if (requestEvent.detail.elt !== this.ownerForm) {
          return;
        }
        for (const fileEntry of this.pendingFileEntries) {
          fileEntry.status = "uploading";
          fileEntry.progressPercent = 0;
        }
      },

      showUploadProgress(progressEvent) {
        const progressDetail = progressEvent.detail;
        if (progressDetail.elt !== this.ownerForm) {
          return;
        }
        if (!progressDetail.lengthComputable) {
          return;
        }
        const uploadingEntries = this.fileEntries.filter(
          (fileEntry) => fileEntry.status === "uploading",
        );
        let uploadingBytes = 0;
        for (const fileEntry of uploadingEntries) {
          uploadingBytes += fileEntry.sizeBytes;
        }

        let bytesBeforeEntry = progressDetail.total - uploadingBytes;
        for (const fileEntry of uploadingEntries) {
          const entryLoadedBytes = Math.min(
            Math.max(progressDetail.loaded - bytesBeforeEntry, 0),
            fileEntry.sizeBytes,
          );
          bytesBeforeEntry += fileEntry.sizeBytes;

          let entryProgressPercent = 100;
          if (fileEntry.sizeBytes > 0) {
            entryProgressPercent = Math.floor(
              (entryLoadedBytes / fileEntry.sizeBytes) * 100,
            );
          }
          fileEntry.progressPercent = Math.max(
            fileEntry.progressPercent,
            entryProgressPercent,
          );
        }
      },

      showUploadResult(requestEvent) {
        if (requestEvent.detail.elt !== this.ownerForm) {
          return;
        }
        let uploadStatus = "failed";
        if (requestEvent.detail.successful) {
          uploadStatus = "uploaded";
        }
        for (const fileEntry of this.fileEntries) {
          if (fileEntry.status !== "uploading") {
            continue;
          }
          fileEntry.status = uploadStatus;
          if (uploadStatus === "uploaded") {
            fileEntry.progressPercent = 100;
            filesByEntryId.delete(fileEntry.entryId);
          }
        }
        this.syncFileInput();
      },
    };
  });
});
