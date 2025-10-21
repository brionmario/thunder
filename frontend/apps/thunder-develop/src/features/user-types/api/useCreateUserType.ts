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

import type { ApiError, ApiUserSchema, CreateUserSchemaRequest } from '../types/user-types';

/**
 * Custom hook to create a new user schema (user type)
 * @returns Object containing createUserType function, data, loading state, error, and reset function
 */
export default function useCreateUserType() {
  const [data, setData] = useState<ApiUserSchema | null>(null);
  const [error, setError] = useState<ApiError | null>(null);
  const [isPending, startCreateTransition] = useTransition();

  const createUserType = async (requestData: CreateUserSchemaRequest): Promise<void> => {
    startCreateTransition(async () => {
      setError(null);
      setData(null);

      try {
        const response = await fetch('https://localhost:8090/user-schemas', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(requestData),
        });

        if (!response.ok) {
          let errorData: ApiError;
          try {
            errorData = (await response.json()) as ApiError;
          } catch {
            errorData = {
              code: 'CREATE_USER_TYPE_ERROR',
              message: `HTTP error! status: ${response.status}`,
              description: await response.text(),
            };
          }
          setError(errorData);
          return;
        }

        const jsonData = (await response.json()) as ApiUserSchema;
        setData(jsonData);
        setError(null);
      } catch (err) {
        const apiError: ApiError = {
          code: 'CREATE_USER_TYPE_ERROR',
          message: err instanceof Error ? err.message : 'An unknown error occurred',
          description: 'Failed to create user type',
        };
        setError(apiError);
      }
    });
  };

  const reset = () => {
    setData(null);
    setError(null);
  };

  return {
    createUserType,
    data,
    loading: isPending,
    error,
    reset,
  };
}
