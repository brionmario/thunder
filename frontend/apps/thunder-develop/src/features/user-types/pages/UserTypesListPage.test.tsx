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

import {describe, it, expect, vi, beforeEach} from 'vitest';
import {screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import render from '@/test/test-utils';
import UserTypesListPage from './UserTypesListPage';

const mockNavigate = vi.fn();
const mockReload = vi.fn();

// Mock react-router
vi.mock('react-router', async () => {
  const actual = await vi.importActual<typeof import('react-router')>('react-router');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Mock the UserTypesList component
vi.mock('../components/UserTypesList', () => ({
  default: () => <div data-testid="user-types-list">User Types List Component</div>,
}));

describe('UserTypesListPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    Object.defineProperty(window, 'location', {
      value: {reload: mockReload},
      writable: true,
    });
  });

  it('renders page title', () => {
    render(<UserTypesListPage />);

    expect(screen.getByText('User Types')).toBeInTheDocument();
  });

  it('renders page description', () => {
    render(<UserTypesListPage />);

    expect(screen.getByText('Define a new user type schema for your organization')).toBeInTheDocument();
  });

  it('renders create user type button', () => {
    render(<UserTypesListPage />);

    const createButton = screen.getByRole('button', {name: /create user type/i});
    expect(createButton).toBeInTheDocument();
  });

  it('renders UserTypesList component', () => {
    render(<UserTypesListPage />);

    expect(screen.getByTestId('user-types-list')).toBeInTheDocument();
  });

  it('navigates to create page when create button is clicked', async () => {
    const user = userEvent.setup();
    render(<UserTypesListPage />);

    const createButton = screen.getByRole('button', {name: /create user type/i});
    await user.click(createButton);

    expect(mockNavigate).toHaveBeenCalledWith('/user-types/create');
  });
});
