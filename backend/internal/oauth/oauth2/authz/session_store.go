/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package authz

import (
	"sync"
	"time"

	"github.com/asgardeo/thunder/internal/oauth/oauth2/model"
	"github.com/asgardeo/thunder/internal/system/utils"
)

// SessionData holds OAuth session information including parameters and authentication time.
type SessionData struct {
	OAuthParameters model.OAuthParameters
	AuthTime        time.Time
}

// sessionDataStoreInterface defines the interface for session data storage.
type sessionDataStoreInterface interface {
	AddSession(value SessionData) string
	GetSession(key string) (bool, SessionData)
	ClearSession(key string)
	ClearSessionStore()
}

// sessionStoreEntry represents an entry in the session data store.
type sessionStoreEntry struct {
	sessionData SessionData
	expiryTime  time.Time
}

// sessionDataStore provides the session data store functionality.
type sessionDataStore struct {
	sessionStore   map[string]sessionStoreEntry
	validityPeriod time.Duration
	mu             sync.RWMutex
}

// newSessionDataStore creates a new instance of sessionDataStore with injected dependencies.
func newSessionDataStore() sessionDataStoreInterface {
	return &sessionDataStore{
		sessionStore:   make(map[string]sessionStoreEntry),
		validityPeriod: 10 * time.Minute, // Set a default validity period.
	}
}

// AddSession adds a session data entry to the session store.
func (sds *sessionDataStore) AddSession(value SessionData) string {
	sds.mu.Lock()
	defer sds.mu.Unlock()

	key := utils.GenerateUUID()
	sds.sessionStore[key] = sessionStoreEntry{
		sessionData: value,
		expiryTime:  time.Now().Add(sds.validityPeriod),
	}
	return key
}

// GetSession retrieves a session data entry from the session store.
func (sdc *sessionDataStore) GetSession(key string) (bool, SessionData) {
	if key == "" {
		return false, SessionData{}
	}

	sdc.mu.RLock()
	entry, exists := sdc.sessionStore[key]
	sdc.mu.RUnlock()

	if exists {
		if time.Now().Before(entry.expiryTime) {
			return true, entry.sessionData
		} else {
			// Remove the expired entry.
			sdc.mu.Lock()
			delete(sdc.sessionStore, key)
			sdc.mu.Unlock()
		}
	}

	return false, SessionData{}
}

// ClearSession removes a specific session data entry from the session store.
func (sdc *sessionDataStore) ClearSession(key string) {
	if key == "" {
		return
	}

	sdc.mu.Lock()
	defer sdc.mu.Unlock()
	delete(sdc.sessionStore, key)
}

// ClearSessionStore removes all session data entries from the session store.
func (sdc *sessionDataStore) ClearSessionStore() {
	sdc.mu.Lock()
	defer sdc.mu.Unlock()

	sdc.sessionStore = make(map[string]sessionStoreEntry)
}
