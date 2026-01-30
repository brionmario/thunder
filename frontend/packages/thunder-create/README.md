# @thunder/create

CLI scaffolding tool for ⚡ Thunder frontends.

## Installation

```bash
# In your Thunder workspace
pnpm install
pnpm build

# Run from anywhere in the workspace
npx thunder-create
```

## Usage

Run the CLI tool from anywhere within a Thunder workspace:

```bash
npx thunder-create
```

Then follow the interactive prompts to select what you want to create.

### Available Commands

- **component** - Create a new React component
- **feature** - Create a new feature module  
- **package** - Create a new shared package
- **app** - Create a new application
- **api** - Create API functions and hooks
- **hook** - Create a custom React hook
- **util** - Create utility functions
- **config** - Create configuration files
- **page** - Create a new page component
- **layout** - Create a new layout component

### Examples

#### Creating a Component

```bash
npx thunder-create component
```

This will prompt you for:
- Component name
- Location (feature, shared, or package)
- Component type (basic, form, modal, page)
- Whether to include tests and Storybook stories

#### Creating a Feature

```bash
npx thunder-create feature
```

This will prompt you for:
- Feature name
- Target app (thunder-develop or thunder-gate)
- Feature type (admin, gate, integration, basic)
- Which directories to include (api, components, pages, etc.)
- Whether to include tests and routing

## Templates

The tool uses Handlebars templates located in `src/templates/` to generate files. Templates include:

- **Thunder naming conventions** - PascalCase for components, kebab-case for directories
- **TypeScript interfaces** - Comprehensive type definitions
- **TanStack Query integration** - API hooks with proper query keys
- **i18n support** - Translation keys and useTranslation hook
- **Testing setup** - Vitest test files with proper mocks
- **@wso2/oxygen-ui integration** - Consistent UI component usage

## Features

- ✅ **Interactive CLI** with beautiful prompts using `@clack/prompts`
- ✅ **Thunder workspace detection** - Only works in valid Thunder workspaces  
- ✅ **File templates** with Handlebars for dynamic content generation
- ✅ **Naming validation** - Ensures proper naming conventions
- ✅ **Duplicate detection** - Prevents overwriting existing files
- ✅ **Auto-imports** - Generates proper import/export statements
- ✅ **Test generation** - Creates test files with proper setup
- ✅ **TypeScript support** - Full TypeScript integration

## Development

```bash
# Install dependencies
pnpm install

# Build the package
pnpm build

# Run tests
pnpm test

# Type checking
pnpm typecheck
```

## Architecture

- **CLI Entry Point**: `src/cli.ts` - Main command interface
- **Commands**: `src/commands/` - Individual command implementations
- **Templates**: `src/templates/` - Handlebars templates for file generation
- **Utils**: `src/utils/` - Workspace detection, file operations, validation

## Contributing

Follow the Thunder project's contribution guidelines. Ensure all new templates follow established patterns and include proper tests.
