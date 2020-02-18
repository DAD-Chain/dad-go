package merkle

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

func TestMerkleLeaf3(t *testing.T) {
	hasher := TreeHasher{}
	leafs := []Uint256{hasher.hash_leaf([]byte{1}),
		hasher.hash_leaf([]byte{2}),
		hasher.hash_leaf([]byte{3})}
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	if tree.Root() != sha256.Sum256(nil) {
		t.Fatal("root error")
	}
	for i, _ := range leafs {
		tree.Append([]byte{byte(i + 1)})
	}

	hashes := make([]Uint256, 5, 5)
	for i := 0; i < 4; i++ {
		hashes[i], _ = tree.hashStore.GetHash(uint32(i))
	}
	hashes[4] = tree.Root()

	cmp := []Uint256{
		leafs[0],
		leafs[1],
		hasher.hash_children(leafs[0], leafs[1]),
		leafs[2],
		hasher.hash_children(hasher.hash_children(leafs[0], leafs[1]),
			leafs[2]),
	}

	for i := 0; i < 5; i++ {
		if hashes[i] != cmp[i] {
			t.Fatal(fmt.Sprintf("error: %d, expected %x, found %x", i, cmp[i], hashes[i]))
		}
	}

}

func TestMerkle(t *testing.T) {
	hasher := TreeHasher{}
	leafs := []Uint256{hasher.hash_leaf([]byte{1}),
		hasher.hash_leaf([]byte{2}),
		hasher.hash_leaf([]byte{3}),
		hasher.hash_leaf([]byte{4})}

	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	if tree.Root() != sha256.Sum256(nil) {
		t.Fatal("root error")
	}
	for i, _ := range leafs {
		tree.Append([]byte{byte(i + 1)})
	}

	hashes := make([]Uint256, 6, 6)
	for i := 0; i < 6; i++ {
		hashes[i], _ = tree.hashStore.GetHash(uint32(i))
	}
	cmp := []Uint256{
		leafs[0],
		leafs[1],
		hasher.hash_children(leafs[0], leafs[1]),
		leafs[2],
		leafs[3],
		hasher.hash_children(leafs[2], leafs[3]),
		hasher.hash_children(hasher.hash_children(leafs[0], leafs[1]),
			hasher.hash_children(leafs[2], leafs[3])),
	}

	for i := 0; i < 6; i++ {
		if hashes[i] != cmp[i] {
			fmt.Println(hashes)
			fmt.Println(cmp)
			t.Fatal(fmt.Sprintf("error: %d, expected %x, found %x", i, cmp[i], hashes[i]))
		}
	}

}

func TestMerkleHashes(t *testing.T) {
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	for i := 0; i < 100; i++ {
		tree.Append([]byte{byte(i + 1)})
	}

	// 100 == 110 0100
	if len(tree.hashes) != 3 {
		t.Fatal(fmt.Sprintf("error tree hashes size"))
	}

}

// zero based return merkle root of D[0:n]
func TestMerkleRoot(t *testing.T) {
	n := 100
	roots := make([]Uint256, n, n)
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	for i := 0; i < n; i++ {
		tree.Append([]byte{byte(i + 1)})
		roots[i] = tree.Root()
	}

	cmp := make([]Uint256, n, n)
	for i := 0; i < n; i++ {
		cmp[i] = tree.merkleRoot(uint32(i) + 1)
		if cmp[i] != roots[i] {
			t.Error(fmt.Sprintf("error merkle root is not equal at %d", i))
		}
	}

}

func TestGetSubTreeSize(t *testing.T) {
	sizes := getSubTreeSize(7)
	fmt.Println("sub tree size", sizes)
}

// zero based return merkle root of D[0:n]
func TestMerkleIncludeProof(t *testing.T) {
	n := uint32(8)
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	for i := uint32(0); i < n; i++ {
		tree.Append([]byte{byte(i + 1)})
	}

	verify := NewMerkleVerifier()

	root := tree.Root()
	for i := uint32(0); i < n; i++ {
		proof := tree.InclusionProof(i, n)
		leaf_hash := tree.hasher.hash_leaf([]byte{byte(i + 1)})
		res := verify.VerifyLeafHashInclusion(leaf_hash, i, proof, root, n)
		if res != nil {
			t.Fatal(res, i, proof)
		}
	}
}

func TestMerkleConsistencyProofLen(t *testing.T) {
	n := uint32(7)
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	for i := uint32(0); i < n; i++ {
		tree.Append([]byte{byte(i + 1)})
	}

	cmp := []int{3, 2, 4, 1, 4, 3, 0}
	for i := uint32(0); i < n; i++ {
		proof := tree.ConsistencyProof(i+1, n)
		if len(proof) != cmp[i] {
			t.Fatal("error: wrong proof length")
		}
	}

}

func TestMerkleConsistencyProof(t *testing.T) {
	n := uint32(140)
	roots := make([]Uint256, n, n)
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	for i := uint32(0); i < n; i++ {
		tree.Append([]byte{byte(i + 1)})
		roots[i] = tree.Root()
	}

	verify := NewMerkleVerifier()

	for i := uint32(0); i < n; i++ {
		proof := tree.ConsistencyProof(i+1, n)
		err := verify.VerifyConsistency(i+1, n, roots[i], roots[n-1], proof)
		if err != nil {
			t.Fatal("verify consistency error:", i, err)
		}

	}
}

//~70w
func BenchmarkMerkleInsert(b *testing.B) {
	store := NewFileHashStore("merkletree.db")
	tree := NewTree(0, nil, store)
	for i := 0; i < b.N; i++ { //use b.N for looping
		tree.Append([]byte(fmt.Sprintf("bench %d", i)))
	}
}

var treeTest *CompactMerkleTree
var storeTest HashStore
var N = 100 //00

func init() {
	storeTest := NewFileHashStore("merkletree.db")
	treeTest := NewTree(0, nil, storeTest)
	for i := 0; i < N; i++ {
		treeTest.Append([]byte(fmt.Sprintf("setup %d", i)))
	}

}

/*
// ~20w
func BenchmarkMerkleInclusionProof(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeTest.InclusionProof(uint32(i), uint32(N))
	}
}

// ~20w
func BenchmarkMerkleConsistencyProof(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeTest.ConsistencyProof(uint32(i+1), uint32(N))
	}
}
*/

//~70w
func BenchmarkMerkleInsert2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeTest.Append([]byte(fmt.Sprintf("bench %d", i)))
	}
}

//

func TestNewFileSeek(t *testing.T) {
	name := "test.txt"
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal("can not open file", err)
	}
	off, err := f.Seek(0, 2)
	f.Write([]byte{12})
	t.Fatal(off, err)
}
