package miniredis

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestWatch(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// Simple WATCH
	err = c.Watch(func(tx *redis.Tx) error {
		return nil
	}, "foo")
	ok(t, err)
}

// Test simple multi/exec block.
func TestSimpleTransaction(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	pipe := c.TxPipeline()

	b, err := pipe.Set("aap", 1, 0).Result()
	ok(t, err)

	// Not set yet.
	equals(t, false, s.Exists("aap"))

	res, err := pipe.Exec()
	ok(t, err)
	equals(t, 1, len(res))

	// SET should be back to normal mode
	b, err = c.Set("aap", 1, 0).Result()
	ok(t, err)
	equals(t, "OK", b)
}

func TestDiscardTransaction(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.Set("aap", "noot")

	pipe := c.TxPipeline()

	_, err = pipe.Set("aap", "mies", 0).Result()
	ok(t, err)

	// Not committed
	s.CheckGet(t, "aap", "noot")

	err = pipe.Discard()
	ok(t, err)

	// TX didn't get executed
	s.CheckGet(t, "aap", "noot")
}

func TestTxQueueErr(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	pipe := c.TxPipeline()

	_, err = pipe.Set("aap", "mies", 0).Result()
	ok(t, err)

	_, err = pipe.SAdd("wat").Result()
	ok(t, err)

	// This one is ok again
	_, err = pipe.Set("noot", "vuur", 0).Result()
	ok(t, err)

	_, err = pipe.Exec()
	assert(t, err != nil, "do EXEC error")

	// Didn't get EXECed
	equals(t, false, s.Exists("wat"))
}

func TestTxWatch(t *testing.T) {
	// Watch with no error.
	s, err := Run()
	ok(t, err)
	defer s.Close()

	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.Set("one", "two")
	err = c.Watch(func(tx *redis.Tx) error {
		return nil
	}, "one")
	ok(t, err)

	pipe := c.TxPipeline()

	_, err = pipe.Get("one").Result()
	ok(t, err)

	v, err := pipe.Exec()
	ok(t, err)
	equals(t, 1, len(v))
	equals(t, "two", v[0].(*redis.StringCmd).Val())
}