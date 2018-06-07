package miniredis
//
//import (
//	"testing"
//	"time"
//
//	"github.com/go-redis/redis"
//)
//
//func setup(t *testing.T) (*Miniredis, *redis.Client, func()) {
//	s, err := Run()
//	ok(t, err)
//	c1 := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//	return s, c1, func() { s.Close() }
//}
//
//func setup2(t *testing.T) (*Miniredis, *redis.Client, *redis.Client, func()) {
//	s, err := Run()
//	ok(t, err)
//	c1 := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//	c2 := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//	return s, c1, c2, func() { s.Close() }
//}
//
//func TestLpush(t *testing.T) {
//	s, c, done := setup(t)
//	defer done()
//
//	{
//		b, err := c.LPush("l", "aap", "noot", "mies").Result()
//		ok(t, err)
//		equals(t, 3, b) // New length.
//
//		r, err := c.LRange("l", 0, 0).Result()
//		ok(t, err)
//		equals(t, []string{"mies"}, r)
//
//		r, err = c.LRange("l", -1, -1).Result()
//		ok(t, err)
//		equals(t, []string{"aap"}, r)
//	}
//
//	// Push more.
//	{
//		b, err := c.LPush("l", "aap2", "noot2", "mies2").Result()
//		ok(t, err)
//		equals(t, 6, b) // New length.
//
//		r, err := c.LRange("l", 0, 0).Result()
//		ok(t, err)
//		equals(t, []string{"mies2"}, r)
//
//		r, err = c.LRange("l", -1, -1).Result()
//		ok(t, err)
//		equals(t, []string{"aap"}, r)
//	}
//
//	// Direct usage
//	{
//		l, err := s.Lpush("l2", "a")
//		ok(t, err)
//		equals(t, 1, l)
//		l, err = s.Lpush("l2", "b")
//		ok(t, err)
//		equals(t, 2, l)
//		list, err := s.List("l2")
//		ok(t, err)
//		equals(t, []string{"b", "a"}, list)
//
//		el, err := s.Lpop("l2")
//		ok(t, err)
//		equals(t, "b", el)
//		el, err = s.Lpop("l2")
//		ok(t, err)
//		equals(t, "a", el)
//		// Key is removed on pop-empty.
//		equals(t, false, s.Exists("l2"))
//	}
//
//	// Various errors
//	{
//		err := c.LPush("l").Err()
//		assert(t, err != nil, "LPUSH error")
//		err = c.Set("str", "value", 0).Err()
//		ok(t, err)
//		err = c.LPush("str", "noot", "mies").Err()
//		assert(t, err != nil, "LPUSH error")
//	}
//
//}
//
//func TestLpushx(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var b int64
//	var r []string
//
//	{
//		b, err = c.LPushX("l", "aap").Result()
//		ok(t, err)
//		equals(t, 0, b)
//		equals(t, false, s.Exists("l"))
//
//		// Create the list with a normal LPUSH
//		b, err = c.LPush("l", "noot").Result()
//		ok(t, err)
//		equals(t, 1, b)
//		equals(t, true, s.Exists("l"))
//
//		b, err = c.LPushX("l", "mies").Result()
//		ok(t, err)
//		equals(t, 2, b)
//		equals(t, true, s.Exists("l"))
//	}
//
//	// Push more.
//	{
//		b, err = c.LPush("l2", "aap1").Result()
//		ok(t, err)
//		equals(t, 1, b)
//		b, err = c.LPush("l2", "aap2", "noot2", "mies2").Result()
//		ok(t, err)
//		equals(t, 4, b)
//
//		r, err = c.LRange("l2", 0, 0).Result()
//		ok(t, err)
//		equals(t, []string{"mies2"}, r)
//
//		r, err = c.LRange("l2", -1, -1).Result()
//		ok(t, err)
//		equals(t, []string{"aap1"}, r)
//	}
//
//	// Errors
//	{
//		err = c.Set("str", "value", 0).Err()
//		ok(t, err)
//		err = c.LPushX("str", "mies").Err()
//		assert(t, err != nil, "LPUSHX error")
//	}
//
//}
//
//func TestLpop(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	b, err := c.LPush("l", "aap", "noot", "mies").Result()
//	ok(t, err)
//	equals(t, 3, b) // New length.
//
//	var el string
//	var i int64
//
//	// Simple pops.
//	{
//		el, err = c.LPop("l").Result()
//		ok(t, err)
//		equals(t, "mies", el)
//
//		el, err = c.LPop("l").Result()
//		ok(t, err)
//		equals(t, "noot", el)
//
//		el, err = c.LPop("l").Result()
//		ok(t, err)
//		equals(t, "aap", el)
//
//		// Last element has been popped. Key is gone.
//		i, err = c.Exists("l").Result()
//		ok(t, err)
//		equals(t, 0, i)
//
//		// Can pop non-existing keys just fine.
//		el, err = c.LPop("l").Result()
//		ok(t, err)
//		equals(t, nil, el)
//	}
//}
//
//func TestRPushPop(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var b int64
//	var r []string
//
//	{
//		b, err = c.RPush("l", "aap", "noot", "mies").Result()
//		ok(t, err)
//		equals(t, 3, b) // New length.
//
//		r, err = c.LRange("l", 0, 0).Result()
//		ok(t, err)
//		equals(t, []string{"aap"}, r)
//
//		r, err = c.LRange("l", -1, -1).Result()
//		ok(t, err)
//		equals(t, []string{"mies"}, r)
//	}
//
//	// Push more.
//	{
//		b, err = c.RPush("l", "aap2", "noop2", "mies2").Result()
//		ok(t, err)
//		equals(t, 6, b) // New length.
//
//		r, err = c.LRange("l", 0, 0).Result()
//		ok(t, err)
//		equals(t, []string{"aap"}, r)
//
//		r, err = c.LRange("l", -1, -1).Result()
//		ok(t, err)
//		equals(t, []string{"mies2"}, r)
//	}
//
//	// Direct usage
//	{
//		l, err := s.Push("l2", "a")
//		ok(t, err)
//		equals(t, 1, l)
//		l, err = s.Push("l2", "b")
//		ok(t, err)
//		equals(t, 2, l)
//		list, err := s.List("l2")
//		ok(t, err)
//		equals(t, []string{"a", "b"}, list)
//
//		el, err := s.Pop("l2")
//		ok(t, err)
//		equals(t, "b", el)
//		el, err = s.Pop("l2")
//		ok(t, err)
//		equals(t, "a", el)
//		// Key is removed on pop-empty.
//		equals(t, false, s.Exists("l2"))
//	}
//
//	// Wrong type of key
//	{
//		err = c.Set("key", "value", 0).Err()
//		ok(t, err)
//		err = c.RPush("str", "noot", "mies").Err()
//		assert(t, err != nil, "RPUSH error")
//	}
//
//}
//
//func TestRpop(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	s.Push("l", "aap", "noot", "mies")
//
//	var el string
//	var i int64
//
//	// Simple pops.
//	{
//		el, err = c.RPop("l").Result()
//		ok(t, err)
//		equals(t, "mies", el)
//
//		el, err = c.RPop("l").Result()
//		ok(t, err)
//		equals(t, "noot", el)
//
//		el, err = c.RPop("l").Result()
//		ok(t, err)
//		equals(t, "aap", el)
//
//		// Last element has been popped. Key is gone.
//		i, err = c.Exists("l").Result()
//		ok(t, err)
//		equals(t, 0, i)
//
//		// Can pop non-existing keys just fine.
//		el, err = c.RPop("l").Result()
//		ok(t, err)
//		equals(t, nil, el)
//	}
//}
//
//func TestLindex(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	s.Push("l", "aap", "noot", "mies", "vuur")
//
//	var el string
//
//	{
//		el, err = c.LIndex("l", 0).Result()
//		ok(t, err)
//		equals(t, "aap", el)
//	}
//	{
//		el, err = c.LIndex("l", 1).Result()
//		ok(t, err)
//		equals(t, "noot", el)
//	}
//	{
//		el, err = c.LIndex("l", 3).Result()
//		ok(t, err)
//		equals(t, "vuur", el)
//	}
//	// Too many
//	{
//		el, err = c.LIndex("l", 3000).Result()
//		ok(t, err)
//		equals(t, nil, el)
//	}
//	{
//		el, err = c.LIndex("l", -1).Result()
//		ok(t, err)
//		equals(t, "vuur", el)
//	}
//	{
//		el, err = c.LIndex("l", -2).Result()
//		ok(t, err)
//		equals(t, "mies", el)
//	}
//	// Too big
//	{
//		el, err = c.LIndex("l", -400).Result()
//		ok(t, err)
//		equals(t, nil, el)
//	}
//	// Non exising key
//	{
//		el, err = c.LIndex("nonexisting", 400).Result()
//		ok(t, err)
//		equals(t, nil, el)
//	}
//
//	// Wrong type of key
//	{
//		err = c.Set("str", "value", 0).Err()
//		ok(t, err)
//		err = c.LIndex("str", 1).Err()
//		assert(t, err != nil, "LINDEX error")
//	}
//}
//
//func TestLlen(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	s.Push("l", "aap", "noot", "mies", "vuur")
//
//	var el int64
//
//	{
//		el, err = c.LLen("l").Result()
//		ok(t, err)
//		equals(t, 4, el)
//	}
//
//	// Non exising key
//	{
//		el, err = c.LLen("nonexisting").Result()
//		ok(t, err)
//		equals(t, 0, el)
//	}
//
//	// Wrong type of key
//	{
//		err = c.Set("str", "value", 0).Err()
//		ok(t, err)
//		err = c.LLen("str").Err()
//		assert(t, err != nil, "LLEN error")
//	}
//}
//
//func TestLtrim(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	s.Push("l", "aap", "noot", "mies", "vuur")
//
//	var el string
//
//	{
//		el, err = c.LTrim("l", 0, 2).Result()
//		ok(t, err)
//		equals(t, "OK", el)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"aap", "noot", "mies"}, l)
//	}
//
//	// Delete key on empty list
//	{
//		el, err = c.LTrim("l", 0, -99).Result()
//		ok(t, err)
//		equals(t, "OK", el)
//		equals(t, false, s.Exists("l"))
//	}
//
//	// Non exising key
//	{
//		el, err = c.LTrim("nonexisting", 0, 1).Result()
//		ok(t, err)
//		equals(t, "OK", el)
//	}
//
//	// Wrong type of key
//	{
//		s.Set("str", "string!")
//		err = c.LTrim("str", 0, 1).Err()
//		assert(t, err != nil, "LTRIM error")
//	}
//}
//
//func TestLrem(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var n int64
//
//	// Reverse
//	{
//		s.Push("l", "aap", "noot", "mies", "vuur", "noot", "noot")
//		n, err = c.LRem("l", -1, "noot").Result()
//		ok(t, err)
//		equals(t, 1, n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"aap", "noot", "mies", "vuur", "noot"}, l)
//	}
//	// Normal
//	{
//		s.Push("l2", "aap", "noot", "mies", "vuur", "noot", "noot")
//		n, err = c.LRem("l2", 2, "noot").Result()
//		ok(t, err)
//		equals(t, 2, n)
//		l, err := s.List("l2")
//		ok(t, err)
//		equals(t, []string{"aap", "mies", "vuur", "noot"}, l)
//	}
//
//	// All
//	{
//		s.Push("l3", "aap", "noot", "mies", "vuur", "noot", "noot")
//		n, err = c.LRem("l3", 0, "noot").Result()
//		ok(t, err)
//		equals(t, 3, n)
//		l, err := s.List("l3")
//		ok(t, err)
//		equals(t, []string{"aap", "mies", "vuur"}, l)
//	}
//
//	// All
//	{
//		s.Push("l4", "aap", "noot", "mies", "vuur", "noot", "noot")
//		n, err = c.LRem("l4", 200, "noot").Result()
//		ok(t, err)
//		equals(t, 3, n)
//		l, err := s.List("l4")
//		ok(t, err)
//		equals(t, []string{"aap", "mies", "vuur"}, l)
//	}
//
//	// Delete key on empty list
//	{
//		s.Push("l5", "noot", "noot", "noot")
//		n, err = c.LRem("l5", 99, "noot").Result()
//		ok(t, err)
//		equals(t, 3, n)
//		equals(t, false, s.Exists("l5"))
//	}
//
//	// Non exising key
//	{
//		n, err = c.LRem("nonexisting", 0, "aap").Result()
//		ok(t, err)
//		equals(t, 0, n)
//	}
//
//	// Error cases
//	{
//		s.Set("str", "string!")
//		err = c.LRem("str", 0, "aap").Err()
//		assert(t, err != nil, "LREM error")
//	}
//}
//
//func TestLset(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	s.Push("l", "aap", "noot", "mies", "vuur", "noot", "noot")
//
//	var n string
//
//	// Simple LSET
//	{
//		n, err = c.LSet("l", 1, "noot!").Result()
//		ok(t, err)
//		equals(t, "OK", n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"aap", "noot!", "mies", "vuur", "noot", "noot"}, l)
//	}
//
//	{
//		n, err = c.LSet("l", -1, "noot?").Result()
//		ok(t, err)
//		equals(t, "OK", n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"aap", "noot!", "mies", "vuur", "noot", "noot?"}, l)
//	}
//
//	// Out of range
//	{
//		err = c.LSet("l", 10000, "aap").Err()
//		assert(t, err != nil, "LSET error")
//
//		err = c.LSet("l", -10000, "aap").Err()
//		assert(t, err != nil, "LSET error")
//	}
//
//	// Non exising key
//	{
//		err = c.LSet("nonexisting", 0, "aap").Err()
//		assert(t, err != nil, "LSET error")
//	}
//
//	// Error cases
//	{
//		s.Set("str", "string!")
//		err = c.LSet("str", 0, "aap").Err()
//		assert(t, err != nil, "LSET error")
//	}
//}
//
//func TestLinsert(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	s.Push("l", "aap", "noot", "mies", "vuur", "noot", "end")
//
//	var n int64
//
//	// Before
//	{
//		n, err = c.LInsert("l", "BEFORE", "noot", "!").Result()
//		ok(t, err)
//		equals(t, 7, n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"aap", "!", "noot", "mies", "vuur", "noot", "end"}, l)
//	}
//
//	// After
//	{
//		n, err = c.LInsert("l", "AFTER", "noot", "?").Result()
//		ok(t, err)
//		equals(t, 8, n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"aap", "!", "noot", "?", "mies", "vuur", "noot", "end"}, l)
//	}
//
//	// Edge case before
//	{
//		n, err = c.LInsert("l", "BEFORE", "aap", "[").Result()
//		ok(t, err)
//		equals(t, 9, n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"[", "aap", "!", "noot", "?", "mies", "vuur", "noot", "end"}, l)
//	}
//
//	// Edge case after
//	{
//		n, err = c.LInsert("l", "AFTER", "end", "]").Result()
//		ok(t, err)
//		equals(t, 10, n)
//		l, err := s.List("l")
//		ok(t, err)
//		equals(t, []string{"[", "aap", "!", "noot", "?", "mies", "vuur", "noot", "end", "]"}, l)
//	}
//
//	// Non exising pivot
//	{
//		n, err = c.LInsert("l", "before", "nosuch", "noot").Result()
//		ok(t, err)
//		equals(t, -1, n)
//	}
//
//	// Non exising key
//	{
//		n, err = c.LInsert("nonexisting", "before", "aap", "noot").Result()
//		ok(t, err)
//		equals(t, 0, n)
//	}
//
//	// Error cases
//	{
//		s.Set("str", "string!")
//		err = c.LInsert("str", "before", "value", "value").Err()
//		assert(t, err != nil, "LINSERT error")
//	}
//}
//
//func TestRpoplpush(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var n string
//
//	s.Push("l", "aap", "noot", "mies")
//	s.Push("l2", "vuur", "noot", "end")
//	{
//		n, err = c.RPopLPush("l", "l2").Result()
//		ok(t, err)
//		equals(t, "mies", n)
//		s.CheckList(t, "l", "aap", "noot")
//		s.CheckList(t, "l2", "mies", "vuur", "noot", "end")
//	}
//	// Again!
//	{
//		n, err = c.RPopLPush("l", "l2").Result()
//		ok(t, err)
//		equals(t, "noot", n)
//		s.CheckList(t, "l", "aap")
//		s.CheckList(t, "l2", "noot", "mies", "vuur", "noot", "end")
//	}
//	// Again!
//	{
//		n, err = c.RPopLPush("l", "l2").Result()
//		ok(t, err)
//		equals(t, "aap", n)
//		assert(t, !s.Exists("l"), "l exists")
//		s.CheckList(t, "l2", "aap", "noot", "mies", "vuur", "noot", "end")
//	}
//
//	// Non exising lists
//	{
//		s.Push("ll", "aap", "noot", "mies")
//
//		n, err = c.RPopLPush("ll", "nosuch").Result()
//		ok(t, err)
//		equals(t, "mies", n)
//		assert(t, s.Exists("nosuch"), "nosuch exists")
//		s.CheckList(t, "ll", "aap", "noot")
//		s.CheckList(t, "nosuch", "mies")
//
//		n, err = c.RPopLPush("nosuch2", "ll").Result()
//		ok(t, err)
//		equals(t, nil, n)
//	}
//
//	// Cycle
//	{
//		s.Push("cycle", "aap", "noot", "mies")
//
//		n, err = c.RPopLPush("cycle", "cycle").Result()
//		ok(t, err)
//		equals(t, "mies", n)
//		s.CheckList(t, "cycle", "mies", "aap", "noot")
//	}
//
//	// Error cases
//	{
//		s.Push("src", "aap", "noot", "mies")
//
//		s.Set("str", "string!")
//		err = c.RPopLPush("str", "src").Err()
//		assert(t, err != nil, "RPOPLPUSH error")
//		err = c.RPopLPush("src", "str").Err()
//		assert(t, err != nil, "RPOPLPUSH error")
//	}
//}
//
//func TestRpushx(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var i, b int64
//	var r []string
//
//	// Simple cases
//	{
//		// No key key
//		i, err = c.RPushX("l", "value").Result()
//		ok(t, err)
//		equals(t, 0, i)
//		assert(t, !s.Exists("l"), "l doesn't exist")
//
//		s.Push("l", "aap", "noot")
//
//		i, err = c.RPushX("l", "mies").Result()
//		ok(t, err)
//		equals(t, 3, i)
//
//		s.CheckList(t, "l", "aap", "noot", "mies")
//	}
//
//	// Push more.
//	{
//		b, err = c.LPush("l2", "aap1").Result()
//		ok(t, err)
//		equals(t, 1, b)
//		b, err = c.RPush("l2", "aap2", "noot2", "mies2").Result()
//		ok(t, err)
//		equals(t, 4, b)
//
//		r, err = c.LRange("l2", 0, 0).Result()
//		ok(t, err)
//		equals(t, []string{"aap1"}, r)
//
//		r, err = c.LRange("l2", -1, -1).Result()
//		ok(t, err)
//		equals(t, []string{"mies2"}, r)
//	}
//
//	// Error cases
//	{
//		s.Set("str", "string!")
//		err = c.RPushX("str", "value").Err()
//		assert(t, err != nil, "RPUSHX error")
//	}
//}
//
//// execute command in a go routine. Used to test blocking commands.
//func goStrings(t *testing.T, c *redis.Client, cmds ...interface{}) <-chan []string {
//	var (
//		got = make(chan []string, 1)
//	)
//	go func() {
//		cmd := redis.NewStringCmd(cmds...)
//		err := c.Process(cmd)
//		if err != nil {
//			got <- []string{err.Error()}
//			return
//		}
//		if cmd.Err() != nil {
//			got <- []string{cmd.Err().Error()}
//			return
//		}
//		if cmd == nil {
//			got <- nil
//		} else {
//			st, _ := cmd.Result()
//			got <- []string{st}
//		}
//	}()
//	return got
//}
//
//func TestBrpop(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var v []string
//
//	// Simple cases
//	{
//		s.Push("ll", "aap", "noot", "mies")
//		v, err = c.BRPop(1, "ll").Result()
//		ok(t, err)
//		equals(t, []string{"ll", "mies"}, v)
//	}
//
//	// Error cases
//	{
//		err = c.BRPop(-1, "ll").Err()
//		assert(t, err != nil, "BRPOP error")
//	}
//}
//
//func TestBrpopSimple(t *testing.T) {
//	_, c1, c2, done := setup2(t)
//	defer done()
//
//	got := goStrings(t, c2, "BRPOP", "mylist", "0")
//	time.Sleep(30 * time.Millisecond)
//
//	b, err := c1.RPush("mylist", "e1", "e2", "e3").Result()
//	ok(t, err)
//	equals(t, 3, b)
//
//	select {
//	case have := <-got:
//		equals(t, []string{"mylist", "e3"}, have)
//	case <-time.After(500 * time.Millisecond):
//		t.Error("BRPOP took too long")
//	}
//}
//
//func TestBrpopMulti(t *testing.T) {
//	_, c1, c2, done := setup2(t)
//	defer done()
//
//	got := goStrings(t, c2, "BRPOP", "l1", "l2", "l3", 0)
//	err := c1.RPush("l0", "e01").Err()
//	ok(t, err)
//	err = c1.RPush("l2", "e21").Err()
//	ok(t, err)
//	err = c1.RPush("l3", "e31").Err()
//	ok(t, err)
//
//	select {
//	case have := <-got:
//		equals(t, []string{"l2", "e21"}, have)
//	case <-time.After(500 * time.Millisecond):
//		t.Error("BRPOP took too long")
//	}
//
//	got = goStrings(t, c2, "BRPOP", "l1", "l2", "l3", 0)
//	select {
//	case have := <-got:
//		equals(t, []string{"l3", "e31"}, have)
//	case <-time.After(500 * time.Millisecond):
//		t.Error("BRPOP took too long")
//	}
//}
//
//func TestBrpopTimeout(t *testing.T) {
//	_, c, done := setup(t)
//	defer done()
//
//	got := goStrings(t, c, "BRPOP", "l1", 1)
//	select {
//	case have := <-got:
//		equals(t, []string(nil), have)
//	case <-time.After(1500 * time.Millisecond):
//		t.Error("BRPOP took too long")
//	}
//}
//
//func TestBrpopTx(t *testing.T) {
//	// BRPOP in a transaction behaves as if the timeout triggers right away
//	m, c, done := setup(t)
//	defer done()
//
//	{
//		pipe := c.TxPipeline()
//		_, err := pipe.BRPop(3, "l1").Result()
//		ok(t, err)
//		_, err = pipe.Set("foo", "bar", 0).Result()
//		ok(t, err)
//
//		v, err := pipe.Exec()
//		ok(t, err)
//		equals(t, 2, len(v))
//		equals(t, nil, v[0])
//		equals(t, "OK", v[1])
//	}
//
//	// Now set something
//	m.Push("l1", "e1")
//
//	{
//		pipe := c.TxPipeline()
//		_, err := pipe.BRPop(3, "l1").Result()
//		ok(t, err)
//		_, err = pipe.Set("foo", "bar", 0).Result()
//		ok(t, err)
//
//		v, err := pipe.Exec()
//		ok(t, err)
//		equals(t, 2, len(v))
//		equals(t, "l1", v[0])
//		equals(t, "e1", v[0])
//		equals(t, "OK", v[1])
//	}
//}
//
//func TestBlpop(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var v []string
//
//	// Simple cases
//	{
//		s.Push("ll", "aap", "noot", "mies")
//		v, err = c.BLPop(1, "ll").Result()
//		ok(t, err)
//		equals(t, []string{"ll", "aap"}, v)
//	}
//
//	// Error cases
//	{
//		err = c.BLPop(-1, "key").Err()
//		assert(t, err != nil, "BLPOP error")
//	}
//}
//
//func TestBrpoplpush(t *testing.T) {
//	s, err := Run()
//	ok(t, err)
//	defer s.Close()
//	c := redis.NewClient(&redis.Options{
//		Network: "tcp",
//		Addr:    s.Addr(),
//	})
//
//	var v string
//
//	// Simple cases
//	{
//		s.Push("l1", "aap", "noot", "mies")
//		v, err = c.BRPopLPush("l1", "l2", 1).Result()
//		ok(t, err)
//		equals(t, "mies", v)
//
//		lv, err := s.List("l2")
//		ok(t, err)
//		equals(t, []string{"mies"}, lv)
//	}
//
//	// Error cases
//	{
//		err = c.BRPopLPush("key", "foo", -1).Err()
//		assert(t, err != nil, "BRPOPLPUSH error")
//	}
//}
//
//func TestBrpoplpushSimple(t *testing.T) {
//	s, c1, c2, done := setup2(t)
//	defer done()
//
//	got := make(chan string, 1)
//	go func() {
//		b, err := c2.BRPopLPush("from", "to", 1).Result()
//		ok(t, err)
//		got <- b
//	}()
//
//	time.Sleep(30 * time.Millisecond)
//
//	b, err := c1.RPush("from", "e1", "e2", "e3").Result()
//	ok(t, err)
//	equals(t, 3, b)
//
//	select {
//	case have := <-got:
//		equals(t, "e3", have)
//	case <-time.After(500 * time.Millisecond):
//		t.Error("BRPOP took too long")
//	}
//
//	lv, err := s.List("from")
//	ok(t, err)
//	equals(t, []string{"e1", "e2"}, lv)
//	lv, err = s.List("to")
//	ok(t, err)
//	equals(t, []string{"e3"}, lv)
//}
//
//func TestBrpoplpushTimeout(t *testing.T) {
//	_, c, done := setup(t)
//	defer done()
//
//	got := goStrings(t, c, "BRPOPLPUSH", "l1", "l2", 1)
//	select {
//	case have := <-got:
//		equals(t, []string(nil), have)
//	case <-time.After(1500 * time.Millisecond):
//		t.Error("BRPOPLPUSH took too long")
//	}
//}
