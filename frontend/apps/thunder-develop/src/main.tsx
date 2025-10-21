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

import * as ReactDOM from 'react-dom/client';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {BrowserRouter, Route, Routes} from 'react-router';
import {theme} from '@thunder/ui';
import CssBaseline from '@mui/material/CssBaseline';
import {StrictMode} from 'react';
import App from './Dashboard';
import UsersListPage from './features/users/pages/UsersListPage';
import HomePage from './features/home/pages/HomePage';
import CreateUserPage from './features/users/pages/CreateUserPage';
import ViewUserPage from './features/users/pages/ViewUserPage';
import UserTypesListPage from './features/user-types/pages/UserTypesListPage';
import CreateUserTypePage from './features/user-types/pages/CreateUserTypePage';
import ViewUserTypePage from './features/user-types/pages/ViewUserTypePage';
import IntegrationsPage from './features/integrations/pages/IntegrationsPage';
import ApplicationsPage from './features/applications/pages/ApplicationsPage';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <StyledEngineProvider injectFirst>
      <ThemeProvider theme={theme} defaultMode="light">
        <CssBaseline enableColorScheme />
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<App />}>
              <Route path="/" element={<HomePage />} />
              <Route path="users" element={<UsersListPage />} />
              <Route path="users/create" element={<CreateUserPage />} />
              <Route path="users/:userId" element={<ViewUserPage />} />
              <Route path="user-types" element={<UserTypesListPage />} />
              <Route path="user-types/create" element={<CreateUserTypePage />} />
              <Route path="user-types/:id" element={<ViewUserTypePage />} />
              <Route path="integrations" element={<IntegrationsPage />} />
              <Route path="applications" element={<ApplicationsPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </ThemeProvider>
    </StyledEngineProvider>
  </StrictMode>,
);
