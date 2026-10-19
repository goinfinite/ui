window.multiColumnRepeatableFieldsetDemoRandomRow = () => ({
  name: `Member ${Math.floor(Math.random() * 900 + 100)}`,
  role: ["owner", "maintainer", "member"][Math.floor(Math.random() * 3)],
  age: String(Math.floor(Math.random() * 60 + 18)),
  active: Math.random() > 0.5,
});
