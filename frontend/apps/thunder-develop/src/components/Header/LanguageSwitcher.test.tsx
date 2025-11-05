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

import {render, screen, fireEvent, waitFor} from '@testing-library/react';
import {describe, it, expect, vi, beforeEach} from 'vitest';
import LanguageSwitcher from './LanguageSwitcher';

// Mock the local useLanguage hook
vi.mock('@/hooks/useLanguage', () => ({
  useLanguage: vi.fn(() => ({
    currentLanguage: 'en-US',
    availableLanguages: [
      {code: 'en-US', name: 'English (United States)', nativeName: 'English (US)', direction: 'ltr'},
    ],
    setLanguage: vi.fn(),
  })),
}));

describe('LanguageSwitcher', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render the language switcher button', () => {
    render(<LanguageSwitcher />);

    const button = screen.getByRole('button', {name: /change language/i});
    expect(button).toBeInTheDocument();
  });

  it('should open language menu when button is clicked', async () => {
    render(<LanguageSwitcher />);

    const button = screen.getByRole('button', {name: /change language/i});
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText('English (US)')).toBeInTheDocument();
    });
  });

  it('should call setLanguage when selecting a language', async () => {
    const mockSetLanguage = vi.fn().mockResolvedValue(undefined);
    const {useLanguage} = await import('@/hooks/useLanguage');
    vi.mocked(useLanguage).mockReturnValue({
      currentLanguage: 'en-US',
      availableLanguages: [
        {code: 'en-US', name: 'English (United States)', nativeName: 'English (US)', direction: 'ltr'},
      ],
      setLanguage: mockSetLanguage,
    });

    render(<LanguageSwitcher />);

    const button = screen.getByRole('button', {name: /change language/i});
    fireEvent.click(button);

    await waitFor(() => {
      const englishOption = screen.getByText('English (US)');
      fireEvent.click(englishOption);
    });

    expect(mockSetLanguage).toHaveBeenCalledWith('en-US');
  });

  it('should close the menu after selecting a language', async () => {
    const mockSetLanguage = vi.fn().mockResolvedValue(undefined);
    const {useLanguage} = await import('@/hooks/useLanguage');
    vi.mocked(useLanguage).mockReturnValue({
      currentLanguage: 'en-US',
      availableLanguages: [
        {code: 'en-US', name: 'English (United States)', nativeName: 'English (US)', direction: 'ltr'},
      ],
      setLanguage: mockSetLanguage,
    });

    render(<LanguageSwitcher />);

    const button = screen.getByRole('button', {name: /change language/i});
    fireEvent.click(button);

    await waitFor(() => {
      const englishOption = screen.getByText('English (US)');
      fireEvent.click(englishOption);
    });

    await waitFor(() => {
      expect(screen.queryByRole('menu')).not.toBeInTheDocument();
    });
  });
});
