#!/usr/bin/env node

/**
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
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

const fs = require('fs');
const path = require('path');
const yaml = require('yaml');

const API_DIR = path.join(__dirname, '../docs/static/api');
const OUTPUT_FILE = path.join(API_DIR, 'combined.yaml');

const API_FILES = [
  'application.yaml',
  'authentication.yaml',
  'branding.yaml',
  'flow-execution.yaml',
  'flow-management.yaml',
  'group.yaml',
  'healthcheck.yaml',
  'i18n.yaml',
  'idp.yaml',
  'notification-sender.yaml',
  'ou.yaml',
  'registration.yaml',
  'resource.yaml',
  'role.yaml',
  'user.yaml',
];

function mergeOpenAPISpecs() {
  console.log('🔄 Merging OpenAPI specifications...');

  // Base structure from the first spec
  const firstSpec = yaml.parse(
    fs.readFileSync(path.join(API_DIR, API_FILES[0]), 'utf8'),
  );

  const combined = {
    openapi: firstSpec.openapi || '3.0.3',
    info: {
      title: 'Thunder API Reference',
      version: '1.0',
      description:
        'Complete API reference for Thunder identity and access management.',
      license: firstSpec.info?.license || {
        name: 'Apache 2.0',
        url: 'https://www.apache.org/licenses/LICENSE-2.0.html',
      },
    },
    servers: firstSpec.servers || [],
    security: firstSpec.security || [],
    tags: [],
    paths: {},
    components: {
      schemas: {},
      securitySchemes: {},
      responses: {},
      parameters: {},
    },
  };

  // Process each API spec
  API_FILES.forEach((file) => {
    console.log(`  ➜ Processing ${file}...`);
    const specPath = path.join(API_DIR, file);
    const spec = yaml.parse(fs.readFileSync(specPath, 'utf8'));

    // Merge tags
    if (spec.tags) {
      spec.tags.forEach((tag) => {
        if (!combined.tags.find((t) => t.name === tag.name)) {
          combined.tags.push(tag);
        }
      });
    }

    // Merge paths
    if (spec.paths) {
      Object.entries(spec.paths).forEach(([path, pathItem]) => {
        if (combined.paths[path]) {
          // Merge methods if path exists
          combined.paths[path] = { ...combined.paths[path], ...pathItem };
        } else {
          combined.paths[path] = pathItem;
        }
      });
    }

    // Merge components
    if (spec.components) {
      if (spec.components.schemas) {
        combined.components.schemas = {
          ...combined.components.schemas,
          ...spec.components.schemas,
        };
      }
      if (spec.components.securitySchemes) {
        combined.components.securitySchemes = {
          ...combined.components.securitySchemes,
          ...spec.components.securitySchemes,
        };
      }
      if (spec.components.responses) {
        combined.components.responses = {
          ...combined.components.responses,
          ...spec.components.responses,
        };
      }
      if (spec.components.parameters) {
        combined.components.parameters = {
          ...combined.components.parameters,
          ...spec.components.parameters,
        };
      }
    }
  });

  // Sort tags alphabetically
  combined.tags.sort((a, b) => a.name.localeCompare(b.name));

  // Write the combined spec
  fs.writeFileSync(OUTPUT_FILE, yaml.stringify(combined), 'utf8');
  console.log(`✅ Combined API spec written to ${OUTPUT_FILE}`);
  console.log(
    `📊 Stats: ${combined.tags.length} tags, ${Object.keys(combined.paths).length} paths`,
  );
}

try {
  mergeOpenAPISpecs();
} catch (error) {
  console.error('❌ Error merging API specs:', error);
  process.exit(1);
}
