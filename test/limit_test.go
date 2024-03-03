package test

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/ratelimit"
)

func TestLimit(t *testing.T) {
	// 一秒多少滴水
	limiter := ratelimit.New(1)
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := limiter.Take()
		if time.Until(now) > 0 {
			fmt.Println("ttt")
		}
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
