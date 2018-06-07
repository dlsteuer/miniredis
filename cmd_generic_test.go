package miniredis

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
)

// Test EXPIRE. Keys with an expiration are called volatile in Redis parlance.
func TestTTL(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var b time.Duration
	var n bool

	// Not volatile yet
	{
		equals(t, time.Duration(0), s.TTL("foo"))
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, -2 * time.Second, b)
	}

	// Set something
	{
		err = c.Set("foo", "bar", 0).Err()
		ok(t, err)
		// key exists, but no Expire set yet
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, -1 * time.Second, b)

		n, err = c.Expire("foo", 1200 * time.Second).Result()
		ok(t, err)
		equals(t, true, n) // EXPIRE returns 1 on success

		equals(t, 1200*time.Second, s.TTL("foo"))
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, 1200 * time.Second, b)
	}

	// A SET resets the expire.
	{
		err = c.Set("foo", "bar", 0).Err()
		ok(t, err)
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, -1 * time.Second, b)
	}

	// Set a non-existing key
	{
		n, err = c.Expire("nokey", 1200 * time.Second).Result()
		ok(t, err)
		equals(t, false, n) // EXPIRE returns 0 on failure.
	}

	// Remove an expire
	{

		// No key yet
		n, err = c.Persist("exkey").Result()
		ok(t, err)
		equals(t, false, n)

		err = c.Set("exkey", "bar", 0).Err()
		ok(t, err)

		// No timeout yet
		n, err = c.Persist("exkey").Result()
		ok(t, err)
		equals(t, false, n)

		err = c.Expire("exkey", 1200 * time.Second).Err()
		ok(t, err)

		// All fine now
		n, err = c.Persist("exkey").Result()
		ok(t, err)
		equals(t, true, n)

		// No TTL left
		b, err = c.TTL("exkey").Result()
		ok(t, err)
		equals(t, -1 * time.Second, b)
	}

	// Hash key works fine, too
	{
		err = c.HSet("wim", "zus", "iet").Err()
		ok(t, err)
		_, err = c.Expire("wim", 1234 * time.Second).Result()
		ok(t, err)
	}

	{
		err = c.Set("wim", "zus", 0).Err()
		ok(t, err)
		err = c.Expire("wim", -1200).Err()
		ok(t, err)
		equals(t, false, s.Exists("wim"))
	}
}

func TestExpireat(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var b time.Duration
	var n bool

	// Not volatile yet
	{
		equals(t, time.Duration(0), s.TTL("foo"))
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, -2 * time.Second, b)
	}

	// Set something
	{
		err = c.Set("foo", "bar", 0).Err()
		ok(t, err)
		// Key exists, but no ttl set.
		// b, err := redis.Int(c.Do("TTL", "foo"))
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, -1 * time.Second, b)

		s.SetTime(time.Unix(1234567890, 0))
		n, err = c.ExpireAt("foo", time.Unix(1234567890+100, 0)).Result()
		ok(t, err)
		equals(t, true, n) // EXPIREAT returns 1 on success.

		equals(t, 100*time.Second, s.TTL("foo"))
		b, err = c.TTL("foo").Result()
		ok(t, err)
		equals(t, 100 * time.Second, b)
		equals(t, 100*time.Second, s.TTL("foo"))
	}
}

func TestPexpire(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var b time.Duration
	var n bool

	// Key exists
	{
		ok(t, s.Set("foo", "bar"))
		// b, err := redis.Int(c.Do("PEXPIRE", "foo", 12))
		n, err = c.PExpire("foo", 12* time.Millisecond).Result()
		ok(t, err)
		equals(t, true, n)

		b, err = c.PTTL("foo").Result()
		ok(t, err)
		equals(t, 12 * time.Millisecond, b)

		equals(t, 12*time.Millisecond, s.TTL("foo"))
	}
	// Key doesn't exist
	{
		n, err = c.PExpire("nosuch", 12*time.Millisecond).Result()
		ok(t, err)
		equals(t, false, n)

		b, err = c.PTTL("nosuch").Result()
		ok(t, err)
		equals(t, -2 * time.Millisecond, b)
	}

	// No expire
	{
		s.Set("aap", "noot")
		b, err = c.PTTL("aap").Result()
		ok(t, err)
		equals(t, -1 * time.Millisecond, b)
	}
}

func TestDel(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.Set("foo", "bar")
	s.HSet("aap", "noot", "mies")
	s.Set("one", "two")
	s.SetTTL("one", time.Second*1234)
	s.Set("three", "four")
	r, err := c.Del("one", "aap", "nosuch").Result()
	ok(t, err)
	equals(t, int64(2), r)
	equals(t, time.Duration(0), s.TTL("one"))

	// Direct also works:
	s.Set("foo", "bar")
	s.Del("foo")
	got, err := s.Get("foo")
	equals(t, ErrKeyNotFound, err)
	equals(t, "", got)
}

func TestType(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v string

	// String key
	{
		s.Set("foo", "bar!")
		v, err = c.Type("foo").Result()
		ok(t, err)
		equals(t, "string", v)
	}

	// Hash key
	{
		s.HSet("aap", "noot", "mies")
		v, err = c.Type("aap").Result()
		ok(t, err)
		equals(t, "hash", v)
	}

	// New key
	{
		// v, err := redis.String(c.Do("TYPE", "nosuch"))
		v, err = c.Type("nosuch").Result()
		ok(t, err)
		equals(t, "none", v)
	}

	// Direct usage:
	{
		equals(t, "hash", s.Type("aap"))
		equals(t, "", s.Type("nokey"))
	}
}

func TestExists(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v int64

	// String key
	{
		s.Set("foo", "bar!")
		v, err = c.Exists("foo").Result()
		ok(t, err)
		equals(t, 1, v)
	}

	// Hash key
	{
		s.HSet("aap", "noot", "mies")
		// v, err := redis.Int(c.Do("EXISTS", "aap"))
		c.Exists("aap")
		ok(t, err)
		equals(t, 1, v)
	}

	// Multiple keys
	{
		v, err = c.Exists("foo", "aap").Result()
		ok(t, err)
		equals(t, 2, v)

		v, err = c.Exists("foo", "noot", "aap").Result()
		ok(t, err)
		equals(t, 2, v)
	}

	// New key
	{
		v, err = c.Exists("nosuch").Result()
		ok(t, err)
		equals(t, 0, v)
	}

	// Wrong usage
	{
		_, err = c.Exists().Result()
		assert(t, err != nil, "do EXISTS error")
	}

	// Direct usage:
	{
		equals(t, true, s.Exists("aap"))
		equals(t, false, s.Exists("nokey"))
	}
}

func TestMove(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v bool

	// No problem.
	{
		s.Set("foo", "bar!")
		v, err = c.Move("foo", 1).Result()
		ok(t, err)
		equals(t, true, v)
	}

	// Src key doesn't exists.
	{
		v, err = c.Move("nosuch", 1).Result()
		ok(t, err)
		equals(t, false, v)
	}

	// Target key already exists.
	{
		s.DB(0).Set("two", "orig")
		s.DB(1).Set("two", "taken")
		v, err = c.Move("two", 1).Result()
		ok(t, err)
		equals(t, false, v)
		s.CheckGet(t, "two", "orig")
	}

	// TTL is also moved
	{
		s.DB(0).Set("one", "two")
		s.DB(0).SetTTL("one", time.Second*4242)
		v, err = c.Move("one", 1).Result()
		ok(t, err)
		equals(t, true, v)
		equals(t, s.DB(1).TTL("one"), time.Second*4242)
	}
}

func TestKeys(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.Set("foo", "bar!")
	s.Set("foobar", "bar!")
	s.Set("barfoo", "bar!")
	s.Set("fooooo", "bar!")

	var v []string

	{
		v, err = c.Keys("foo").Result()
		ok(t, err)
		equals(t, []string{"foo"}, v)
	}

	// simple '*'
	{
		v, err = c.Keys("foo*").Result()
		ok(t, err)
		equals(t, []string{"foo", "foobar", "fooooo"}, v)
	}
	// simple '?'
	{
		v, err = c.Keys("fo?").Result()
		ok(t, err)
		equals(t, []string{"foo"}, v)
	}

	// Don't die on never-matching pattern.
	{
		v, err = c.Keys(`f\`).Result()
		ok(t, err)
		equals(t, []string{}, v)
	}
}

func TestRandom(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v string

	s.Set("one", "bar!")
	s.Set("two", "bar!")
	s.Set("three", "bar!")

	// No idea which key will be returned.
	{
		v, err = c.RandomKey().Result()
		ok(t, err)
		assert(t, v == "one" || v == "two" || v == "three", "RANDOMKEY looks sane")
	}
}

func TestRename(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// Non-existing key
	{
		err = c.Rename("nosuch", "to").Err()
		assert(t, err != nil, "do RENAME error")
	}

	// Same key
	{
		err = c.Rename("from", "from").Err()
		assert(t, err != nil, "do RENAME error")
	}

	var str string

	// Move a string key
	{
		s.Set("from", "value")
		str, err = c.Rename("from", "to").Result()
		ok(t, err)
		equals(t, "OK", str)
		equals(t, false, s.Exists("from"))
		equals(t, true, s.Exists("to"))
		s.CheckGet(t, "to", "value")
	}

	// Move a hash key
	{
		s.HSet("from", "key", "value")
		str, err = c.Rename("from", "to").Result()
		ok(t, err)
		equals(t, "OK", str)
		equals(t, false, s.Exists("from"))
		equals(t, true, s.Exists("to"))
		equals(t, "value", s.HGet("to", "key"))
	}

	// Move over something which exists
	{
		s.Set("from", "string value")
		s.HSet("to", "key", "value")
		s.SetTTL("from", time.Second*999999)

		str, err = c.Rename("from", "to").Result()
		ok(t, err)
		equals(t, "OK", str)
		equals(t, false, s.Exists("from"))
		equals(t, true, s.Exists("to"))
		s.CheckGet(t, "to", "string value")
		equals(t, time.Duration(0), s.TTL("from"))
		equals(t, time.Second*999999, s.TTL("to"))
	}
}

func TestScan(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// We cheat with scan. It always returns everything.

	s.Set("key", "value")

	var res []string
	var cur uint64

	// No problem
	{
		res, cur, err = c.Scan(0, "", 2).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"key"}, res)
	}

	// Invalid cursor
	{
		res, cur, err = c.Scan(42, "", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{}, res)
	}

	// COUNT (ignored)
	{
		res, cur, err = c.Scan(0, "", 200).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"key"}, res)
	}

	// MATCH
	{
		s.Set("aap", "noot")
		s.Set("mies", "wim")
		res, cur, err = c.Scan(0, "mi*", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"mies"}, res)
	}
}

func TestRenamenx(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// Non-existing key
	{
		err = c.RenameNX("nosuch", "to").Err()
		assert(t, err != nil, "do RENAMENX error")
	}

	// Same key
	{
		err = c.RenameNX("from", "from").Err()
		assert(t, err != nil, "do RENAMENX error")
	}

	var n bool

	// Move a string key
	{
		s.Set("from", "value")
		n, err = c.RenameNX("from", "to").Result()
		ok(t, err)
		equals(t, true, n)
		equals(t, false, s.Exists("from"))
		equals(t, true, s.Exists("to"))
		s.CheckGet(t, "to", "value")
	}

	// Move over something which exists
	{
		s.Set("from", "string value")
		s.Set("to", "value")

		n, err = c.RenameNX("from", "to").Result()
		ok(t, err)
		equals(t, false, n)
		equals(t, true, s.Exists("from"))
		equals(t, true, s.Exists("to"))
		s.CheckGet(t, "from", "string value")
		s.CheckGet(t, "to", "value")
	}
}
