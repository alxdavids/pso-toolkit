package pso

import (
	"crypto/rand"
	"github.com/alxdavids/bloom-filter/encbf"
	"github.com/alxdavids/bloom-filter/standard"
	"log"
	"math/big"
	"time"
	"xojoc.pw/bitset"
)

func computePSO(n, mode, keySize, max, maxConc int, eps float64, set1, set2 []*big.Int, eblof *encbf.EncBloom) ([]*big.Int, int, *encbf.EncBloom) {
	startTime := time.Now()

	// Give the option of providing a ready-made EBF for faster tests
	if eblof == nil {
		sblof := standard.New(uint(n), eps).(*standard.StandardBloom)

		// Add elements to original Bloom filter
		// We create random sets if non are provided
		if set1 == nil {
			set1 = generateSet(n, int64(max))
		}
		for _, v := range set1 {
			sblof.Add(v.Bytes())
		}

		eblof = encbf.New(sblof, keySize, mode, maxConc).(*encbf.EncBloom)
	}
	pub := eblof.GetPubKey()

	// generate new random set for checking
	if set2 == nil {
		set2 = generateSet(n, int64(max))
	}
	checkTime := time.Now()

	// Setup temp map for performing hom ops in concurrent exec
	for _, v := range set2 {
		eblof.Check(v.Bytes())
	}
	// Compute combined ciphertexts
	eblof.HomCombine()
	log.Println("Hom time: " + time.Since(checkTime).String())

	decTime := time.Now()
	ptxts := eblof.Decrypt()
	log.Println("Dec time: " + time.Since(decTime).String())

	// Generate correct output for set operation
	outTime := time.Now()
	items := []*big.Int{}
	count := 0
	if mode == 0 {
		// union
		for _, v := range ptxts {
			kinv := new(big.Int).ModInverse(new(big.Int).SetBytes(v[1]), pub.N)
			item := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).SetBytes(v[0]), kinv), pub.N)

			if item.Cmp(big.NewInt(0)) != 0 {
				items = append(items, item)
			}
		}
	} else if mode == 1 {
		// inter
		for _, v := range ptxts {
			item := big.NewInt(0)
			if new(big.Int).SetBytes(v[1]).Cmp(big.NewInt(0)) == 0 {
				item = new(big.Int).SetBytes(v[0])
			}

			if item.Cmp(big.NewInt(0)) != 0 {
				items = append(items, item)
			}
		}
	} else if mode == 2 {
		// cardinality
		for _, v := range ptxts {
			if new(big.Int).SetBytes(v[0]).Cmp(big.NewInt(0)) == 0 {
				count++
			}
		}
	}
	log.Println("Out time: " + time.Since(outTime).String())

	log.Println("Full time: " + time.Since(startTime).String())
	eblof.ResetForTesting()
	return items, count, eblof
}

// Generates random sets
func generateSet(n int, max int64) []*big.Int {
	arr := make([]*big.Int, n)
	used := &bitset.BitSet{}
	for i := 0; i < int(n); i++ {
		r, e := rand.Int(rand.Reader, big.NewInt(max))
		if e != nil {
			log.Fatalln(e)
		}

		// Make sure we get unique elements
		if !used.Get(int(r.Int64())) {
			arr[i] = r
		} else {
			// redo loop if not unique
			i--
		}
		used.Set(int(r.Int64()))
	}

	return arr
}
