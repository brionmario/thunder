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

import {existsSync, readFileSync} from 'fs';
import {join, resolve} from 'path';

export interface WorkspaceInfo {
  isThunderWorkspace: boolean;
  frontendPath: string | null;
  packagePath: string | null;
  appsPath: string | null;
  currentWorkingDirectory: string;
}

/**
 * Detects if the current directory is part of a Thunder workspace
 * and returns relevant paths.
 */
export async function getWorkspaceInfo(): Promise<WorkspaceInfo> {
  const cwd = process.cwd();

  // Look for Thunder workspace indicators by walking up the directory tree
  let currentDir = cwd;
  let thunderRoot: string | null = null;

  while (currentDir !== '/' && currentDir !== '.') {
    // Check for Thunder root indicators - looking for a frontend directory with nx.json
    const frontendDir = join(currentDir, 'frontend');
    const frontendNxJson = join(frontendDir, 'nx.json');
    const frontendPackageJson = join(frontendDir, 'package.json');

    if (existsSync(frontendDir) && existsSync(frontendNxJson) && existsSync(frontendPackageJson)) {
      try {
        const packageJson = JSON.parse(readFileSync(frontendPackageJson, 'utf8'));
        if (packageJson.name?.includes('thunder')) {
          thunderRoot = currentDir;
          break;
        }
      } catch {
        // Continue searching
      }
    }

    // Also check if we're already in the frontend directory
    const packageJsonPath = join(currentDir, 'package.json');
    const nxJsonPath = join(currentDir, 'nx.json');

    if (existsSync(packageJsonPath) && existsSync(nxJsonPath)) {
      try {
        const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf8'));
        if (packageJson.name === '@thunder/frontend') {
          thunderRoot = resolve(currentDir, '..');
          break;
        }
      } catch {
        // Continue searching
      }
    }

    currentDir = resolve(currentDir, '..');
  }

  if (!thunderRoot) {
    return {
      isThunderWorkspace: false,
      frontendPath: null,
      packagePath: null,
      appsPath: null,
      currentWorkingDirectory: cwd,
    };
  }

  const frontendPath = join(thunderRoot, 'frontend');
  const packagePath = join(frontendPath, 'packages');
  const appsPath = join(frontendPath, 'apps');

  return {
    isThunderWorkspace: true,
    frontendPath,
    packagePath: existsSync(packagePath) ? packagePath : null,
    appsPath: existsSync(appsPath) ? appsPath : null,
    currentWorkingDirectory: cwd,
  };
}

/**
 * Convert a string to PascalCase (for React components).
 */
export function toPascalCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1 $2') // Split camelCase
    .replace(/[_-]/g, ' ') // Replace underscores and hyphens with spaces
    .split(' ')
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join('');
}

/**
 * Convert a string to camelCase (for functions, variables).
 */
export function toCamelCase(str: string): string {
  const pascalCase = toPascalCase(str);
  return pascalCase.charAt(0).toLowerCase() + pascalCase.slice(1);
}

/**
 * Convert a string to kebab-case (for directories, files).
 */
export function toKebabCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1-$2') // Split camelCase
    .replace(/[_\s]+/g, '-') // Replace underscores and spaces with hyphens
    .toLowerCase();
}

/**
 * Convert a string to CONSTANT_CASE (for constants).
 */
export function toConstantCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1_$2') // Split camelCase
    .replace(/[_-\s]+/g, '_') // Replace hyphens and spaces with underscores
    .toUpperCase();
}

/**
 * Validates that a name is suitable for a package, component, etc.
 */
export function validateName(name: string, type: string): string | undefined {
  if (!name || name.trim().length === 0) {
    return `${type} name cannot be empty`;
  }

  const trimmed = name.trim();

  // Check for valid characters
  if (!/^[a-zA-Z][a-zA-Z0-9_-]*$/.test(trimmed)) {
    return `${type} name must start with a letter and contain only letters, numbers, underscores, and hyphens`;
  }

  // Check length
  if (trimmed.length > 50) {
    return `${type} name must be 50 characters or less`;
  }

  // Check for reserved words
  const reserved = ['index', 'src', 'dist', 'build', 'node_modules', 'package', 'test', '__tests__'];
  if (reserved.includes(trimmed.toLowerCase())) {
    return `${type} name '${trimmed}' is reserved`;
  }

  return undefined;
}
