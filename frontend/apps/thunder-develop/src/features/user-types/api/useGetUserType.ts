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

import { useEffect, useRef, useTransition } from 'react';
import { useState } from 'react';

import type { ApiError, ApiUserSchema } from '../types/user-types';

/**
 * Custom hook to fetch a single user schema (user type) by ID
 * Includes double-fetch prevention for React Strict Mode
 * @param id - The user schema ID to fetch
 * @returns Object containing data, loading state, error, and refetch function
 */
export default function useGetUserType(id?: string) {
  const [data, setData] = useState<ApiUserSchema | null>(null);
  const [error, setError] = useState<ApiError | null>(null);
  const [isPending, startFetchTransition] = useTransition();

  // Refs to prevent double-fetch in React Strict Mode
  const hasFetchedRef = useRef(false);
  const lastIdRef = useRef<string | undefined>(undefined);

  const fetchUserType = async (userTypeId: string) => {
    startFetchTransition(async () => {
      setError(null);

      try {
        const response = await fetch(`https://localhost:8090/user-schemas/${userTypeId}`);

        if (!response.ok) {
          let errorData: ApiError;
          try {
            errorData = (await response.json()) as ApiError;
          } catch {
            errorData = {
              code: 'FETCH_USER_TYPE_ERROR',
              message: `HTTP error! status: ${response.status}`,
              description: await response.text(),
            };
          }
          setError(errorData);
          setData(null);
          return;
        }

        const jsonData = (await response.json()) as ApiUserSchema;
        setData(jsonData);
        setError(null);
      } catch (err) {
        const apiError: ApiError = {
          code: 'FETCH_USER_TYPE_ERROR',
          message: err instanceof Error ? err.message : 'An unknown error occurred',
          description: 'Failed to fetch user type',
        };
        setError(apiError);
        setData(null);
      }
    });
  };

  const refetch = (newId?: string) => {
    const idToFetch = newId ?? id;
    if (!idToFetch) return;

    // Reset the hasFetched flag when explicitly refetching
    hasFetchedRef.current = false;
    lastIdRef.current = idToFetch;
    fetchUserType(idToFetch);
  };

  useEffect(() => {
    if (!id) {
      setData(null);
      setError(null);
      return;
    }

    // Check if we've already fetched for this ID
    if (hasFetchedRef.current && lastIdRef.current === id) {
      return;
    }

    // Mark as fetched and store the ID
    hasFetchedRef.current = true;
    lastIdRef.current = id;

    fetchUserType(id);
  }, [id]);

  return {
    data,
    loading: isPending,
    error,
    refetch,
  };
}
