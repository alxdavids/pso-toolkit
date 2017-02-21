package main

import (
	"crypto/rand"
	"github.com/alxdavids/bloom-filter/encbf"
	"github.com/alxdavids/bloom-filter/standard"
	"log"
	"math"
	"math/big"
)

var (
	max     = int64(1000)      // max size of elements (change this for influence over intersection size)
	n       = 100              // size of sets
	keySize = 512              // key size for paillier
	mode    = 0                // 0 = PSU, 1 = PSI, 2 = PSI/PSU-CA
	eps     = math.Pow(2, -20) // false-positive prob for BF
)

func main() {
	sblof := standard.New(uint(n), eps).(*standard.StandardBloom)

	// Add elements to original Bloom filter
	set := generateSet(n)
	for _, v := range set {
		sblof.Add(v.Bytes())
	}

	eblof := encbf.New(sblof, keySize, mode).(*encbf.EncBloom)
	pub := eblof.GetPubKey()

	// generate new random set for checking
	set2 := generateSet(n)
	for _, v := range set2 {
		eblof.Check(v.Bytes())
	}

	ptxts := eblof.Decrypt()

	for _, v := range ptxts {
		kinv := new(big.Int).ModInverse(new(big.Int).SetBytes(v[1]), pub.N)
		item := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).SetBytes(v[0]), kinv), pub.N)

		log.Println(item)
	}
}

func generateSet(n int) []*big.Int {
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
