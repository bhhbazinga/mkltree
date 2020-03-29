package mkltree

import (
	"crypto/sha256"
	"hash"
)

// Merkle Tree. It can only be build by using New* method, for example: mkltree.NewMklTree()
// You can choose whether to store blocks in this merkle tree, if yes, it
// will cause more memory. You can only build a tree once.
// If you want a new tree, try to build a new one.
type mklTree struct {
	hasher hash.Hash

	hashes [][][]byte
	blocks [][]byte // optional: store original blocks
}

// NewMklTree return an merkle tree with default hash method sha256
func NewMklTree(blocks [][]byte, storeBlocks bool) *mklTree {
	return NewMklTreeCustomHash(blocks, storeBlocks, sha256.New())
}

// NewMklTreeCustomHash return an empty merkle tree with inputted hash method
func NewMklTreeCustomHash(blocks [][]byte, storeBlocks bool, h hash.Hash) *mklTree {
	m := &mklTree{
		h,
		[][][]byte{},
		[][]byte{},
	}

	// build hash tree // TODO: move to Add method
	leafLevelHashes := [][]byte{}
	for _, b := range blocks {
		if storeBlocks {
			m.blocks = append(m.blocks, b)
		}
		bhash := m.hasher.Sum(b)
		leafLevelHashes = append(leafLevelHashes, bhash)
	}
	m.hashes = append(m.hashes, leafLevelHashes)

	// TODO: move to Build method
	// hash into tree structure, leaf to root, bottom-up
	for i := 0; i < len(m.hashes); i++ {
		hashes := m.hashes[i]
		if len(hashes) > 1 {
			nextLevelHashes := [][]byte{}
			for i := 0; i < len(hashes); i = i + 2 {
				leftI, rightI := i, i+1
				left := hashes[leftI]
				right := []byte{}
				if rightI < len(hashes) {
					right = hashes[rightI]
				}

				m.hasher.Write(left)
				newHash := m.hasher.Sum(right)
				m.hasher.Reset()
				nextLevelHashes = append(nextLevelHashes, newHash)
			}
			m.hashes = append(m.hashes, nextLevelHashes)
		}
	}

	return m
}

//TODO: implement Add and Build methods, to allow incrementally reading big file's block to build a hash tree.

func (m *mklTree) Root() []byte {
	// last level element, which should only contain one, if not empty.
	if (len(m.hashes)) > 0 {
		return m.hashes[len(m.hashes)-1][0]
	}
	return []byte{}
}
