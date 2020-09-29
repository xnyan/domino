package latency

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestLatManager(t *testing.T) {
	wdwLen, _ := time.ParseDuration("1s")
	lm := NewLatManager(wdwLen, 10)
	wait, _ := time.ParseDuration("10ms")
	fmt.Printf("%d, window 95th %dms, all 95th %dms\n", 0, lm.GetWindow95th(), lm.GetAll95th())
	for i := 1; i <= 20; i++ {
		lat, _ := time.ParseDuration(strconv.Itoa(i) + "ms")
		lm.AddLat(lat)
		time.Sleep(wait)
		//t.Logf("%d, window 95th %dms, all 95th %dms\n", i, lm.GetWindow95th(), lm.GetAll95th())
		fmt.Printf("%d, window 95th %dms, all 95th %dms\n", i, lm.GetWindow95th(), lm.GetAll95th())
	}
}
