# Contributing to WSO2 Thunder ⚡

Thank you for your interest in contributing to WSO2 Thunder! This guide will help you set up your development environment and understand the contribution process.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
  - [Required Tools (Core Development)](#required-tools-core-development)
  - [Required Tools (Frontend Development)](#required-tools-frontend-development)
- [Development Setup](#development-setup)

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### Required Tools (Core Development)

- **[Git](https://git-scm.com/downloads)** - Version control system
- **[Go](https://golang.org/doc/install)** - Version 1.25 or higher
- **[Node.js](https://nodejs.org/en/download/)** - Version 22 or higher (`LTS` is recommended ✅)
- **[pnpm](https://pnpm.io/installation)** - Version 9 or higher (`LTS` is recommended ✅)

### Required Tools (Frontend Development)

- **[ESLint VSCode Extension](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)** - For linting support in VSCode
- **[Prettier VSCode Extension](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode)** - For code formatting in VSCode

## Development Setup

1. Fork the Repository

Start by forking the WSO2 Thunder repository to your GitHub account.

2. Clone the Repository

Clone your forked repository to your local machine.

```bash
git clone https://github.com/<your-username>/thunder.git
cd thunder
```

3. Configure the `deployment.yaml`.

Open up `backend/cmd/server/repository/conf/deployment.yaml` in your favorite text editor and make the following changes:

- Under the `cors` section, add the following allowed origins:

```yaml
  cors:
    allowed_origins:
        - "https://localhost:5190"
        - "https://localhost:5191"
```

> [!TIP]
> This configuration allows the Thunder Gate and Thunder Develop applications to communicate with the backend server during development.

- Under the `gate_client` section, add the following configuration:

```yaml
  gate_client:
    port: 5190
    scheme: "https"
    login_path: "/signin"
```

> [!TIP]
> This configuration points the backend server to the local Thunder Gate application for authentication.

4. Run the Project

Start the Thunder server & apps in development mode.

```bash
make run
```

> [!NOTE]
> This command will automatically set up your complete development environment:
>
> **What gets created:**
>
> - `Develop` application with redirect URI pointing to the local Thunder Develop application
> - Default admin user with credentials: `admin` / `admin`
>
> **Services that start:**
>
> - **Backend server**: [https://localhost:8090](https://localhost:8090)
> - **Thunder Gate** (Login/Register): [https://localhost:5190/signin](https://localhost:5190/signin)
> - **Thunder Develop** (Admin Console): [https://localhost:5191/develop](https://localhost:5191/develop)

💡 Alternatively, if you would rather run the backend and frontend servers manually, run the following commands in separate terminal windows:

**Backend Server:**

```bash
make run_backend
```

**Frontend Server:**

```bash
make run_frontend
```

If you would rather run the frontend applications manually without `make` commands, follow the instructions in the **📖 Development Setup (Frontend) - Manual Mode** to setup the frontend.

<details>
<summary><strong>📖 Development Setup (Frontend) - Manual Mode</strong></summary>

## Installing Dependencies

1. Navigate to the Thunder frontend directory.

```bash
cd frontend
```

2. Install the dependencies using `pnpm`.

```bash
pnpm install
```

### Building the Project

Execute the build command to compile the project. This will build all the necessary packages and applications.

```bash
pnpm build
```

### Seed Data (Optional)

💡 **NOTE**: This step is only necessary if you are running the backend server manually and have not yet set up the initial data.

If you have not already created the `Develop` application and the default admin user, you can do so by running the following command:

```bash
cd backend/scripts
sh ./setup_initial_data.sh $BACKEND_PORT --develop-redirect-uris "https://localhost:$DEVELOP_APP_DEFAULT_PORT/develop"
```

### Setting up the Thunder Gate Application

1. Point the `gate_client` in `deployment.yaml` to the local Thunder Gate application.

- If you are running `make run`, change the `gate_client` section in `backend/cmd/server/repository/conf/deployment.yaml`:
- If you are running the backend server manually, change the `gate_client` section in `<THUNDER_HOME>/repository/conf/deployment.yaml`:

```yaml
gate_client:
  port: 5190
  scheme: "https"
  login_path: "/signin"
```

2. Add the local development origin of the Thunder Gate application (https://localhost:5190) to the CORS allowed origins in `<THUNDER_HOME>/repository/conf/deployment.yaml`.

```yaml
cors:
  allowed_origins:
    - "https://localhost:5190"
```

3. Run the Thunder Gate application.

```bash
cd frontend
pnpm --filter @thunder/gate dev
```

### Setting up the Thunder Develop Application

💡 **IMPORTANT**: This section assumes that you have already created the `Develop` application using the initial data setup script. If not, please refer to the [Seed Data](#seed-data-optional) section above.

1. First, retrieve the application ID of the **Develop** application from the Thunder server. This will be the application with the `client_id` **DEVELOP**.

```bash
curl -k -X GET "https://localhost:8090/applications"
```

2. Then, get the current **Develop** application configuration:

```bash
curl -k -X GET "https://localhost:8090/applications/<develop-application-id>"
```

💡 **NOTE**: Replace `<develop-application-id>` with the actual application ID (e.g., `6100bc91-ba99-4ce9-87dd-6d4d80178c38`) obtained from the previous step. The `-k` flag allows curl to work with self-signed SSL certificates in development.

2. Copy the response from `step 2` and update the `redirect_uris` in the JSON object to include the local development URL (ex: https://localhost:5191/develop). Locate the `inbound_auth_config > config` section and modify the `redirect_uris` array:

```diff
"redirect_uris": [
  "https://localhost:8090/develop",
+  "https://localhost:5191/develop"
]
```

3. Update the **Develop** application with the modified configuration by passing the updated JSON directly:

```bash
curl -k -X PUT "https://localhost:8090/applications/<develop-application-id>" \
  -H "Content-Type: application/json" \
  -d '<paste-the-modified-json-here>'
```

4. Add the local development origin of the thunder-develop application (https://localhost:5191) to the CORS allowed origins in `thunder-home/repository/conf/deployment.yaml`.

```yaml
cors:
  allowed_origins:
    - "https://localhost:5191"
```

5. Run the Thunder Develop application.

```bash
pnpm --filter @thunder/develop dev
```

This will run the Thunder Develop application on `https://localhost:5191/develop`.

</details>
