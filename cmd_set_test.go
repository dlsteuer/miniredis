package miniredis

import (
	"sort"
	"testing"

	"github.com/go-redis/redis"
)

// Test SADD / SMEMBERS.
func TestSadd(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	var b int64
	var m []string
	var str string

	{
		b, err = c.SAdd("s", "aap", "noot", "mies").Result()
		ok(t, err)
		equals(t, 3, b) // New elements.

		members, err := s.Members("s")
		ok(t, err)
		equals(t, []string{"aap", "mies", "noot"}, members)

		m, err = c.SMembers("s").Result()
		ok(t, err)
		equals(t, []string{"aap", "mies", "noot"}, m)
	}

	{
		str, err = c.Type("s").Result()
		ok(t, err)
		equals(t, "set", str)
	}

	// SMEMBERS on an nonexisting key
	{
		m, err = c.SMembers("nosuch").Result()
		ok(t, err)
		equals(t, []string{}, m)
	}

	{
		b, err = c.SAdd("s", "new", "noot", "mies").Result()
		ok(t, err)
		equals(t, 1, b) // Only one new field.

		members, err := s.Members("s")
		ok(t, err)
		equals(t, []string{"aap", "mies", "new", "noot"}, members)
	}

	// Direct usage
	{
		added, err := s.SetAdd("s1", "aap")
		ok(t, err)
		equals(t, 1, added)

		members, err := s.Members("s1")
		ok(t, err)
		equals(t, []string{"aap"}, members)
	}

	// Wrong type of key
	{
		err = c.Set("str", "value", 0).Err()
		ok(t, err)
		err = c.SAdd("str", "hi").Err()
		assert(t, err != nil, "SADD error")
		err = c.SMembers("str").Err()
		assert(t, err != nil, "MEMBERS error")
	}

}

// Test SISMEMBER
func TestSismember(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s", "aap", "noot", "mies")

	var b bool

	{
		b, err = c.SIsMember("s", "aap").Result()
		ok(t, err)
		equals(t, true, b)

		b, err = c.SIsMember("s", "nosuch").Result()
		ok(t, err)
		equals(t, false, b)
	}

	// a nonexisting key
	{
		b, err = c.SIsMember("nosuch", "nosuch").Result()
		ok(t, err)
		equals(t, false, b)
	}

	// Direct usage
	{
		isMember, err := s.IsMember("s", "noot")
		ok(t, err)
		equals(t, true, isMember)
	}

	// Wrong type of key
	{
		err = c.Set("str", "value", 0).Err()
		ok(t, err)
		err = c.SIsMember("str", "thing").Err()
		assert(t, err != nil, "SISMEMBER error")
	}

}

// Test SREM
func TestSrem(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s", "aap", "noot", "mies", "vuur")

	var b int64

	{
		b, err = c.SRem("s", "aap", "noot").Result()
		ok(t, err)
		equals(t, 2, b)

		members, err := s.Members("s")
		ok(t, err)
		equals(t, []string{"mies", "vuur"}, members)
	}

	// a nonexisting field
	{
		b, err = c.SRem("s", "nosuch").Result()
		ok(t, err)
		equals(t, 0, b)
	}

	// a nonexisting key
	{
		b, err = c.SRem("nosuch", "nosuch").Result()
		ok(t, err)
		equals(t, 0, b)
	}

	// Direct usage
	{
		b, err := s.SRem("s", "mies")
		ok(t, err)
		equals(t, 1, b)

		members, err := s.Members("s")
		ok(t, err)
		equals(t, []string{"vuur"}, members)
	}

	// Wrong type of key
	{
		err = c.Set("str", "value", 0).Err()
		ok(t, err)
		err = c.SRem("str", "value").Err()
		assert(t, err != nil, "SREM error")
	}
}

// Test SMOVE
func TestSmove(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s", "aap", "noot")
	var b bool

	{
		b, err = c.SMove("s", "s2", "aap").Result()
		ok(t, err)
		equals(t, true, b)

		m, err := s.IsMember("s", "aap")
		ok(t, err)
		equals(t, false, m)
		m, err = s.IsMember("s2", "aap")
		ok(t, err)
		equals(t, true, m)
	}

	// Move away the last member
	{
		b, err = c.SMove("s", "s2", "noot").Result()
		ok(t, err)
		equals(t, true, b)

		equals(t, false, s.Exists("s"))

		m, err := s.IsMember("s2", "noot")
		ok(t, err)
		equals(t, true, m)
	}

	// a nonexisting member
	{
		b, err = c.SMove("s", "s2", "nosuch").Result()
		ok(t, err)
		equals(t, false, b)
	}

	// a nonexisting key
	{
		b, err = c.SMove("nosuch", "nosuch2", "nosuch").Result()
		ok(t, err)
		equals(t, false, b)
	}

	// Wrong type of key
	{
		err = c.Set("str", "value", 0).Err()
		ok(t, err)
		err = c.SMove("str", "dst", "value").Err()
		assert(t, err != nil, "SMOVE error")
		err = c.SMove("s2", "str", "value").Err()
		assert(t, err != nil, "SMOVE error")
	}
}

// Test SPOP
func TestSpop(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s", "aap", "noot")

	var el string

	{
		el, err = c.SPop("s").Result()
		ok(t, err)
		assert(t, el == "aap" || el == "noot", "spop got something")

		el, err = c.SPop("s").Result()
		ok(t, err)
		assert(t, el == "aap" || el == "noot", "spop got something")

		assert(t, !s.Exists("s"), "all spopped away")
	}

	// a nonexisting key
	{
		_, err = c.SPop("nosuch").Result()
		nilCheck(t, err)
	}

	// various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SPop("str").Err()
		assert(t, err != nil, "SPOP error")
	}
}

// Test SRANDMEMBER
func TestSrandmember(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s", "aap", "noot", "mies")

	var el string
	var els []string

	// No count
	{
		el, err = c.SRandMember("s").Result()
		ok(t, err)
		assert(t, el == "aap" || el == "noot" || el == "mies", "srandmember got something")
	}

	// Positive count
	{
		els, err = c.SRandMemberN("s", 2).Result()
		ok(t, err)
		equals(t, 2, len(els))
	}

	// Negative count
	{
		els, err = c.SRandMemberN("s", -2).Result()
		ok(t, err)
		equals(t, 2, len(els))
	}

	// a nonexisting key
	{
		_, err = c.SRandMember("nosuch").Result()
		nilCheck(t, err)
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SRandMember("str").Err()
		assert(t, err != nil, "SRANDMEMBER error")
	}
}

// Test SDIFF
func TestSdiff(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s1", "aap", "noot", "mies")
	s.SetAdd("s2", "noot", "mies", "vuur")
	s.SetAdd("s3", "aap", "mies", "wim")

	var els []string

	// Simple case
	{
		els, err = c.SDiff("s1", "s2").Result()
		ok(t, err)
		equals(t, []string{"aap"}, els)
	}

	// No other set
	{
		els, err = c.SDiff("s1").Result()
		ok(t, err)
		sort.Strings(els)
		equals(t, []string{"aap", "mies", "noot"}, els)
	}

	// 3 sets
	{
		els, err = c.SDiff("s1", "s2", "s3").Result()
		ok(t, err)
		equals(t, []string{}, els)
	}

	// A nonexisting key
	{
		els, err = c.SDiff("s9").Result()
		ok(t, err)
		equals(t, []string{}, els)
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SDiff().Err()
		assert(t, err != nil, "SDIFF error")
		err = c.SDiff("str").Err()
		assert(t, err != nil, "SDIFF error")
		err = c.SDiff("chk", "str").Err()
		assert(t, err != nil, "SDIFF error")
	}
}

// Test SDIFFSTORE
func TestSdiffstore(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s1", "aap", "noot", "mies")
	s.SetAdd("s2", "noot", "mies", "vuur")
	s.SetAdd("s3", "aap", "mies", "wim")

	var i int64

	// Simple case
	{
		i, err = c.SDiffStore("res", "s1", "s3").Result()
		ok(t, err)
		equals(t, 1, i)
		s.CheckSet(t, "res", "noot")
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SDiffStore("t").Err()
		assert(t, err != nil, "SDIFFSTORE error")
		err = c.SDiffStore("t", "str").Err()
		assert(t, err != nil, "SDIFFSTORE error")
	}
}

// Test SINTER
func TestSinter(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s1", "aap", "noot", "mies")
	s.SetAdd("s2", "noot", "mies", "vuur")
	s.SetAdd("s3", "aap", "mies", "wim")

	var els []string

	// Simple case
	{
		els, err = c.SInter("s1", "s2").Result()
		ok(t, err)
		sort.Strings(els)
		equals(t, []string{"mies", "noot"}, els)
	}

	// No other set
	{
		els, err = c.SInter("s1").Result()
		ok(t, err)
		sort.Strings(els)
		equals(t, []string{"aap", "mies", "noot"}, els)
	}

	// 3 sets
	{
		els, err = c.SInter("s1", "s2", "s3").Result()
		ok(t, err)
		equals(t, []string{"mies"}, els)
	}

	// A nonexisting key
	{
		els, err = c.SInter("s9").Result()
		ok(t, err)
		equals(t, []string{}, els)
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SInter().Err()
		assert(t, err != nil, "SINTER error")
		err = c.SInter("str").Err()
		assert(t, err != nil, "SINTER error")
	}
}

// Test SINTERSTORE
func TestSinterstore(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s1", "aap", "noot", "mies")
	s.SetAdd("s2", "noot", "mies", "vuur")
	s.SetAdd("s3", "aap", "mies", "wim")

	var i int64

	// Simple case
	{
		i, err = c.SInterStore("res", "s1", "s3").Result()
		ok(t, err)
		equals(t, 2, i)
		s.CheckSet(t, "res", "aap", "mies")
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SInterStore("t").Err()
		assert(t, err != nil, "SINTERSTORE error")
		err = c.SInterStore("t", "str").Err()
		assert(t, err != nil, "SINTERSTORE error")
	}
}

// Test SUNION
func TestSunion(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s1", "aap", "noot", "mies")
	s.SetAdd("s2", "noot", "mies", "vuur")
	s.SetAdd("s3", "aap", "mies", "wim")

	var els []string

	// Simple case
	{
		els, err = c.SUnion("s1", "s2").Result()
		ok(t, err)
		sort.Strings(els)
		equals(t, []string{"aap", "mies", "noot", "vuur"}, els)
	}

	// No other set
	{
		els, err = c.SUnion("s1").Result()
		ok(t, err)
		sort.Strings(els)
		equals(t, []string{"aap", "mies", "noot"}, els)
	}

	// 3 sets
	{
		els, err = c.SUnion("s1", "s2", "s3").Result()
		ok(t, err)
		sort.Strings(els)
		equals(t, []string{"aap", "mies", "noot", "vuur", "wim"}, els)
	}

	// A nonexisting key
	{
		els, err = c.SUnion("s9").Result()
		ok(t, err)
		equals(t, []string{}, els)
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SUnion().Err()
		assert(t, err != nil, "SUNION error")
		err = c.SUnion("str").Err()
		assert(t, err != nil, "SUNION error")
		err = c.SUnion("chk", "str").Err()
		assert(t, err != nil, "SUNION error")
	}
}

// Test SUNIONSTORE
func TestSunionstore(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.SetAdd("s1", "aap", "noot", "mies")
	s.SetAdd("s2", "noot", "mies", "vuur")
	s.SetAdd("s3", "aap", "mies", "wim")

	var i int64

	// Simple case
	{
		i, err = c.SUnionStore("res", "s1", "s3").Result()
		ok(t, err)
		equals(t, 4, i)
		s.CheckSet(t, "res", "aap", "mies", "noot", "wim")
	}

	// Various errors
	{
		s.SetAdd("chk", "aap", "noot")
		s.Set("str", "value")

		err = c.SUnionStore("t").Err()
		assert(t, err != nil, "SUNIONSTORE error")
		err = c.SUnionStore("t", "str").Err()
		assert(t, err != nil, "SUNIONSTORE error")
	}
}

func TestSscan(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	// We cheat with sscan. It always returns everything.

	s.SetAdd("set", "value1", "value2")

	var keys []string
	var cur uint64

	// No problem
	{
		keys, cur, err = c.SScan("set", 0, "", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"value1", "value2"}, keys)
	}

	// Invalid cursor
	{
		keys, _, err = c.SScan("set", 42, "", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{}, keys)
	}

	// COUNT (ignored)
	{
		keys, _, err = c.SScan("set", 0, "", 200).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"value1", "value2"}, keys)
	}

	// MATCH
	{
		s.SetAdd("set", "aap", "noot", "mies")
		keys, cur, err = c.SScan("set", 0, "mi*", 0).Result()
		ok(t, err)
		equals(t, 0, cur)
		equals(t, []string{"mies"}, keys)
	}
}
