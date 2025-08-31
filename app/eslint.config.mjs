import js from '@eslint/js';
import globals from 'globals';
import tseslint from 'typescript-eslint';

import reactPlugin from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import jsxA11y from 'eslint-plugin-jsx-a11y';

import importPlugin from 'eslint-plugin-import';
import boundaries from 'eslint-plugin-boundaries';
import sortPlugin from 'eslint-plugin-simple-import-sort';
import unusedImports from 'eslint-plugin-unused-imports';
import checkFile from 'eslint-plugin-check-file';

import prettierPlugin from 'eslint-plugin-prettier';
import prettierFlat from 'eslint-config-prettier/flat';
import * as prettierRc from './prettier.config.mjs';

import testingLibrary from 'eslint-plugin-testing-library';
import jestDom from 'eslint-plugin-jest-dom';

import pluginQuery from '@tanstack/eslint-plugin-query';

import pluginRouter from '@tanstack/eslint-plugin-router'

const ignoreConfig = { ignores: ['dist', 'coverage', 'target', '**/*.d.ts'] };

const tsPreset = tseslint.configs.recommended;
const jsPreset = js.configs.recommended;
const prettier = prettierFlat;

const projectConfig = {
  files: ['./src/**/*.{ts,tsx,js,jsx}'],

  languageOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    globals: globals.browser,
    parserOptions: { project: './tsconfig.app.json' },
  },

  plugins: {
    react: reactPlugin,
    'react-hooks': reactHooks,
    'react-refresh': reactRefresh,
    'jsx-a11y': jsxA11y,
    import: importPlugin,
    boundaries,
    'simple-import-sort': sortPlugin,
    'unused-imports': unusedImports,
    'check-file': checkFile,
    prettier: prettierPlugin,
  },

  settings: {
    react: { version: 'detect' },
    'import/resolver': {
      typescript: { project: './tsconfig.app.json' },
      node: true,
    },
    // Groups by type, applies the pattern to each folder or file the applied
    // pattern creates a public module, i.e. `src/components/button` and every
    // descendant is considered private.
    //
    // Assets are not included as boundary elements since they are individual
    // files that don't need architectural protection and should be importable
    // directly.
    'boundaries/elements': [
      { type: 'features', pattern: 'src/features/*', mode: 'folder', capture: ['feature'] },
      { type: 'components', pattern: 'src/components/*', mode: 'folder' },
      { type: 'entities', pattern: 'src/entities/*', mode: 'folder' },
      { type: 'contexts', pattern: 'src/contexts/*', mode: 'folder' },
      { type: 'providers', pattern: 'src/providers/*', mode: 'folder' },
      { type: 'theme', pattern: 'src/theme', mode: 'folder' },
      { type: 'utilities', pattern: 'src/utilities/*', mode: 'folder' },
      { type: 'app', pattern: 'src/app', mode: 'folder' },
      { type: 'config', pattern: ['./config.js', './config.d.ts', 'config.ts'], mode: 'file' },
    ],
    'boundaries/ignore': [
      '**/*.png',
      '**/*.jpg',
      '**/*.jpeg',
      '**/*.svg',
      '**/*.gif',
      '**/*.webp',
      '**/*.ico',
      'src/assets/**',
      'node_modules/**',
      'public/**',
    ],
  },

  rules: {
    ...reactPlugin.configs.recommended.rules,
    ...reactHooks.configs.recommended.rules,

    ...jsxA11y.configs.recommended.rules,
    ...boundaries.configs.strict.rules,

    'prettier/prettier': ['error', prettierRc.default ?? prettierRc],

    'eol-last': ['error', 'always'],

    'react/react-in-jsx-scope': 'off',
    'react-refresh/only-export-components': ['warn', { allowConstantExport: true }],

    'simple-import-sort/imports': 'error',
    'simple-import-sort/exports': 'error',

    'unused-imports/no-unused-imports': 'error',
    'no-unused-vars': [
      'warn',
      { varsIgnorePattern: '^_', argsIgnorePattern: '^_', ignoreRestSiblings: true }
    ],
    '@typescript-eslint/no-unused-vars': [
      'warn',
      { varsIgnorePattern: '^_', argsIgnorePattern: '^_', ignoreRestSiblings: true }
    ],

    'import/no-extraneous-dependencies': [
      'error',
      {
        devDependencies: ['**/*.{test,spec}.*', 'vite.config.*'],
        includeInternal: false,
        includeTypes: true
      },
    ],
    'import/first': 'error',
    'import/newline-after-import': 'error',
    'import/no-duplicates': 'error',
    'import/no-cycle': ['error', { maxDepth: 1 }],

    // Prevents disallowed cross-layer imports (e.g., features importing other
    // features or app).
    //
    // See
    // <https://github.com/javierbrea/eslint-plugin-boundaries/tree/master/docs/rules>
    // for more details.
    'boundaries/element-types': ['error', {
      default: 'disallow',
      message: "This file is unallowed to import this dependency.",
      rules: [
        { from: ['entities'], allow: ['generated'] },
        { from: ['features'], allow: ['components', 'utilities', 'providers', 'contexts', 'theme', 'config', 'entities', 'generated'] },
        { from: ['components'], allow: ['components'] },
        { from: ['contexts'], allow: ['theme', 'config', 'utilities', 'entities', 'generated'] },
        { from: ['providers'], allow: ['contexts', 'config', 'theme', 'components', 'entities', 'generated'] },
        { from: ['app'], allow: ['features', 'app', 'providers', 'entities', 'contexts', 'entities', 'theme', 'components', 'generated'] },
      ],
    }],
    // Cross-folder importing only through allowed public entry-points
    // (`index.ts[x]`, `index.js[x]`). Assets are excluded since they are
    // individual files that cannot have index files.
    'boundaries/entry-point': ['error', {
      default: 'disallow',
      message: "Import via the public API (index.ts, index.js, index.tsx, or index.jsx)",
      rules: [{
        target: ['features', 'entities', 'providers', 'config', 'utilities', 'components', 'contexts', 'theme', 'generated'],
        allow: ['index.ts', 'index.tsx', 'index.js', 'index.jsx', 'config.js', 'config.ts', 'config.d.ts'],
        message: "Import cross-folder via the public API (index.{ts, js, tsx, jsx})"
      }],
    }],
    // The following are rules for private elements:
    //
    // * An element becomes private when it is under another element.
    // * Private elements can't be used by anyone except its parent (and any
    //   other descendant of the parent when allowUncles option is enabled).
    // * Private elements can import public elements.
    // * Private elements can import another private element when both have the
    //   same parent ("sibling elements").
    // * Private elements can import another private element if it is a direct
    //   child of a common ancestor, and the `allowUncles` option is enabled.
    'boundaries/no-private': ['error', {
      allowUncles: true
    }],

    '@typescript-eslint/consistent-type-imports': [
    'error',
      {
        fixStyle: 'inline-type-imports',
      }
     ],

    '@typescript-eslint/naming-convention': ['warn',
      { selector: 'typeLike', format: ['PascalCase'] },
      {
        selector: 'variable',
        types: ['boolean'],
        format: ['PascalCase'],
        prefix: ['is', 'has', 'should', 'can', 'did', 'will'],
      },
      {
        selector: 'function',
        format: ['camelCase'],
        filter: { regex: '^[A-Z]', match: false },
      },
    ],
    'react/jsx-pascal-case': ['error', { allowAllCaps: true }],
    '@typescript-eslint/no-shadow': ['error'],
    'no-shadow': 'off',

    'check-file/folder-naming-convention': ['error', { '**/*': 'SNAKE_CASE' }],
    'check-file/filename-naming-convention': ['error', {
      'src/**/!(main|routes|index).{tsx,jsx}': 'PASCAL_CASE',
      'src/**/(!*.d).{ts,js}': 'CAMEL_CASE',
    }],

    'react/function-component-definition': ['error', {
      namedComponents: 'function-declaration',
      unnamedComponents: 'function-expression',
    }],
  },
};

const testOverride = {
  files: ['**/*.{test,spec}.{ts,tsx,js,jsx}'],
  plugins: { 'testing-library': testingLibrary, 'jest-dom': jestDom },
  rules: {
    ...testingLibrary.configs.react.rules,
    ...jestDom.configs.recommended.rules,
    'testing-library/no-node-access': 'off',
  },
  languageOptions: { globals: { jest: true } },
};

const buildScriptsOverride = {
  files: ['vite.config.{js,ts}', '*.cjs'],
  languageOptions: { parserOptions: { project: null } },
  rules: { '@typescript-eslint/no-var-requires': 'off' },
};

const scriptsOverride = {
  files: ['./scripts/**/*.{ts,tsx,js,jsx}'],
  languageOptions: { globals: globals.node },
};

export default tseslint.config(
  ignoreConfig,
  tsPreset,
  jsPreset,
  prettier,
  projectConfig,
  pluginQuery.configs['flat/recommended'],
  pluginRouter.configs['flat/recommended'],
  testOverride,
  buildScriptsOverride,
  scriptsOverride,
);
