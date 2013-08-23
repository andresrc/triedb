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
	depth    int16
	count    int32
	value    Value
	next     *node
	children *node
}

// Creates a new node
func createNode(key byte, value Value, next *node, children *node) *node {
	var d int16
	var c int32
	if value == nil {
		c = 0
	} else {
		c = 1
	}
	if next != nil {
		c += next.count
	}
	if children != nil {
		d = 1 + children.depth
	} else {
		d = 1
	}
	return &node{key, d, c, value, next, children}
}

// Creates a leaf node
func createLeaf(key byte, value Value) *node {
	return createNode(key, value, nil, nil)
}

// Gets the value for a key.
func (n *node) get(key []byte, pos int) Value {
	var last = len(key) - 1
	if pos >= last {
		panic("Position exceeded")
	}
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
		if pos < last {
			if n.children != nil {
				return n.children.get(key, pos+1)
			}
			return nil
		}
		if pos == last {
			return n.value
		}
	}
	// Previous sibling
	return nil
}

// Changes the next pointer of a node
func (n *node) setNext(next *node) *node {
	if n.next == next {
		return n
	}
	return createNode(n.key, n.value, next, n.children)
}

// Changes the children pointer of a node
func (n *node) setChildren(children *node) *node {
	if n.children == children {
		return n
	}
	return createNode(n.key, n.value, n.next, children)
}

// Sets the value of a key.
func (n *node) set(key []byte, pos int, value Value) *node {
	var last = len(key) - 1
	if pos >= last {
		panic("Position exceeded")
	}
	var k = key[pos]
	// Previous sibling
	if k < n.key {
		return n // Not found, unmodified
	}
	var ret *node // return value
	// Next sibling
	if k > n.key {
		if n.next != nil {
			ret = n.setNext(n.next.set(key, pos, value))
		} else {
			return n // Not found, unmodified
		}
	}
	// Current node
	if k == n.key {
		if pos < last {
			if n.children != nil {
				ret = n.setChildren(n.children.set(key, pos+1, value))
			} else {
				return n // Not found, unmodified
			}
		} else { // this is the node to change
			if n.value == value {
				return n // Unmodified
			} else {
				ret = createNode(n.key, value, n.next, n.children)
			}
		}
	}
	// Check if the node has to disappear
	if ret.value == nil && ret.next == nil && ret.children == nil {
		return nil
	}
	return ret
}

// Removes a key
func (n *node) remove(key []byte, pos int) *node {
	return n.set(key, pos, nil)
}
