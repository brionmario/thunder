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

import Handlebars from 'handlebars';
import fsExtra from 'fs-extra';
import {join, dirname} from 'path';
import {readFileSync} from 'fs';
import {fileURLToPath} from 'url';

const {existsSync, mkdirSync, writeFileSync} = fsExtra;

export type TemplateContext = Record<string, any>;

/**
 * Renders a template string with the given context.
 */
export function renderTemplate(templateContent: string, context: TemplateContext): string {
  const template = Handlebars.compile(templateContent);
  return template(context);
}

/**
 * Renders a template file with the given context.
 */
export function renderTemplateFile(templatePath: string, context: TemplateContext): string {
  if (!existsSync(templatePath)) {
    throw new Error(`Template file not found: ${templatePath}`);
  }

  const templateContent = readFileSync(templatePath, 'utf8');
  return renderTemplate(templateContent, context);
}

/**
 * Creates a file from a template with the given context.
 */
export async function createFileFromTemplate(
  templatePath: string,
  outputPath: string,
  context: TemplateContext,
): Promise<void> {
  const content = renderTemplateFile(templatePath, context);

  // Ensure the directory exists
  const dir = dirname(outputPath);
  if (!existsSync(dir)) {
    mkdirSync(dir, {recursive: true});
  }

  writeFileSync(outputPath, content, 'utf8');
}

/**
 * Creates multiple files from templates.
 */
export async function createFilesFromTemplates(
  templates: {
    templatePath: string;
    outputPath: string;
    context: TemplateContext;
  }[],
): Promise<void> {
  for (const template of templates) {
    await createFileFromTemplate(template.templatePath, template.outputPath, template.context);
  }
}

/**
 * Checks if a file or directory already exists.
 */
export function fileExists(path: string): boolean {
  return existsSync(path);
}

/**
 * Creates a directory if it doesn't exist.
 */
export function ensureDir(path: string): void {
  if (!existsSync(path)) {
    mkdirSync(path, {recursive: true});
  }
}

/**
 * Gets the template directory path.
 */
export function getTemplateDir(): string {
  // For the linked global package, find the package root by looking for package.json
  let currentDir = dirname(fileURLToPath(import.meta.url));
  
  // Go up directories until we find package.json
  while (currentDir !== '/' && !existsSync(join(currentDir, 'package.json'))) {
    currentDir = dirname(currentDir);
  }
  
  // If we found package.json, templates should be in dist/templates or src/templates
  if (existsSync(join(currentDir, 'package.json'))) {
    const distTemplates = join(currentDir, 'dist', 'templates');
    const srcTemplates = join(currentDir, 'src', 'templates');
    
    if (existsSync(distTemplates)) {
      return distTemplates;
    }
    if (existsSync(srcTemplates)) {
      return srcTemplates;
    }
  }
  
  // Fallback to current directory relative
  return join(dirname(fileURLToPath(import.meta.url)), '..', 'templates');
}

/**
 * Registers common Handlebars helpers for use in templates.
 */
export function registerHandlebarsHelpers(): void {
  // Helper to convert to PascalCase
  Handlebars.registerHelper('pascalCase', (str: string) =>
    str
      .replace(/([a-z])([A-Z])/g, '$1 $2')
      .replace(/[_-]/g, ' ')
      .split(' ')
      .filter(Boolean)
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(''),
  );

  // Helper to convert to camelCase
  Handlebars.registerHelper('camelCase', (str: string) => {
    const pascalCase = str
      .replace(/([a-z])([A-Z])/g, '$1 $2')
      .replace(/[_-]/g, ' ')
      .split(' ')
      .filter(Boolean)
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join('');
    return pascalCase.charAt(0).toLowerCase() + pascalCase.slice(1);
  });

  // Helper to convert to kebab-case
  Handlebars.registerHelper('kebabCase', (str: string) =>
    str
      .replace(/([a-z])([A-Z])/g, '$1-$2')
      .replace(/[_\s]+/g, '-')
      .toLowerCase(),
  );

  // Helper to convert to CONSTANT_CASE
  Handlebars.registerHelper('constantCase', (str: string) =>
    str
      .replace(/([a-z])([A-Z])/g, '$1_$2')
      .replace(/[_-\s]+/g, '_')
      .toUpperCase(),
  );

  // Helper for conditional inclusion
  Handlebars.registerHelper('if_eq', function (this: any, a: any, b: any, options: any) {
    return a === b ? options.fn(this) : options.inverse(this);
  });

  // Helper for array inclusion
  Handlebars.registerHelper('if_includes', function (this: any, array: any, item: any, options: any) {
    return array?.includes(item) ? options.fn(this) : options.inverse(this);
  });
}
