package pso

import (
	"flag"
	"github.com/alxdavids/yabf/encbf"
	"log"
	"math"
	"math/big"
	"os"
	"runtime"
	"testing"
)

var (
	domainFactor int                                                          // domain size of elements (change this for influence over intersection size)
	domainSize   int                                                          // domain size of elements (change this for influence over intersection size)
	n            int                                                          // size of set
	maxProcs     int                                                          // Max number of threads
	maxConc      int                                                          // Maximum number of initiated goroutines
	keySize      int                                                          // key size for paillier
	k            int                                                          // -log_2(k) = eps (number of hash functions)
	mode         = 0                                                          // 0 = PSU, 1 = PSI, 2 = PSI/PSU-CA, 3 = all
	eps          float64                                                      // false-positive prob for BF
	set1         []*big.Int                                                   // set stored in blof
	set2         []*big.Int                                                   // set used for querying
	currEblof    *encbf.EncBloom                                        = nil // Used for redoing tests without re-encrypting
	outFile      string                                                       // logging output file
	psoLog       = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)       // logger for file
)

func init() {
	flag.IntVar(&keySize, "key_length", 1024, "Sets the key size, choose 1024 or 2048")
	flag.IntVar(&n, "set_size", 64, "Size of the sets considered")
	flag.IntVar(&k, "false_positive", 30, "False positive probability (-log_2)")
	flag.IntVar(&maxProcs, "max_threads", 4, "Sets the max number of threads to use")
	flag.IntVar(&maxConc, "max_conc", 10000, "Sets the max number of goroutines")
	flag.IntVar(&domainFactor, "domain_factor", 5, "Size of domain (actual_domain_size = domain_factor*n)")
	flag.IntVar(&mode, "mode", 0, "Mode (0 = PSU, 1 = PSI, 2 = PSU/I-CA), 3 = all")
	flag.StringVar(&outFile, "out", "", "File name for log output")
	prev := runtime.GOMAXPROCS(maxProcs)

	eps = math.Pow(2, -30)

	// Print params
	log.Printf("Previous number of threads used: %v\n", prev)
	log.Printf("Max number of threads: %v\n", maxProcs)
	log.Printf("Key size: %v\n", keySize)
	log.Printf("Set size: %v\n", n)
	log.Printf("False positive: %v\n", k)

	// open log file if specified
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		psoLog = log.New(f, "", log.LstdFlags|log.Lshortfile)
		defer f.Close()
	}

	// Generate the sets that are to be used
	domainSize := int64(n * domainFactor)
	set1 = generateSet(n, int64(domainSize))
	set2 = generateSet(n, int64(domainSize))
}

func TestUnion(t *testing.T) {
	if mode != 0 && mode != 3 {
		psoLog.Printf("Not testing union; mode=%v", mode)
		return
	}

	psoLog.Println("******TESTING UNION******")
	newItems, _, eblof := computePSO(n, 0, keySize, domainSize, maxConc, eps, set1, set2, currEblof, psoLog)
	eblof.DumpParams()
	currEblof = eblof

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
			psoLog.Println(v)
			log.Fatalln("Element found in set1")
		}

		if !b2 {
			psoLog.Println(v)
			log.Fatalln("Element not found in set2")
		}
	}
	psoLog.Println("******FINISHED UNION******")
}

func TestInter(t *testing.T) {
	if mode != 1 && mode != 3 {
		psoLog.Printf("Not testing intersection; mode=%v", mode)
		return
	}

	psoLog.Println("******TESTING INTERSECTION******")
	newItems, _, eblof := computePSO(n, 1, keySize, domainSize, maxConc, eps, set1, set2, currEblof, psoLog)
	currEblof = eblof

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
			psoLog.Println(v)
			log.Fatalln("Element not found in set1")
		}

		if !b2 {
			psoLog.Println(v)
			log.Fatalln("Element not found in set2")
		}
	}
	psoLog.Println("******FINISHED INTERSECTION******")
}

func TestCA(t *testing.T) {
	if mode != 2 && mode != 3 {
		psoLog.Printf("Not testing cardinality; mode=%v", mode)
		return
	}
	psoLog.Println("******TESTING CARDINALITY******")
	_, count, _ := computePSO(n, 2, keySize, domainSize, maxConc, eps, set1, set2, currEblof, psoLog)

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
		psoLog.Println(chkCount)
		psoLog.Println(count)
		log.Fatalln("Cardinality check incorrect")
	}
	psoLog.Println("******FINISHED CARDINALITY******")
}
