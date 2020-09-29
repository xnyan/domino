package workload

import (
	"math/rand"
	"strconv"
)

// Acknowledgement: this implementation is based on TAPIR's retwis benchmark.
type RetwisWorkload struct {
	*AbstractWorkload

	add_user_ratio        int
	follow_unfollow_ratio int
	post_tweet_ratio      int
	load_timeline_ratio   int
}

func NewRetwisWorkload(
	workload *AbstractWorkload,
	retwis_add_user_ratio int,
	retwis_follow_unfollow_ratio int,
	retwis_post_tweet_ratio int,
	retwis_load_timeline_ratio int,
) *RetwisWorkload {
	retwis := &RetwisWorkload{
		AbstractWorkload:      workload,
		add_user_ratio:        retwis_add_user_ratio,
		follow_unfollow_ratio: retwis_follow_unfollow_ratio,
		post_tweet_ratio:      retwis_post_tweet_ratio,
		load_timeline_ratio:   retwis_load_timeline_ratio,
	}

	return retwis
}

// Generates a retwis txn. This function is currently not thread-safe
func (retwis *RetwisWorkload) GenTxn() *Txn {
	retwis.txnCount++
	txnId := strconv.FormatInt(retwis.txnCount, 10)

	txnType := rand.Intn(100) //[0,100)
	if txnType < retwis.add_user_ratio {
		// Add user txn. read 1, write 3 keys
		return retwis.buildTxn(txnId, 1, 3)
	} else if txnType <
		(retwis.add_user_ratio + retwis.follow_unfollow_ratio) {
		// Follow/Unfollow txn. read 2, write 2 keys
		return retwis.buildTxn(txnId, 2, 2)
	} else if txnType <
		(retwis.add_user_ratio + retwis.follow_unfollow_ratio + retwis.post_tweet_ratio) {
		// Post tweet txn. read 3, write 5 keys
		return retwis.buildTxn(txnId, 3, 5)
	} else if txnType <
		(retwis.add_user_ratio + retwis.follow_unfollow_ratio + retwis.post_tweet_ratio +
			retwis.load_timeline_ratio) {
		// Load timeline txn. read [1, 10] keys
		rN := rand.Intn(10) + 1 // [1,10]
		return retwis.buildTxn(txnId, rN, 0)
	} else {
		logger.Fatal("Txn generation error: uncovered percentage to generate a txn")
		return nil
	}
}

func (retwis *RetwisWorkload) String() string {
	return "RetwisWorkload"
}
