package main

import (
	"testing"

	"github.com/meteorhacks/bddp"
)

func BenchmarkPut(b *testing.B) {
	// create a client and connect
	c = bddp.NewClient()
	if err := c.Connect(Address); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := SendMetric(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}
