package pso

import (
	"crypto/rand"
	"github.com/alxdavids/bloom-filter/encbf"
	"github.com/alxdavids/bloom-filter/standard"
	"log"
	"math/big"
	"time"
)

func computePSO(n, mode, keySize, max int, eps float64, set1, set2 []*big.Int) []*big.Int {
	startTime := time.Now()
	sblof := standard.New(uint(n), eps).(*standard.StandardBloom)

	// Add elements to original Bloom filter
	if set1 == nil {
		set1 = generateSet(n, int64(max))
	}
	for _, v := range set1 {
		sblof.Add(v.Bytes())
	}

	eblof := encbf.New(sblof, keySize, mode).(*encbf.EncBloom)
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

	outTime := time.Now()
	items := []*big.Int{}
	for _, v := range ptxts {
		kinv := new(big.Int).ModInverse(new(big.Int).SetBytes(v[1]), pub.N)
		item := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).SetBytes(v[0]), kinv), pub.N)

		if item.Cmp(big.NewInt(0)) != 0 {
			items = append(items, item)
		}
	}
	log.Println("Out time: " + time.Since(outTime).String())

	log.Println("Full time: " + time.Since(startTime).String())
	return items
}

func generateSet(n int, max int64) []*big.Int {
	arr := make([]*big.Int, n)
	for i := 0; i < int(n); i++ {
		r, e := rand.Int(rand.Reader, big.NewInt(max))
		if e != nil {
			log.Fatalln(e)
		}
		arr[i] = r
	}

	return arr
}
