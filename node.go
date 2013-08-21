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

package triedb

// Persistent trie node interface
type node struct {
	key      byte
	value    Value
	next     *node
	children *node
}

func (n *node) get(key []byte, pos int) Value {
	var k = key[pos]
	// Next sibling
	if k > n.key {
		if n.next != nil {
			return n.next.get(key, pos)
		}
		return nil
	}
	// Current node
	if k == n.key {
		var last = len(key) - 1
		if pos < last {
			if n.children != nil {
				return n.children.get(key, pos+1)
			}
			return nil
		}
		if pos == last {
			return n.value
		}
		panic("Position exceeded")
	}
	// Previous sibling, should not happen
	panic("Previous sibling")
}

func (n *node) remove(key []byte, pos int) *node {
	return nil // todo
}
