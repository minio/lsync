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
	"fmt"
	"testing"
	"time"

	. "github.com/minio/lsync"
)

func testSimpleWriteLock(t *testing.T, duration time.Duration) (locked bool) {

	lrwm := NewLRWMutex("resource")

	if !lrwm.GetRLock(time.Second) {
		panic("Failed to acquire read lock")
	}
	// fmt.Println("1st read lock acquired, waiting...")

	if !lrwm.GetRLock(time.Second) {
		panic("Failed to acquire read lock")
	}
	// fmt.Println("2nd read lock acquired, waiting...")

	go func() {
		time.Sleep(2*time.Second)
		lrwm.RUnlock()
		// fmt.Println("1st read lock released, waiting...")
	}()

	go func() {
		time.Sleep(3*time.Second)
		lrwm.RUnlock()
		// fmt.Println("2nd read lock released, waiting...")
	}()

	// fmt.Println("Trying to acquire write lock, waiting...")
	locked = lrwm.GetLock(duration)
	if locked {
		// fmt.Println("Write lock acquired, waiting...")
		time.Sleep(1*time.Second)

		lrwm.Unlock()
	} else {
		// fmt.Println("Write lock failed due to timeout")
	}
	return
}

func TestSimpleWriteLockAcquired(t *testing.T) {
	locked := testSimpleWriteLock(t, 5*time.Second)

	expected := true
	if locked != expected {
		t.Errorf("TestSimpleWriteLockAcquired(): \nexpected %#v\ngot      %#v", expected, locked)
	}
}

func TestSimpleWriteLockTimedOut(t *testing.T) {
	locked := testSimpleWriteLock(t, time.Second)

	expected := false
	if locked != expected {
		t.Errorf("TestSimpleWriteLockTimedOut(): \nexpected %#v\ngot      %#v", expected, locked)
	}
}


func testDualWriteLock(t *testing.T, duration time.Duration) (locked bool) {

	lrwm := NewLRWMutex("resource")

	// fmt.Println("Getting initial write lock")
	if !lrwm.GetLock(time.Second) {
		panic("Failed to acquire initial write lock")
	}

	go func() {
		time.Sleep(2*time.Second)
		lrwm.Unlock()
		// fmt.Println("Initial write lock released, waiting...")
	}()

	// fmt.Println("Trying to acquire 2nd write lock, waiting...")
	locked = lrwm.GetLock(duration)
	if locked {
		// fmt.Println("2nd write lock acquired, waiting...")
		time.Sleep(time.Second)

		lrwm.Unlock()
	} else {
		// fmt.Println("2nd write lock failed due to timeout")
	}
	return
}

func TestDualWriteLockAcquired(t *testing.T) {
	locked := testDualWriteLock(t, 3*time.Second)

	expected := true
	if locked != expected {
		t.Errorf("TestDualWriteLockAcquired(): \nexpected %#v\ngot      %#v", expected, locked)
	}

}

func TestDualWriteLockTimedOut(t *testing.T) {
	locked := testDualWriteLock(t, time.Second)

	expected := false
	if locked != expected {
		t.Errorf("TestDualWriteLockTimedOut(): \nexpected %#v\ngot      %#v", expected, locked)
	}

}

	}
}
