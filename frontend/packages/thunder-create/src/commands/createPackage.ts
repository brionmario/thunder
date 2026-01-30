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

async function createPackage(): Promise<void> {
  const workspaceInfo = await getWorkspaceInfo();

  // Get package type
  const packageType = await select({
    message: 'Package type:',
    options: [
      {value: 'javascript', label: 'JavaScript package (like thunder-logger)'},
      {value: 'react', label: 'React package (like thunder-shared-contexts)'},
    ],
  });

  if (typeof packageType !== 'string') {
    cancel(colors.red('Operation cancelled.'));
    process.exit(1);
  }

  // Get package name
  const name = await text({
    message: 'Package name:',
    placeholder: 'shared-ui-components',
    validate: (value) => validateName(value, 'Package'),
  });

  if (typeof name !== 'string') {
    cancel(colors.red('Operation cancelled.'));
    process.exit(1);
  }

  const packageName = toKebabCase(name);
  const fullPackageName = `@thunder/${packageName}`;
  const packageDir = join(workspaceInfo.packagePath!, `thunder-${packageName}`);

  // Check if package already exists
  if (await fileExists(packageDir)) {
    cancel(colors.red(`Package '${packageName}' already exists at ${packageDir}`));
    process.exit(1);
  }

  const s = spinner();
  s.start(colors.cyan(`Creating package ${packageName}...`));

  try {
    registerHandlebarsHelpers();

    const context = {
      packageName: packageName,
      fullPackageName: fullPackageName,
      packageType,
      isReactPackage: packageType === 'react',
    };

    // Create package directory structure
    await ensureDir(join(packageDir, 'src'));

    // Create package.json based on type
    const templateBase = packageType === 'react' ? 'package-react' : 'package-js';
    await createFileFromTemplate(
      join(getTemplateDir(), 'package', `${templateBase}`, 'package.json.hbs'),
      join(packageDir, 'package.json'),
      context,
    );

    // Create common config files based on package type
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
        join(getTemplateDir(), 'package', `${templateBase}`, `${file}.hbs`),
        join(packageDir, file),
        context,
      );
    }

    // Create src/index.ts
    await createFileFromTemplate(
      join(getTemplateDir(), 'package', `${templateBase}`, 'src', 'index.ts.hbs'),
      join(packageDir, 'src', 'index.ts'),
      context,
    );

    // Create README.md
    await createFileFromTemplate(
      join(getTemplateDir(), 'package', `${templateBase}`, 'README.md.hbs'),
      join(packageDir, 'README.md'),
      context,
    );

    s.stop(colors.green(`✅ Package '${packageName}' created successfully!`));

    logger.info(`Package '${packageName}' created at ${packageDir}`);
    // eslint-disable-next-line no-console
    console.log();
    // eslint-disable-next-line no-console
    console.log(colors.cyan('Next steps:'));
    // eslint-disable-next-line no-console
    console.log(colors.gray('1. Install dependencies: cd to the package directory and run pnpm install'));
    // eslint-disable-next-line no-console
    console.log(colors.gray('2. Start building your package functionality'));
    // eslint-disable-next-line no-console
    console.log(colors.gray('3. Run tests: pnpm test'));
    // eslint-disable-next-line no-console
    console.log(colors.gray('4. Build package: pnpm build'));
  } catch (error) {
    s.stop(colors.red('❌ Failed to create package'));
    logger.error('Package creation failed:', {error: error instanceof Error ? error.message : String(error)});
    process.exit(1);
  }
}

export default createPackage;
