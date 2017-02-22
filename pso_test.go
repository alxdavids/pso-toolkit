package pso

import (
	"log"
	"math"
	"testing"
)

var (
	max     = 5000             // max size of elements (change this for influence over intersection size)
	n       = 1000             // size of sets
	keySize = 1024             // key size for paillier
	mode    = 0                // 0 = PSU, 1 = PSI, 2 = PSI/PSU-CA
	eps     = math.Pow(2, -20) // false-positive prob for BF
)

func TestUnion(t *testing.T) {
	log.Println("******TESTING UNION******")
	set1 := generateSet(n, int64(max))
	set2 := generateSet(n, int64(max))

	newItems, _ := computePSO(n, 0, keySize, max, eps, set1, set2)

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
	set1 := generateSet(n, int64(max))
	set2 := generateSet(n, int64(max))

	newItems, _ := computePSO(n, 1, keySize, max, eps, set1, set2)

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
	set1 := generateSet(n, int64(max))
	set2 := generateSet(n, int64(max))

	_, count := computePSO(n, 2, keySize, max, eps, set1, set2)

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
