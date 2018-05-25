package pso

import (
	"crypto/rand"
	"github.com/alxdavids/yabf/encbf"
	"github.com/alxdavids/yabf/standard"
	"log"
	"math/big"
	"time"
	"xojoc.pw/bitset"
)

func computePSO(n, mode, keySize, max, maxConc int, eps float64, set1, set2 []*big.Int, eblof *encbf.EncBloom, psoLog *log.Logger) ([]*big.Int, int, *encbf.EncBloom) {
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
	psoLog.Printf("Hom time: %v", time.Since(checkTime).Seconds())

	decTime := time.Now()
	ptxts := eblof.Decrypt()
	psoLog.Printf("Dec time: %v", time.Since(decTime).Seconds())

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
	psoLog.Printf("Out time: %v", time.Since(outTime).Seconds())

	psoLog.Printf("Full time: %v", time.Since(startTime).Seconds())
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
			psoLog.Fatalln(e)
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
