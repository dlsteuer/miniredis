package server

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/go-redis/redis"
)

const (
	errWrongNumberOfArgs = "ERR Wrong number of args"
)

func Test(t *testing.T) {
	s, err := NewServer(":0")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if have := s.Addr().Port; have <= 0 {
		t.Fatalf("have %v, want > 0", have)
	}

	s.Register("PING", func(c *Peer, cmd string, args []string) {
		c.WriteInline("PONG")
	})
	s.Register("ECHO", func(c *Peer, cmd string, args []string) {
		if len(args) != 1 {
			c.WriteError(errWrongNumberOfArgs)
			return
		}
		c.WriteBulk(args[0])
	})
	s.Register("dWaRfS", func(c *Peer, cmd string, args []string) {
		if len(args) != 0 {
			c.WriteError(errWrongNumberOfArgs)
			return
		}
		c.WriteLen(7)
		c.WriteBulk("Blick")
		c.WriteBulk("Flick")
		c.WriteBulk("Glick")
		c.WriteBulk("Plick")
		c.WriteBulk("Quee")
		c.WriteBulk("Snick")
		c.WriteBulk("Whick")
	})
	s.Register("PLUS", func(c *Peer, cmd string, args []string) {
		if len(args) != 2 {
			c.WriteError(errWrongNumberOfArgs)
			return
		}
		a, err := strconv.Atoi(args[0])
		if err != nil {
			c.WriteError(fmt.Sprintf("ERR not an int: %q", args[0]))
			return
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			c.WriteError(fmt.Sprintf("ERR not an int: %q", args[1]))
			return
		}
		c.WriteInt(a + b)
	})
	s.Register("NULL", func(c *Peer, cmd string, args []string) {
		c.WriteNull()
	})

	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr().String(),
	})
	if err != nil {
		t.Fatal(err)
	}

	{
		res, err := c.Ping().Result()
		if err != nil {
			t.Fatal(err)
		}
		if have, want := res, "PONG"; have != want {
			t.Errorf("have: %s, want: %s", have, want)
		}
	}

	{
		res, err := c.Ping().Result()
		if err != nil {
			t.Fatal(err)
		}
		if have, want := res, "PONG"; have != want {
			t.Errorf("have: %s, want: %s", have, want)
		}
	}

	{
		echo, err := c.Echo("hello\nworld").Result()
		if err != nil {
			t.Fatal(err)
		}
		if have, want := echo, "hello\nworld"; have != want {
			t.Errorf("have: %s, want: %s", have, want)
		}
	}

	{
		cmd := redis.NewStringSliceCmd("dwaRFS")
		err := c.Process(cmd)
		if err != nil {
			t.Fatal(err)
		}
		dwarfs, err := cmd.Result()
		if err != nil {
			t.Fatal(err)
		}
		if have, want := dwarfs, []string{"Blick",
			"Flick",
			"Glick",
			"Plick",
			"Quee",
			"Snick",
			"Whick",
		}; !reflect.DeepEqual(have, want) {
			t.Errorf("have: %s, want: %s", have, want)
		}
	}

	{
		bigPayload := strings.Repeat("X", 1<<24)
		echo, err := c.Echo(bigPayload).Result()
		if err != nil {
			t.Fatal(err)
		}
		if have, want := echo, bigPayload; have != want {
			t.Errorf("have: %s, want: %s", have, want)
		}
	}
}
