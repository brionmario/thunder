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

import { useState, useTransition } from 'react';

import type { ApiError } from '../types/user-types';

/**
 * Custom hook to delete a user schema (user type)
 * @returns Object containing deleteUserType function, loading state, error, and reset function
 */
export default function useDeleteUserType() {
  const [error, setError] = useState<ApiError | null>(null);
  const [isPending, startDeleteTransition] = useTransition();

  const deleteUserType = async (userTypeId: string): Promise<boolean> => {
    let success = false;

    await new Promise<void>((resolve) => {
      startDeleteTransition(async () => {
        setError(null);

        try {
          const response = await fetch(`https://localhost:8090/user-schemas/${userTypeId}`, {
            method: 'DELETE',
          });

          if (!response.ok) {
            let errorData: ApiError;
            try {
              errorData = (await response.json()) as ApiError;
            } catch {
              errorData = {
                code: 'DELETE_USER_TYPE_ERROR',
                message: `HTTP error! status: ${response.status}`,
                description: await response.text(),
              };
            }
            setError(errorData);
            success = false;
          } else {
            setError(null);
            success = true;
          }
        } catch (err) {
          const apiError: ApiError = {
            code: 'DELETE_USER_TYPE_ERROR',
            message: err instanceof Error ? err.message : 'An unknown error occurred',
            description: 'Failed to delete user type',
          };
          setError(apiError);
          success = false;
        }

        resolve();
      });
    });

    return success;
  };

  const reset = () => {
    setError(null);
  };

  return {
    deleteUserType,
    loading: isPending,
    error,
    reset,
  };
}
