package pso

import (
	"github.com/alxdavids/bloom-filter/encbf"
	"log"
	"math"
	"math/big"
	"testing"
)

var (
	max       = 1200                // max size of elements (change this for influence over intersection size)
	n         = int(math.Pow(2, 8)) // size of sets
	maxConc   = 10000               // Maximum number of initiated goroutines
	keySize   = 1024                // key size for paillier
	mode      = 0                   // 0 = PSU, 1 = PSI, 2 = PSI/PSU-CA
	eps       = math.Pow(2, -25)    // false-positive prob for BF
	eblofCopy *encbf.EncBloom
	set1      []*big.Int
	set2      []*big.Int
)

func TestUnion(t *testing.T) {
	log.Println("******TESTING UNION******")
	set1 = generateSet(n, int64(max))
	set2 = generateSet(n, int64(max))

	newItems, _, eblof := computePSO(n, 0, keySize, max, maxConc, eps, set1, set2, nil)
	eblofCopy = eblof

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
	log.Println("******FINISHED UNION******")
}

func TestInter(t *testing.T) {
	log.Println("******TESTING INTERSECTION******")
	newItems, _, eblof := computePSO(n, 1, keySize, max, maxConc, eps, set1, set2, eblofCopy)
	eblofCopy = eblof

	// Check item exists in set2 and not in set1
	for _, v := range newItems {
		b1 := false
		for _, u := range set1 {
			if v.Cmp(u) == 0 {
				b1 = true
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
			log.Fatalln("Element not found in set1")
		}

		if !b2 {
			log.Println(v)
			log.Fatalln("Element not found in set2")
		}
	}
	log.Println("******FINISHED INTERSECTION******")
}

func TestCA(t *testing.T) {
	log.Println("******TESTING CARDINALITY******")
	_, count, _ := computePSO(n, 2, keySize, max, maxConc, eps, set1, set2, eblofCopy)

	// Check item exists in set2 and not in set1
	chkCount := 0
	for _, v := range set1 {
		for _, u := range set2 {
			if v.Cmp(u) == 0 {
				chkCount++
			}
		}
	}

	if chkCount != count {
		log.Println(set1)
		log.Println(set2)
		log.Println(chkCount)
		log.Println(count)
		log.Fatalln("Cardinality check incorrect")
	}
	log.Println("******FINISHED CARDINALITY******")
}
