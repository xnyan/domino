package workload

import (
	"github.com/op/go-logging"
	"math"
	"math/rand"
)

var logger = logging.MustGetLogger("workload")

var defaultValSize int = 1024

type Txn struct {
	TxnId     string
	ReadKeys  []string
	WriteData map[string]string
}

type Workload interface {
	GenTxn() *Txn
}

type AbstractWorkload struct {
	Workload // interface

	KeyList   []string  // all of the keys
	KeyNum    int64     // total number of the keys
	alpha     float64   // zipfian alpha value
	zipf      []float64 // zipfian values
	zipfReady bool      // If zipfian distribution has been initialized
	txnCount  int64

	valSize int // the size of value in Bytes
	val     string
}

func NewAbstractWorkload(
	keyList []string,
	zipfAlpha float64,
	valSize int,
) *AbstractWorkload {
	workload := &AbstractWorkload{
		KeyList: keyList,
		alpha:   zipfAlpha,
		valSize: valSize,
	}
	workload.KeyNum = int64(len(workload.KeyList))
	workload.zipf = nil
	workload.zipfReady = false
	workload.txnCount = 0

	if workload.valSize <= 0 {
		logger.Infof("Value size is not positive. Using the default %d", defaultValSize)
		workload.valSize = defaultValSize
	}

	workload.genVal()

	return workload
}

func (workload *AbstractWorkload) genVal() string {
	workload.val = ""
	for i := 0; i < workload.valSize; i++ {
		workload.val += "a"
	}
	return workload.val
}

func (workload *AbstractWorkload) getVal() string {
	return workload.val
}

// Currently, the read and write keys overlap, that is, one set is a subset of the other.
// TODO Generates transactions that have read and write keys not fully overlapped.
func (workload *AbstractWorkload) buildTxn(
	txnId string,
	rN, wN int,
) *Txn {
	txn := &Txn{
		TxnId:     txnId,
		ReadKeys:  make([]string, 0),
		WriteData: make(map[string]string),
	}

	max := wN
	if rN > max {
		max = rN
	}
	// Generates keys
	keyList := workload.genKeyList(max)

	// Read keys
	for i := 0; i < rN; i++ {
		txn.ReadKeys = append(txn.ReadKeys, keyList[i])
	}
	// Write keys
	for i := 0; i < wN; i++ {
		//txn.WriteData[keyList[i]] = workload.randKey()
		//txn.WriteData[keyList[i]] = strconv.Itoa(rand.Int())
		//txn.WriteData[keyList[i]] = keyList[i] // uses the key as the data, like in TAPIR's benchmark
		txn.WriteData[keyList[i]] = workload.getVal()
	}

	return txn
}

func (workload *AbstractWorkload) genKeyList(num int) []string {
	kList := make([]string, num)
	for i := 0; i < len(kList); i++ {
		kList[i] = workload.KeyList[workload.randKey()]
	}
	return kList
}

// Acknowledgement: this implementation is based on TAPIR's retwis benchmark.
func (workload *AbstractWorkload) randKey() int64 {
	if workload.alpha < 0 {
		// Uniform selection of keys.
		return rand.Int63n(workload.KeyNum)
	} else {
		// Zipf-like selection of keys.
		if !workload.zipfReady {
			workload.zipf = make([]float64, workload.KeyNum)

			var c float64 = 0.0
			var i int64
			for i = 1; i <= workload.KeyNum; i++ {
				c = c + (1.0 / math.Pow(float64(i), workload.alpha))
			}
			c = 1.0 / c

			var sum float64 = 0.0
			for i = 1; i <= workload.KeyNum; i++ {
				sum += (c / math.Pow(float64(i), workload.alpha))
				workload.zipf[i-1] = sum
			}
			workload.zipfReady = true
		}

		var rndNum float64 = 0.0
		for rndNum == 0.0 {
			rndNum = rand.Float64() //[0.0,1.0)
		}

		// Uses binary search to find the key's index
		var l, r, mid int64 = 0, workload.KeyNum, 0
		for l < r {
			if rndNum > workload.zipf[mid] {
				l = mid + 1
			} else if rndNum < workload.zipf[mid] {
				r = mid - 1
			} else {
				break
			}
			// Updates the mid for the last round
			mid = (l + r) / 2
		}

		if mid >= int64(len(workload.zipf)) {
			mid = int64(len(workload.zipf)) - 1
		}

		if workload.zipf[mid] < rndNum {
			// Takes the right one
			if mid+1 < int64(len(workload.zipf)) {
				return mid + 1
			}
		}
		return mid
	}
}

func (workload *AbstractWorkload) String() string {
	return "AbstractWorkload"
}
