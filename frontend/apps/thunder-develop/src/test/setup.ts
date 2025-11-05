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

/* eslint-disable max-classes-per-file */

import '@testing-library/jest-dom';
import {cleanup} from '@testing-library/react';
import {afterEach, beforeAll} from 'vitest';
import i18n from 'i18next';
import {initReactI18next} from 'react-i18next';
import enUS from '@thunder/i18n/locales/en-US';

// Initialize i18n for tests
beforeAll(async () => {
  await i18n.use(initReactI18next).init({
    resources: {
      'en-US': {
        common: enUS.common,
        navigation: enUS.navigation,
        users: enUS.users,
        userTypes: enUS.userTypes,
        integrations: enUS.integrations,
        applications: enUS.applications,
        dashboard: enUS.dashboard,
        auth: enUS.auth,
        mfa: enUS.mfa,
        social: enUS.social,
        consent: enUS.consent,
        errors: enUS.errors,
      },
    },
    lng: 'en-US',
    fallbackLng: 'en-US',
    defaultNS: 'common',
    interpolation: {
      escapeValue: false,
    },
    // Disable Suspense in tests for faster execution
    react: {
      useSuspense: false,
    },
  });
});

// Cleanup after each test
afterEach(() => {
  cleanup();
});

// Mock HTMLMediaElement methods that don't exist in jsdom
Object.defineProperty(window.HTMLMediaElement.prototype, 'play', {
  configurable: true,
  value: () => Promise.resolve(),
});

Object.defineProperty(window.HTMLMediaElement.prototype, 'pause', {
  configurable: true,
  value: () => {
    // Intentionally empty
  },
});

Object.defineProperty(window.HTMLMediaElement.prototype, 'load', {
  configurable: true,
  value: () => {
    // Intentionally empty
  },
});

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
  readonly root = null;

  readonly rootMargin = '';

  readonly thresholds = [];

  observe() {
    return this;
  }

  disconnect() {
    return this;
  }

  unobserve() {
    return this;
  }
} as unknown as typeof IntersectionObserver;

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
  observe() {
    return this;
  }

  disconnect() {
    return this;
  }

  unobserve() {
    return this;
  }
} as unknown as typeof ResizeObserver;
