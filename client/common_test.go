package main

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkPut(b *testing.B) {
	c := Prepare()

	ts := time.Now().UnixNano()
	pld := make([]byte, 16, 16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pno := int64(i % 4)

		vals := []string{
			"test-app",
			"test-type",
			"test-host",
			"test-data",
		}

		vals[i%4] = vals[i%4] + strconv.Itoa(i)

		err := SendMetric(c, ts, pno, vals, pld)
		if err != nil {
			b.Fatal(err)
		}
	}
}
