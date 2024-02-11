module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    // "plugin:react-hooks/recommended",
    'airbnb',
    'airbnb-typescript',
    // "plugin:react/jsx-runtime",
    'prettier',
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    project: './tsconfig.json',
    tsconfigRootDir: __dirname,
  },
  plugins: [
    '@typescript-eslint',
    // "react"
  ],
  reportUnusedDisableDirectives: true,
  rules: {
    'import/prefer-default-export': 'off',
    'import/extensions': ['error', 'always', { ignorePackages: true }],

    '@typescript-eslint/naming-convention': [
      'error',
      { selector: 'enum', format: ['PascalCase'] },
      { selector: 'enumMember', format: ['PascalCase'] },
      {
        // everything should be camelCase by default
        selector: 'default',
        format: ['camelCase'],
      },
      {
        // functions can also be PascalCase (but have to be globally accessed)
        selector: 'variable',
        format: ['PascalCase', 'camelCase'],
        modifiers: ['global'],
        types: ['function'],
      },
      {
        // other variables can be UPPER_CASE (but have to be globally accessed)
        selector: 'variable',
        format: ['UPPER_CASE', 'camelCase'],
        modifiers: ['global'],
      },
      {
        // unused values must start with an underscore
        selector: 'parameter',
        format: ['camelCase'],
        leadingUnderscore: 'require',
        modifiers: ['unused'],
      },
      {
        // types must be PascalCase
        selector: 'typeLike',
        format: ['PascalCase'],
      },
    ],
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/no-unused-vars': 'error',
    '@typescript-eslint/lines-between-class-members': 'off',
    '@typescript-eslint/ban-types': 'error',
    '@typescript-eslint/ban-ts-comment': [
      'error',
      { 'ts-expect-error': 'allow-with-description' },
    ],
    '@typescript-eslint/no-unsafe-return': 'error',
  },
  overrides: [
    {
      // prevent trying to type-check javascript files
      // https://typescript-eslint.io/linting/troubleshooting#fixing-the-error
      extends: ['plugin:@typescript-eslint/disable-type-checked'],
      files: ['./**/*.{js,mjs}'],
    },
    {
      // cdk stacks
      files: ['{bin,lib}/**/*.ts'],
      rules: {
        'no-new': 'off',
      },
    },
    {
      env: {
        node: true,
      },
      files: ['.eslintrc.{js,cjs}'],
      parserOptions: {
        sourceType: 'script',
      },
    },
  ],
};
