/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import {text, select, spinner, cancel} from '@clack/prompts';
import colors from 'picocolors';
import {join} from 'path';
import {createLogger} from '@thunder/logger';
import {getWorkspaceInfo, validateName, toKebabCase} from '../utils/workspace.js';
import {
  createFileFromTemplate,
  getTemplateDir,
  registerHandlebarsHelpers,
  fileExists,
  ensureDir,
} from '../utils/templates.js';

const logger = createLogger();

async function createFeature(): Promise<void> {
  const workspaceInfo = await getWorkspaceInfo();

  // Get feature type
  const featureType = await select({
    message: 'Feature type:',
    options: [
      {value: 'admin', label: 'Admin feature (thunder-admin-xxx)'},
      {value: 'gate', label: 'Gate feature (thunder-gate-xxx)'},
    ],
  });

  if (typeof featureType !== 'string') {
    cancel(colors.red('Operation cancelled.'));
    process.exit(1);
  }

  // Get feature name
  const name = await text({
    message: 'Feature name:',
    placeholder: 'user-management',
    validate: (value) => validateName(value, 'Feature'),
  });

  if (typeof name !== 'string') {
    cancel(colors.red('Operation cancelled.'));
    process.exit(1);
  }

  const featureName = toKebabCase(name);
  const packageName = `thunder-${featureType}-${featureName}`;
  const featureDir = join(workspaceInfo.packagePath!, packageName);

  // Check if feature already exists
  if (await fileExists(featureDir)) {
    cancel(colors.red(`Feature '${featureName}' already exists at ${featureDir}`));
    process.exit(1);
  }

  const s = spinner();
  s.start(colors.cyan(`Creating feature ${featureName}...`));

  try {
    registerHandlebarsHelpers();

    const context = {
      featureName,
      packageName,
      featureType,
    };

    // Create feature directory structure
    const directories = [
      'api',
      'components', 
      'config',
      'constants',
      'contexts',
      'data',
      'hooks',
      'models',
      'pages',
      'utils',
    ];

    for (const dir of directories) {
      await ensureDir(join(featureDir, dir));
    }

    // Create .gitkeep in data directory
    await createFileFromTemplate(
      join(getTemplateDir(), 'feature', 'gitkeep.hbs'),
      join(featureDir, 'data', '.gitkeep'),
      context,
    );

    // Create package.json (like thunder-shared-contexts but as feature)
    await createFileFromTemplate(
      join(getTemplateDir(), 'feature', 'package.json.hbs'),
      join(featureDir, 'package.json'),
      context,
    );

    // Create other config files
    const configFiles = [
      'tsconfig.json',
      'tsconfig.lib.json', 
      'tsconfig.spec.json',
      'tsconfig.eslint.json',
      'eslint.config.js',
      'vitest.config.ts',
      'rolldown.config.js',
      'prettier.config.js',
      '.editorconfig',
      '.gitignore',
      '.prettierignore',
    ];

    for (const file of configFiles) {
      await createFileFromTemplate(
        join(getTemplateDir(), 'feature', `${file}.hbs`),
        join(featureDir, file),
        context,
      );
    }

    // Create src/index.ts
    await ensureDir(join(featureDir, 'src'));
    await createFileFromTemplate(
      join(getTemplateDir(), 'feature', 'src', 'index.ts.hbs'),
      join(featureDir, 'src', 'index.ts'),
      context,
    );

    s.stop(colors.green(`✅ Feature '${featureName}' created successfully!`));

    logger.info(`Feature '${featureName}' created at ${featureDir}`);
    // eslint-disable-next-line no-console
    console.log();
    // eslint-disable-next-line no-console
    console.log(colors.cyan('Next steps:'));
    // eslint-disable-next-line no-console
    console.log(colors.gray('1. Install dependencies: cd to the feature directory and run pnpm install'));
    // eslint-disable-next-line no-console
    console.log(colors.gray('2. Start building your feature components and logic'));
  } catch (error) {
    s.stop(colors.red('❌ Failed to create feature'));
    logger.error('Feature creation failed:', {error: error instanceof Error ? error.message : String(error)});
    process.exit(1);
  }
}

export default createFeature;
