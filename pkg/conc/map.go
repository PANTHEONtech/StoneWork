// SPDX-License-Identifier: Apache-2.0

// Copyright 2023 PANTHEON.tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conc

import "sync"

// Map is a very simple concurrent map implementation.
type Map[K comparable, V any] struct {
	mu   *sync.RWMutex
	data map[K]V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{mu: &sync.RWMutex{}, data: make(map[K]V)}
}

func (m Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

func (m Map[K, V]) Set(key K, val V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

func (m Map[K, V]) Del(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Iter returns channel filled with snapshot of the data inside of map.
// The returned channel can be used in a for-range loop.
// Order of map items coming from the channel is not specified.
func (m Map[K, V]) Iter() <-chan KVPair[K, V] {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ch := make(chan KVPair[K, V], len(m.data))
	for k, v := range m.data {
		ch <- KVPair[K, V]{k, v}
	}
	close(ch)
	return ch
}

// KVPair represents a single map key-value pair.
type KVPair[K comparable, V any] struct {
	Key K
	Val V
}
