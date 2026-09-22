function randomNumberGenerator(rangeSize) {
  const randomValues = crypto.getRandomValues(new Uint32Array(1));
  const rawRandomInteger = randomValues[0];
  const numberWithinRange = rawRandomInteger % rangeSize;
  return numberWithinRange;
}

function createRandomPassword() {
  const passwordLength = 16;
  const letterChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
  const numberChars = "0123456789";
  const specialChars = "!@#$%^&*()_+";
  const allChars = letterChars + numberChars + specialChars;

  const passwordChars = [];
  for (let charIndex = 0; charIndex < passwordLength; charIndex++) {
    const randomPosition = randomNumberGenerator(allChars.length);
    passwordChars.push(allChars[randomPosition]);
  }

  const letterPosition = randomNumberGenerator(passwordLength);

  let numberPosition = randomNumberGenerator(passwordLength);
  while (numberPosition === letterPosition) {
    numberPosition = randomNumberGenerator(passwordLength);
  }

  let specialPosition = randomNumberGenerator(passwordLength);
  while (
    specialPosition === letterPosition ||
    specialPosition === numberPosition
  ) {
    specialPosition = randomNumberGenerator(passwordLength);
  }

  passwordChars[letterPosition] =
    letterChars[randomNumberGenerator(letterChars.length)];
  passwordChars[numberPosition] =
    numberChars[randomNumberGenerator(numberChars.length)];
  passwordChars[specialPosition] =
    specialChars[randomNumberGenerator(specialChars.length)];

  return passwordChars.join("");
}
