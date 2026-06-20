module.exports = {
  testEnvironment: 'node',
  testMatch: [
    '**/tests/**/*.test.[jt]s?(x)'
  ],
  testPathIgnorePatterns: [
    '/node_modules/',
    '/.aether/',
    '/packages/'
  ]
};
