/*
Copyright (C) Andres Rodriguez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package triedb provides the TrieDB runtime plus a byte slice-based
// reference back-end implementation, used for testing pruposes.
package triedb

import "io"

// A value in the datastore
type Value interface {
	Len() int64
	io.WriterTo
}

// A TrieDB datastore
type DB interface {
}

// A TrieDB backend
type Backend interface {
	Load(key []byte, reader io.Reader, current Value) Value
}
