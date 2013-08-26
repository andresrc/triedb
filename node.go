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

// A node represents a node in the persistent trie used as an index.
type node struct {
	key      byte
	depth    int16
	count    int32
	value    Value
	next     *node
	children *node
}

// createNode creates a new node. Every node must be created using this function
// (or other based on this) in order to maintain the invariants.
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
		c += children.count
	} else {
		d = 1
	}
	return &node{key, d, c, value, next, children}
}

// clear returns nil if the node is empty, or the current node if not
func (n *node) clear() *node {
	if n.value == nil && n.next == nil && n.children == nil {
		return nil
	}
	return n
}

// createLeaf creates a new node with no children and no siblings.
func createLeaf(key byte, value Value) *node {
	return createNode(key, value, nil, nil)
}

// getKey returns the byte at the current position and the last index for the provided key and position
func getKey(key []byte, pos int) (k byte, last int) {
	last = len(key) - 1
	if pos >= last {
		panic("Position exceeded")
	}
	k = key[pos]
	return
}

// get returns the value for a key slice, if any.
func (n *node) get(key []byte, pos int) Value {
	k, last := getKey(key, pos)
	// Refer to next sibling
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
	// It was for a previous sibling
	return nil
}

// setNext creates a new node with the provided next sibling pointer if it is different from the current one.
func (n *node) setNext(next *node) *node {
	if n.next == next {
		return n
	}
	return createNode(n.key, n.value, next, n.children)
}

// setChildren creates a new node with the provided first sibling pointer if it is different from the current one.
func (n *node) setChildren(children *node) *node {
	if n.children == children {
		return n
	}
	return createNode(n.key, n.value, n.next, children)
}

// createSibling creates a new sibling node.
func (n *node) createSibling(key []byte, pos int, value Value) *node {
	k, last := getKey(key, pos)
	prev := k < n.key
	// If no value is provided it's done.
	if value == nil {
		if prev {
			return n
		} else {
			return nil
		}
	}
	// Calculate next pointer
	var next *node
	if prev {
		next = n
	} else {
		next = nil
	}
	// If we are not at the end of the key, create, complete and return an intermediate node
	if pos < last {
		return createNode(k, nil, next, nil).set(key, pos, value)
	}
	// End of the key, create and return a sibling node.
	return createNode(k, value, next, nil)
}

// set creates a new node with the provided value for the provided key slice if it is different from the current one.
func (n *node) set(key []byte, pos int, value Value) *node {
	k, last := getKey(key, pos)
	// It is for a previous sibling
	if k < n.key {
		// We create a previous sibling node if needed
		return n.createSibling(key, pos, value)
	}
	// Next sibling
	if k > n.key {
		var next *node
		if n.next != nil {
			next = n.next.set(key, pos, value)
		} else {
			next = n.createSibling(key, pos, value)
		}
		return n.setNext(next).clear()
	}
	// Current node
	if k == n.key {
		if pos < last {
			var children *node
			if n.children != nil {
				children = n.children
			} else {
				if value == nil {
					return n // Unmodified
				}
				children = createNode(key[pos+1], nil, nil, nil)
			}
			return n.setChildren(children.set(key, pos+1, value)).clear()
		} else { // this is the node to change
			if n.value == value {
				return n // Unmodified
			} else {
				return createNode(n.key, value, n.next, n.children).clear()
			}
		}
	}
	panic("Unreachable")
}

// remove clears the value for a key, removing the nodes that are not needed any more.
func (n *node) remove(key []byte, pos int) *node {
	return n.set(key, pos, nil)
}
