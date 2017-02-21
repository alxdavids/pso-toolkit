package pso

import (
	"log"
	"math"
	"testing"
)

var (
	max     = 1000             // max size of elements (change this for influence over intersection size)
	n       = 200              // size of sets
	keySize = 1024             // key size for paillier
	mode    = 0                // 0 = PSU, 1 = PSI, 2 = PSI/PSU-CA
	eps     = math.Pow(2, -20) // false-positive prob for BF
)

func TestUnion(t *testing.T) {
	set1 := generateSet(n, int64(max))
	set2 := generateSet(n, int64(max))

	newItems := computePSO(n, 0, keySize, max, eps, set1, set2)

	// Check item exists in set2 and not in set1
	for _, v := range newItems {
		b1 := true
		for _, u := range set1 {
			if v.Cmp(u) == 0 {
				b1 = false
			}
		}
		b2 := false
		for _, u := range set2 {
			if v.Cmp(u) == 0 {
				b2 = true
			}
		}

		if !b1 {
			log.Println(v)
			log.Fatalln("Element found in set1")
		}

		if !b2 {
			log.Println(v)
			log.Fatalln("Element not found in set2")
		}
	}
}
