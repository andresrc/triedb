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
type node interface {
	// Node value.
	Value() Value

	// Obtains a value
	Get(key []byte, pos int) Value
}

// The empty node
var (
	empty = emptyNode{}
)

// Returns the empty node
func Empty() node {
	return &empty
}

// Empty node definition
type emptyNode struct {
}

func (*emptyNode) Value() Value {
	return nil
}

func (*emptyNode) Get(key []byte, pos int) Value {
	return nil
}

// Leaf node, only contains a value
type leafNode struct {
	value Value
}

func (n *leafNode) Value() Value {
	return n.value
}

func (n *leafNode) Get(key []byte, pos int) Value {
	if pos == len(key)-1 {
		return n.value
	}
	return nil
}

// Full node. Worst case
type fullNode struct {
	nodes [256]node
	value Value
}

func (n *fullNode) Value() Value {
	return n.value
}
