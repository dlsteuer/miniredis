package miniredis

import (
	"testing"

	"github.com/go-redis/redis"
)

// Test DBSIZE, FLUSHDB, and FLUSHALL.
func TestCmdServer(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// Set something
	{
		s.Set("aap", "niet")
		s.Set("roos", "vuur")
		s.DB(1).Set("noot", "mies")
	}

	var n int64
	var b string

	{
		n, err = c.DBSize().Result()
		ok(t, err)
		equals(t, 2, n)

		b, err = c.FlushDB().Result()
		ok(t, err)
		equals(t, "OK", b)

		n, err = c.DBSize().Result()
		ok(t, err)
		equals(t, 0, n)

		c := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    s.Addr(),
			DB: 1,
		})

		n, err = c.DBSize().Result()
		ok(t, err)
		equals(t, 1, n)

		b, err = c.FlushAll().Result()
		ok(t, err)
		equals(t, "OK", b)

		n, err = c.DBSize().Result()
		ok(t, err)
		equals(t, 0, n)

		c = redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    s.Addr(),
			DB: 4,
		})

		n, err = c.DBSize().Result()
		ok(t, err)
		equals(t, 0, n)
	}

	c = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	{
		b, err = c.FlushDBAsync().Result()
		ok(t, err)
		equals(t, "OK", b)

		b, err = c.FlushAllAsync().Result()
		ok(t, err)
		equals(t, "OK", b)
	}
}
