package miniredis_test

import (
	"time"

	"github.com/dlsteuer/miniredis"
	"github.com/go-redis/redis"
)

func Example() {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// Configure you application to connect to redis at s.Addr()
	// Any redis client should work, as long as you use redis commands which
	// miniredis implements.
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	if _, err = c.Set("foo", "bar", 0).Result(); err != nil {
		panic(err)
	}

	// You can ask miniredis about keys directly, without going over the network.
	if got, err := s.Get("foo"); err != nil || got != "bar" {
		panic("Didn't get 'bar' back")
	}
	// Or with a DB id
	if _, err := s.DB(42).Get("foo"); err != miniredis.ErrKeyNotFound {
		panic("didn't use a different database")
	}

	// Test key with expiration
	s.SetTTL("foo", 60*time.Second)
	s.FastForward(60 * time.Second)
	if s.Exists("foo") {
		panic("expect key to be expired")
	}

	// Or use a Check* function which Fail()s if the key is not what we expect
	// (checks for existence, key type and the value)
	// s.CheckGet(t, "foo", "bar")

	// Check if there really was only one connection.
	if s.TotalConnectionCount() != 1 {
		panic("too many connections made")
	}

	// Output:
}
