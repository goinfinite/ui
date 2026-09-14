function isNonBlankString(value) {
  return typeof value === "string" && value.trim() !== "";
}

function resolveApiResponseMessage(apiResponse, responseOutcome) {
  const readableMessage = apiResponse?.readableMessage;
  if (isNonBlankString(readableMessage)) {
    return readableMessage;
  }

  const humanReadableMessage = apiResponse?.humanReadableMessage;
  const outcomeMessage = humanReadableMessage?.[responseOutcome];
  if (isNonBlankString(outcomeMessage)) {
    return outcomeMessage;
  }

  const stringBodyMessage = apiResponse?.body;
  if (isNonBlankString(stringBodyMessage)) {
    return stringBodyMessage;
  }

  return "";
}

function resolveApiResponseOutcome(httpStatusCode) {
  if (httpStatusCode === 207) {
    return "partialSuccess";
  }
  if (httpStatusCode >= 400) {
    return "error";
  }
  return "success";
}

function resolveApiResponseDisplay(apiResponse, httpStatusCode) {
  const responseOutcome = resolveApiResponseOutcome(httpStatusCode);
  const responseMessage = resolveApiResponseMessage(
    apiResponse,
    responseOutcome,
  );
  return { message: responseMessage, outcome: responseOutcome };
}
