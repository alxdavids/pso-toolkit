package pso

import (
	"flag"
	"github.com/alxdavids/bloom-filter/encbf"
	"log"
	"math"
	"math/big"
	"runtime"
	"testing"
)

var (
	domain    int                // domain size of elements (change this for influence over intersection size)
	n         int                // size of set
	maxProcs  int                // Max number of threads
	maxConc   = 10000            // Maximum number of initiated goroutines
	keySize   = 1024             // key size for paillier
	mode      = 0                // 0 = PSU, 1 = PSI, 2 = PSI/PSU-CA
	eps       = math.Pow(2, -50) // false-positive prob for BF
	set1      []*big.Int         // set stored in blof
	set2      []*big.Int         // set used for querying
	eblofCopy *encbf.EncBloom    // Used for redoing tests without re-encrypting
	outFile   string             // logging goes to a file
)

func init() {
	flag.IntVar(&keySize, "k", 1024, "Sets the key size, choose 1024 or 2048")
	flag.IntVar(&n, "n", 64, "Sets the set size")
	flag.IntVar(&maxProcs, "m", 4, "Sets the max number of threads to use")
	flag.StringVar(&outFile, "f", "", "File name for log output")
	prev := runtime.GOMAXPROCS(maxProcs)
	log.Printf("Previous number of threads used: %v\n", prev)

	// Enable logging to file
	f, err := os.OpenFile("outFile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}

func TestUnion(t *testing.T) {
	log.Printf("Max number of threads: %v\n", maxProcs)
	log.Printf("Key size: %v\n", keySize)
	log.Println("******TESTING UNION******")

	// set the size of the domain here
	domain = 5 * n
	set1 = generateSet(n, int64(domain))
	set2 = generateSet(n, int64(domain))

	newItems, _, eblof := computePSO(n, 0, keySize, domain, maxConc, eps, set1, set2, nil)
	eblof.DumpParams()
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
	newItems, _, eblof := computePSO(n, 1, keySize, domain, maxConc, eps, set1, set2, eblofCopy)
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
	_, count, _ := computePSO(n, 2, keySize, domain, maxConc, eps, set1, set2, eblofCopy)

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
