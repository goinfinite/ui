function randomNumberGenerator(rangeSize) {
  const maxUnbiasedValue = Math.floor(0x100000000 / rangeSize) * rangeSize;
  const randomValues = new Uint32Array(1);
  let rawRandomInteger;
  do {
    crypto.getRandomValues(randomValues);
    rawRandomInteger = randomValues[0];
  } while (rawRandomInteger >= maxUnbiasedValue);
  return rawRandomInteger % rangeSize;
}

const randomPasswordDefaultOptions = {
  length: 16,
  minLength: 6,
  maxLength: 64,
  includeLowercase: true,
  includeUppercase: true,
  includeNumbers: true,
  includeSpecialChars: true,
};

function randomPasswordCharsetsResolver(options) {
  const charsets = [];
  if (options.includeLowercase) {
    charsets.push("abcdefghijklmnopqrstuvwxyz");
  }
  if (options.includeUppercase) {
    charsets.push("ABCDEFGHIJKLMNOPQRSTUVWXYZ");
  }
  if (options.includeNumbers) {
    charsets.push("0123456789");
  }
  if (options.includeSpecialChars) {
    charsets.push("!@#$%^&*()_+-=.");
  }
  return charsets;
}

function createRandomPassword(options = {}) {
  const resolvedOptions = { ...randomPasswordDefaultOptions, ...options };
  const charsets = randomPasswordCharsetsResolver(resolvedOptions);
  if (charsets.length === 0) {
    throw new Error("RandomPasswordCharsetsEmpty");
  }

  const minLength = Math.max(
    Number(resolvedOptions.minLength) || charsets.length,
    charsets.length,
  );
  const maxLength = Math.max(
    Number(resolvedOptions.maxLength) || minLength,
    minLength,
  );
  const passwordLength = Math.min(
    Math.max(Number(resolvedOptions.length) || 16, minLength),
    maxLength,
  );

  const allChars = charsets.join("");
  const passwordChars = [];
  for (let charIndex = 0; charIndex < passwordLength; charIndex++) {
    passwordChars.push(allChars[randomNumberGenerator(allChars.length)]);
  }

  const availablePositions = [];
  for (let position = 0; position < passwordLength; position++) {
    availablePositions.push(position);
  }

  charsets.forEach((charset) => {
    const availableIndex = randomNumberGenerator(availablePositions.length);
    const targetPosition = availablePositions[availableIndex];
    availablePositions.splice(availableIndex, 1);
    passwordChars[targetPosition] =
      charset[randomNumberGenerator(charset.length)];
  });

  return passwordChars.join("");
}
