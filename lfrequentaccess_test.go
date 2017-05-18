/*
 * Minio Cloud Storage, (C) 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// GOMAXPROCS=10 go test

package lsync_test

import (
	"testing"

	. "github.com/minio/lsync"
	"sync"
)

type Map map[string]string

const key = "test"
const val = "value"

func TestLFrequentAccess(t *testing.T) {
	// Create new LFrequentAccess for type Map
	freqaccess := NewLFrequentAccess(make(Map))

	cur := freqaccess.LockBeforeSet().(Map) // Lock in order to update
	mp := make(Map)                         // Create new Map
	for k, v := range cur {                 // Copy over old contents
		mp[k] = v
	}
	mp[key] = val                      // Add new value
	freqaccess.SetNewCopyAndUnlock(mp) // Exchange old version of map with new version

	mp2 := freqaccess.ReadOnlyAccess().(Map) // Get read only access to Map
	if mp2[key] != val {
		t.Errorf("TestLFrequentAccess(): \nexpected %#v\ngot      %#v", val, mp2[key])
	}
}

func BenchmarkLFrequentAccessMap(b *testing.B) {
	m := NewLFrequentAccess(make(Map))

	cur := m.LockBeforeSet().(Map)
	mp := make(Map)
	for k, v := range cur {
		mp[k] = v
	}
	mp[key] = val
	m.SetNewCopyAndUnlock(mp)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mp := m.ReadOnlyAccess().(Map)
			if mp[key] != val {
				panic("Expected key value")
			}
		}
	})
}

func BenchmarkLFrequentAccessMapRegularMutex(b *testing.B) {
	type MapMutex struct {
		mu sync.Mutex
		mp Map
	}

	m := MapMutex{}

	m.mp = make(Map)
	m.mu.Lock()
	m.mp[key] = val
	m.mu.Unlock()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.mu.Lock()
			mp := m.mp
			if mp[key] != val {
				panic("Expected key value")
			}
			m.mu.Unlock()
		}
	})
}

type Slice []string

func BenchmarkLFrequentAccessSlice(b *testing.B) {

	m := NewLFrequentAccess(make(Slice, 0))

	cur := m.LockBeforeSet().(Slice)
	slc := make(Slice, len(cur))
	for i, v := range cur {
		slc[i] = v
	}
	slc = append(slc, val)
	m.SetNewCopyAndUnlock(slc)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			slc := m.ReadOnlyAccess().(Slice)
			if slc[0] != val {
				panic("Expected key value")
			}
		}
	})
}

func BenchmarkLFrequentAccessSliceRegularMutex(b *testing.B) {
	type SliceMutex struct {
		mu  sync.Mutex
		slc Slice
	}

	s := SliceMutex{}

	s.slc = make(Slice, 0)
	s.mu.Lock()
	s.slc = append(s.slc, val)
	s.mu.Unlock()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.mu.Lock()
			slc := s.slc
			if slc[0] != val {
				panic("Expected key value")
			}
			s.mu.Unlock()
		}
	})
}
