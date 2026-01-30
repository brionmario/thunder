#!/usr/bin/env node

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

/* eslint-disable @thunder/copyright-header */

import {Command} from 'commander';
import {intro, outro, cancel} from '@clack/prompts';
import colors from 'picocolors';
import {createLogger} from '@thunder/logger';
import createFeature from './commands/createFeature';
import createPackage from './commands/createPackage';
import {getWorkspaceInfo} from './utils/workspace';

const logger = createLogger();

const program: Command = new Command();

async function main(): Promise<void> {
  console.clear();

  intro(colors.bgBlue(colors.black(' ⚡ Thunder Create ')));

  // Check if we're in a Thunder workspace
  const workspaceInfo = await getWorkspaceInfo();
  if (!workspaceInfo.isThunderWorkspace) {
    cancel(colors.red('This command must be run from a Thunder workspace.'));
    process.exit(1);
  }

  program.name('thunder-create').description('CLI scaffolding tool for ⚡ Thunder frontends').version('0.0.0');

  program.command('feature').description('Create a new feature module').action(createFeature);

  program.command('package').description('Create a new shared package').action(createPackage);

  await program.parseAsync();

  outro(colors.green('✅ Done! Happy coding!'));
}

export default main;

if (import.meta.url === `file://${process.argv[1]}`) {
  main().catch((error) => {
    logger.error('CLI execution failed:', error);
    process.exit(1);
  });
}
