package miniredis

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestAuth(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	s.RequireAuth("nocomment")
	err = c.Ping().Err()
	mustFail(t, err, "NOAUTH Authentication required.")

	c = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     s.Addr(),
		Password: "wrongpasswd",
	})
	_, err = c.Ping().Result()
	mustFail(t, err, "ERR invalid password")

	c = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     s.Addr(),
		Password: "nocomment",
	})

	err = c.Ping().Err()
	ok(t, err)
}

func TestEcho(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	r, err := c.Echo("hello\nworld").Result()
	ok(t, err)
	equals(t, "hello\nworld", r)
}

func TestSelect(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	err = c.Set("foo", "bar", 0).Err()
	ok(t, err)

	c = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
		DB:      5,
	})

	err = c.Set("foo", "baz", 0).Err()
	ok(t, err)

	// Direct access.
	got, err := s.Get("foo")
	ok(t, err)
	equals(t, "bar", got)
	s.Select(5)
	got, err = s.Get("foo")
	ok(t, err)
	equals(t, "baz", got)

	// Another connection should have its own idea of the db:
	c2 := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	v, err := c2.Get("foo").Result()
	ok(t, err)
	equals(t, "bar", v)
}

func TestQuit(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	cmd := redis.NewStringCmd("QUIT")
	err = c.Process(cmd)
	ok(t, err)
	ok(t, cmd.Err())
	equals(t, "OK", cmd.Val())

	res := c.Ping()
	assert(t, res.Err() != nil, "QUIT closed the client")
	equals(t, "", res.Val())
}
