import eslint from '@eslint/js';
import nodePlugin from 'eslint-plugin-n';
import prettierConfig from 'eslint-config-prettier';
// @ts-ignore
import pluginPromise from 'eslint-plugin-promise';
import sonarjs from 'eslint-plugin-sonarjs';
import unicorn from 'eslint-plugin-unicorn';

import { defineConfig, globalIgnores } from 'eslint/config';

export default defineConfig(
  globalIgnores(['./**/node_modules/']),

  {
    linterOptions: {
      reportUnusedDisableDirectives: 'error',
      reportUnusedInlineConfigs: 'error',
    },
  },

  /** Core JS rules */
  eslint.configs.recommended,

  /** https://github.com/sindresorhus/eslint-plugin-unicorn */
  unicorn.configs['recommended'],

  /** https://github.com/eslint-community/eslint-plugin-n */
  nodePlugin.configs['flat/recommended-module'],

  /** https://github.com/SonarSource/eslint-plugin-sonarjs */
  sonarjs.configs.recommended,

  /** https://github.com/eslint-community/eslint-plugin-promise */
  pluginPromise.configs['flat/recommended'],

  {
    rules: {
      /** https://github.com/eslint-community/eslint-plugin-promise/blob/main/docs/rules/no-multiple-resolved.md */
      'promise/no-multiple-resolved': 'error',

      /** https://github.com/eslint-community/eslint-plugin-promise/blob/main/docs/rules/prefer-await-to-callbacks.md */
      'promise/prefer-await-to-callbacks': 'error',

      /** https://github.com/eslint-community/eslint-plugin-promise/blob/main/docs/rules/prefer-await-to-then.md */
      'promise/prefer-await-to-then': 'error',
    },
  },

  /** Turns off all rules that are unnecessary or might conflict with Prettier */
  prettierConfig,
);
