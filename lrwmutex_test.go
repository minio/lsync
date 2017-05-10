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

func TestSimpleWriteLock(t *testing.T) {

	lrwm := NewLRWMutex("resource")

	if !lrwm.GetRLock(1000 * time.Millisecond) {
		panic("Failed to acquire read lock")
	}
	fmt.Println("1st read lock acquired, waiting...")

	if !lrwm.GetRLock(1000 * time.Millisecond) {
		panic("Failed to acquire read lock")
	}
	fmt.Println("2nd read lock acquired, waiting...")

	go func() {
		time.Sleep(1000 * time.Millisecond)
		lrwm.RUnlock()
		fmt.Println("1st read lock released, waiting...")
	}()

	 go func() {
		time.Sleep(2000 * time.Millisecond)
		lrwm.RUnlock()
		fmt.Println("2nd read lock released, waiting...")
	}()

	fmt.Println("Trying to acquire write lock, waiting...")
	if lrwm.GetLock(1000 * time.Millisecond) {
		fmt.Println("Write lock acquired, waiting...")
		time.Sleep(2500 * time.Millisecond)

		lrwm.Unlock()
	} else {
		fmt.Println("Write lock failed due to timeout")
	}
}

func TestDualWriteLock(t *testing.T) {

	lrwm := NewLRWMutex("resource")

	fmt.Println("Getting initial write lock")
	if !lrwm.GetLock(1500 * time.Millisecond) {
		panic("Failed to acquire initial write lock")
	}

	go func() {
		time.Sleep(2500 * time.Millisecond)
		lrwm.Unlock()
		fmt.Println("Initial write lock released, waiting...")
	}()

	fmt.Println("Trying to acquire 2nd write lock, waiting...")
	if lrwm.GetLock(1000 * time.Millisecond) {
		fmt.Println("2nd write lock acquired, waiting...")
		time.Sleep(2500 * time.Millisecond)

		lrwm.Unlock()
	} else {
		fmt.Println("2nd write lock failed due to timeout")
	}
}
