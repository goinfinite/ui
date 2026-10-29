function unixTimestampDateResolver(unixTimestampSeconds) {
  if (
    typeof unixTimestampSeconds !== "number" ||
    !Number.isFinite(unixTimestampSeconds)
  ) {
    throw new Error("InvalidUnixTimestamp");
  }
  const timestampDate = new Date(unixTimestampSeconds * 1000);
  const utcFullYear = timestampDate.getUTCFullYear();
  const isFourDigitYear = utcFullYear >= 0 && utcFullYear <= 9999;
  if (Number.isNaN(timestampDate.getTime()) || !isFourDigitYear) {
    throw new Error("InvalidUnixTimestamp");
  }
  return timestampDate;
}

function toIsoDate(unixTimestampSeconds) {
  const timestampDate = unixTimestampDateResolver(unixTimestampSeconds);
  return timestampDate.toISOString().slice(0, 10);
}

function toRfc1123(unixTimestampSeconds) {
  const timestampDate = unixTimestampDateResolver(unixTimestampSeconds);
  return timestampDate.toUTCString();
}
