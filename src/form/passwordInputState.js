UiToolset.RegisterAlpineState(() => {
  const passwordInputDefaultRules = {
    minLength: 6,
    maxLength: 64,
    length: 16,
    includeLowercase: true,
    includeUppercase: true,
    includeNumbers: true,
    includeSpecialChars: true,
  };

  Alpine.data("passwordInput", (rulesScriptId) => ({
    rules: { ...passwordInputDefaultRules },
    isPasswordVisible: false,
    passwordStrengthPercentage: 0,
    passwordStrengthCriteria: {
      isLongEnough: false,
      hasNumbers: false,
      hasUppercaseChars: false,
      hasLowercaseChars: false,
      hasSpecialChars: false,
    },

    init() {
      if (!rulesScriptId) {
        return;
      }
      try {
        const rulesElement = document.getElementById(rulesScriptId);
        this.rules = {
          ...this.rules,
          ...JSON.parse(rulesElement?.textContent || "{}"),
        };
      } catch (parseError) {
        console.error(`PasswordInputInvalidRulesJson: ${parseError.message}`);
      }
    },

    updatePasswordStrength(passwordValue) {
      const passwordText = String(passwordValue ?? "");
      const criteria = {
        isLongEnough:
          passwordText.length >= this.rules.minLength &&
          passwordText.length <= this.rules.maxLength,
        hasNumbers: false,
        hasUppercaseChars: false,
        hasLowercaseChars: false,
        hasSpecialChars: false,
      };

      let criteriaTotal = 1;
      let criteriaPassed = criteria.isLongEnough ? 1 : 0;

      if (this.rules.includeNumbers) {
        criteriaTotal++;
        criteria.hasNumbers = /[0-9]/.test(passwordText);
        criteriaPassed += criteria.hasNumbers ? 1 : 0;
      }
      if (this.rules.includeUppercase) {
        criteriaTotal++;
        criteria.hasUppercaseChars = /[A-Z]/.test(passwordText);
        criteriaPassed += criteria.hasUppercaseChars ? 1 : 0;
      }
      if (this.rules.includeLowercase) {
        criteriaTotal++;
        criteria.hasLowercaseChars = /[a-z]/.test(passwordText);
        criteriaPassed += criteria.hasLowercaseChars ? 1 : 0;
      }
      if (this.rules.includeSpecialChars) {
        criteriaTotal++;
        criteria.hasSpecialChars = /[!@#$%^&*()_+=.-]/.test(passwordText);
        criteriaPassed += criteria.hasSpecialChars ? 1 : 0;
      }

      this.passwordStrengthCriteria = criteria;
      this.passwordStrengthPercentage = Math.round(
        (criteriaPassed / criteriaTotal) * 100,
      );
    },

    generateRandomPassword() {
      return UiToolset.CreateRandomPassword(this.rules);
    },
  }));
});
