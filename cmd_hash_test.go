package miniredis

import (
	"sort"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

// Test Hash.
func TestHash(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var b bool
	var v string

	{
		b, err = c.HSet("aap", "noot", "mies").Result()
		ok(t, err)
		equals(t, true, b) // New field.
	}

	{
		v, err = c.HGet("aap", "noot").Result()
		ok(t, err)
		equals(t, "mies", v)
		equals(t, "mies", s.HGet("aap", "noot"))
	}

	{
		b, err = c.HSet("aap", "noot", "mies").Result()
		ok(t, err)
		equals(t, false, b) // Existing field.
	}

	// Wrong type of key
	{
		err = c.Set("foo", "bar", 0).Err()
		ok(t, err)
		err = c.HSet("foo", "noot", "mies").Err()
		assert(t, err != nil, "HSET error")
	}

	// hash exists, key doesn't.
	{
		v, err = c.HGet("app", "nosuch").Result()
		nilCheck(t, err)
	}

	// hash doesn't exists.
	{
		v, err = c.HGet("nosuch", "nosuch").Result()
		assert(t, err != nil, "")
		equals(t, "", s.HGet("nosuch", "nosuch"))
	}

	// Direct HSet()
	{
		s.HSet("wim", "zus", "jet")
		v, err = c.HGet("wim", "zus").Result()
		ok(t, err)
		equals(t, "jet", v)
	}
}

func TestHashSetNX(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// New Hash
	v, err := c.HSetNX("wim", "zus", "jet").Result()
	ok(t, err)
	equals(t, true, v)

	v, err = c.HSetNX("wim", "zus", "jet").Result()
	ok(t, err)
	equals(t, false, v)

	// Just a new key
	v, err = c.HSetNX("wim", "aap", "noot").Result()
	ok(t, err)
	equals(t, true, v)

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HSetNX("foo", "nosuch", "nosuch").Err()
	assert(t, err != nil, "no HSETNX error")
}

func TestHashMSet(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v string

	// New Hash
	{
		v, err = c.HMSet("hash", map[string]interface{} {
			"wim": "zus",
			"jet": "vuur",
		}).Result()
		ok(t, err)
		equals(t, "OK", v)

		equals(t, "zus", s.HGet("hash", "wim"))
		equals(t, "vuur", s.HGet("hash", "jet"))
	}

	// Doesn't touch ttl.
	{
		s.SetTTL("hash", time.Second*999)
		v, err = c.HMSet("hash", map[string]interface{} {
			"gijs": "lam",
		}).Result()
		ok(t, err)
		equals(t, "OK", v)
		equals(t, time.Second*999, s.TTL("hash"))
	}

	{
		// Wrong key type
		s.Set("str", "value")
		err = c.HMSet("str", map[string]interface{} {
			"key": "value",
		}).Err()
		assert(t, err != nil, "no HSETerror")
	}
}

func TestHashDel(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	s.HSet("wim", "gijs", "lam")
	s.HSet("wim", "kees", "bok")
	v, err := c.HDel("wim", "zus", "gijs").Result()
	ok(t, err)
	equals(t, 2, v)

	v, err = c.HDel("wim", "nosuch").Result()
	ok(t, err)
	equals(t, 0, v)

	// Deleting all makes the key disappear
	v, err = c.HDel("wim", "teun", "kees").Result()
	ok(t, err)
	equals(t, 2, v)
	assert(t, !s.Exists("wim"), "no more wim key")

	// Key doesn't exists.
	v, err = c.HDel("nosuch", "nosuch").Result()
	ok(t, err)
	equals(t, 0, v)

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HDel("foo", "nosuch").Err()
	assert(t, err != nil, "no HDEL error")

	// Direct HDel()
	s.HSet("aap", "noot", "mies")
	s.HDel("aap", "noot")
	equals(t, "", s.HGet("aap", "noot"))
}

func TestHashExists(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	v, err := c.HExists("wim", "zus").Result()
	ok(t, err)
	equals(t, true, v)

	v, err = c.HExists("wim", "nosuch").Result()
	ok(t, err)
	equals(t, false, v)

	v, err = c.HExists("nosuch", "nosuch").Result()
	ok(t, err)
	equals(t, false, v)

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HExists("foo", "nosuch").Err()
	assert(t, err != nil, "no HDEL error")
}

func TestHashGetall(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	s.HSet("wim", "gijs", "lam")
	s.HSet("wim", "kees", "bok")
	v, err := c.HGetAll("wim").Result()
	ok(t, err)
	equals(t, map[string]string{
		"zus":  "jet",
		"teun": "vuur",
		"gijs": "lam",
		"kees": "bok",
	}, v)

	v, err = c.HGetAll("nosuch").Result()
	ok(t, err)
	equals(t, 0, len(v))

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HGetAll("foo").Err()
	assert(t, err != nil, "no HGETALL error")
}

func TestHashKeys(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	s.HSet("wim", "gijs", "lam")
	s.HSet("wim", "kees", "bok")

	var v []string

	{
		v, err = c.HKeys("wim").Result()
		ok(t, err)
		equals(t, 4, len(v))
		sort.Strings(v)
		equals(t, []string{
			"gijs",
			"kees",
			"teun",
			"zus",
		}, v)
	}

	// Direct command
	{
		direct, err := s.HKeys("wim")
		ok(t, err)
		sort.Strings(direct)
		equals(t, []string{
			"gijs",
			"kees",
			"teun",
			"zus",
		}, direct)
		_, err = s.HKeys("nosuch")
		equals(t, err, ErrKeyNotFound)
	}

	v, err = c.HKeys("nosuch").Result()
	ok(t, err)
	equals(t, 0, len(v))

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HKeys("foo").Err()
	assert(t, err != nil, "no HKEYS error")
}

func TestHashValues(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	s.HSet("wim", "gijs", "lam")
	s.HSet("wim", "kees", "bok")
	v, err := c.HVals("wim").Result()
	ok(t, err)
	equals(t, 4, len(v))
	sort.Strings(v)
	equals(t, []string{
		"bok",
		"jet",
		"lam",
		"vuur",
	}, v)

	v, err = c.HVals("nosuch").Result()
	ok(t, err)
	equals(t, 0, len(v))

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HVals("foo").Err()
	assert(t, err != nil, "no HVALS error")
}

func TestHashLen(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	s.HSet("wim", "gijs", "lam")
	s.HSet("wim", "kees", "bok")
	v, err := c.HLen("wim").Result()
	ok(t, err)
	equals(t, 4, v)

	v, err = c.HLen("nosuch").Result()
	ok(t, err)
	equals(t, 0, v)

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HLen("foo").Err()
	assert(t, err != nil, "no HLEN error")
}

func TestHashMget(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.HSet("wim", "zus", "jet")
	s.HSet("wim", "teun", "vuur")
	s.HSet("wim", "gijs", "lam")
	s.HSet("wim", "kees", "bok")
	v, err := c.HMGet("wim", "zus", "nosuch", "kees").Result()
	ok(t, err)
	equals(t, 3, len(v))
	equals(t, "jet", v[0].(string))
	equals(t, nil, v[1])
	equals(t, "bok", v[2].(string))

	v, err = c.HMGet("nosuch", "zus", "kees").Result()
	ok(t, err)
	equals(t, 2, len(v))
	equals(t, nil, v[0])
	equals(t, nil, v[1])

	// Wrong key type
	s.Set("foo", "bar")
	err = c.HMGet("foo", "boo").Err()
	assert(t, err != nil, "no HMGET error")
}

func TestHashIncrby(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v int64

	// New key
	{
		v, err = c.HIncrBy("hash", "field", 1).Result()
		ok(t, err)
		equals(t, 1, v)
	}

	// Existing key
	{
		v, err = c.HIncrBy("hash", "field", 100).Result()
		ok(t, err)
		equals(t, 101, v)
	}

	// Minus works.
	{
		v, err = c.HIncrBy("hash", "field", -12).Result()
		ok(t, err)
		equals(t, 101-12, v)
	}

	// Direct usage
	s.HIncr("hash", "field", -3)
	equals(t, "86", s.HGet("hash", "field"))

	// Error cases.
	{
		// Wrong key type
		s.Set("str", "cake")
		err = c.HIncrBy("str", "case", 4).Err()
		assert(t, err != nil, "no HINCRBY error")
	}
}

func TestHashIncrbyfloat(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var v float64

	// Existing key
	{
		s.HSet("hash", "field", "12")
		v, err = c.HIncrByFloat("hash", "field", 400.12).Result()
		ok(t, err)
		equals(t, 412.12, v)
		equals(t, "412.12", s.HGet("hash", "field"))
	}

	// Existing key, not a number
	{
		s.HSet("hash", "field", "noint")
		v, err = c.HIncrByFloat("hash", "field", 400).Result()
		assert(t, err != nil, "do HINCRBYFLOAT error")
	}

	// New key
	{
		v, err = c.HIncrByFloat("hash", "newfield", 40.33).Result()
		ok(t, err)
		equals(t, 40.33, v)
		equals(t, "40.33", s.HGet("hash", "newfield"))
	}

	// Direct usage
	{
		s.HSet("hash", "field", "500.1")
		f, err := s.HIncrfloat("hash", "field", 12)
		ok(t, err)
		equals(t, 512.1, f)
		equals(t, "512.1", s.HGet("hash", "field"))
	}

	// Wrong type of existing key
	{
		s.Set("wrong", "type")
		err = c.HIncrByFloat("wrong", "type", 400).Err()
		assert(t, err != nil, "do HINCRBYFLOAT error")
	}
}

func TestHscan(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// We cheat with hscan. It always returns everything.

	s.HSet("h", "field1", "value1")
	s.HSet("h", "field2", "value2")

	var res []string
	var cur uint64

	// No problem
	{
		res, cur, err = c.HScan("h", 0, "", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"field1", "value1", "field2", "value2"}, res)
	}

	// Invalid cursor
	{
		res, cur, err = c.HScan("h", 42, "", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{}, res)
	}

	// COUNT (ignored)
	{
		res, cur, err = c.HScan("h", 0, "", 200).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"field1", "value1", "field2", "value2"}, res)
	}

	// MATCH
	{
		s.HSet("h", "aap", "a")
		s.HSet("h", "noot", "b")
		s.HSet("h", "mies", "m")
		res, cur, err = c.HScan("h", 0, "mi*", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"mies", "m"}, res)
	}
}
