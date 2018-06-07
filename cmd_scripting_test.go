package miniredis

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestEval(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	var b interface{}

	{
		b, err = c.Eval("return 42", []string{}).Result()
		ok(t, err)
		equals(t, 42, b)
	}

	{
		b, err = c.Eval("return {KEYS[1], ARGV[1]}", []string{"key1"}, "key2").Result()
		ok(t, err)
		equals(t, []interface{}{"key1", "key2"}, b)
	}

	{
		b, err = c.Eval("return {ARGV[1]}", []string{}, "key1").Result()
		ok(t, err)
		equals(t, []interface{}{"key1"}, b)
	}

	// Invalid args
	err = c.Eval("42", []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval("[", []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval("os.Exit(42)", []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	{
		b, err = c.Eval(`return string.gsub("foo", "o", "a")`, []string{}).Result()
		ok(t, err)
		equals(t, "faa", b)
	}
}

func TestEvalCall(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	err = c.Eval("redis.call()", []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval("redis.call({})", []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval("redis.call(1)", []string{}).Err()
	assert(t, err != nil, "no EVAL error")
}

func TestScript(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	var (
		script1sha = "a42059b356c875f0717db19a51f6aaca9ae659ea"
		script2sha = "1fa00e76656cc152ad327c13fe365858fd7be306" // "return 42"
	)

	var v string
	var b []bool

	{
		v, err = c.ScriptLoad("return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}").Result()
		ok(t, err)
		equals(t, script1sha, v)
	}

	{
		v, err = c.ScriptLoad("return 42").Result()
		ok(t, err)
		equals(t, script2sha, v)
	}

	{
		b, err = c.ScriptExists(script1sha, script2sha).Result()
		ok(t, err)
		equals(t, []bool{true, true}, b)
	}

	{
		v, err = c.ScriptFlush().Result()
		ok(t, err)
		equals(t, "OK", v)
	}

	{
		b, err = c.ScriptExists(script1sha).Result()
		ok(t, err)
		equals(t, []bool{false}, b)
	}

	{
		b, err = c.ScriptExists().Result()
		ok(t, err)
		equals(t, []bool{}, b)
	}
}

func TestCJSON(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	test := func(expr, want string) {
		t.Helper()
		str, err := c.Eval(expr, []string{}).Result()
		ok(t, err)
		equals(t, str, want)
	}
	test(
		`return cjson.decode('{"id":"foo"}')['id']`,
		"foo",
	)
	test(
		`return cjson.encode({foo=42})`,
		`{"foo":42}`,
	)

	err = c.Eval("redis.encode()", []string{}).Err()
	assert(t, err != nil, "lua error")
	err = c.Eval(`redis.encode("1", "2")`, []string{}).Err()
	assert(t, err != nil, "lua error")
	err = c.Eval(`redis.decode()`, []string{}).Err()
	assert(t, err != nil, "lua error")
	err = c.Eval(`redis.decode("{")`, []string{}).Err()
	assert(t, err != nil, "lua error")
	err = c.Eval(`redis.decode("1", "2")`, []string{}).Err()
	assert(t, err != nil, "lua error")
}

func TestSha1Hex(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	test1 := func(val interface{}, want string) {
		t.Helper()
		str, err := c.Eval("return redis.sha1hex(ARGV[1])", []string{}, val).Result()
		ok(t, err)
		equals(t, str, want)
	}
	test1("foo", "0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33")
	test1("bar", "62cdb7020ff920e5aa642c3d4066950dd1f01f4d")
	test1("0", "b6589fc6ab0dc82cf12099d1c2d40ab994e8410c")
	test1(0, "b6589fc6ab0dc82cf12099d1c2d40ab994e8410c")
	test1(nil, "da39a3ee5e6b4b0d3255bfef95601890afd80709")

	test2 := func(eval, want string) {
		t.Helper()
		have, err := c.Eval(eval, []string{}).Result()
		ok(t, err)
		equals(t, have, want)
	}
	test2("return redis.sha1hex({})", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
	test2("return redis.sha1hex(nil)", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
	test2("return redis.sha1hex(42)", "92cfceb39d57d914ed8b14d0e37643de0797ae56")

	err = c.Eval("redis.sha1hex()", []string{}).Err()
	assert(t, err != nil, "lua error")
}

func TestEvalsha(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	var v string
	var b interface{}

	script1sha := "d006f1a90249474274c76f5be725b8f5804a346b"
	{
		v, err = c.ScriptLoad("return {KEYS[1], ARGV[1]}").Result()
		ok(t, err)
		equals(t, script1sha, v)
	}

	{
		b, err = c.EvalSha(script1sha, []string{"key1"}, "key2").Result()
		ok(t, err)
		equals(t, []interface{}{"key1", "key2"}, b)
	}

	err = c.EvalSha("foo", []string{}).Err()
	mustFail(t, err, msgNoScriptFound)

	err = c.EvalSha("foo", []string{"bar"}).Err()
	mustFail(t, err, msgNoScriptFound)
}

func TestCmdEvalReply(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	test := func(script string, keyCount int, args []interface{}, expected interface{}) {
		t.Helper()
		keys := []string{}
		for i := 0; i < keyCount; i++ {
			keys = append(keys, args[i].(string))
		}
		reply, err := c.Eval(script, keys, args[keyCount:]...).Result()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		equals(t, expected, reply)
	}

	// return boolean true
	test(
		"return true",
		0,
		[]interface{}{},
		int64(1),
	)
	// return boolean false
	test(
		"return false",
		0,
		[]interface{}{},
		int64(0),
	)
	// return single number
	test(
		"return 10",
		0,
		[]interface{}{},
		int64(10),
	)
	// return single float
	test(
		"return 12.345",
		0,
		[]interface{}{},
		int64(12),
	)
	// return multiple numbers
	test(
		"return 10, 20",
		0,
		[]interface{}{},
		int64(10),
	)
	// return single string
	test(
		"return 'test'",
		0,
		[]interface{}{},
		[]byte("test"),
	)
	// return multiple string
	test(
		"return 'test1', 'test2'",
		0,
		[]interface{}{},
		[]byte("test1"),
	)
	// return single table multiple integer
	test(
		"return {10, 20}",
		0,
		[]interface{}{},
		[]interface{}{
			int64(10),
			int64(20),
		},
	)
	// return single table multiple string
	test(
		"return {'test1', 'test2'}",
		0,
		[]interface{}{},
		[]interface{}{
			"test1",
			"test2",
		},
	)
	// return nested table
	test(
		"return {10, 20, {30, 40}}",
		0,
		[]interface{}{},
		[]interface{}{
			int64(10),
			int64(20),
			[]interface{}{
				int64(30),
				int64(40),
			},
		},
	)
	// return combination table
	test(
		"return {10, 20, {30, 'test', true, 40}, false}",
		0,
		[]interface{}{},
		[]interface{}{
			int64(10),
			int64(20),
			[]interface{}{
				int64(30),
				"test",
				int64(1),
				int64(40),
			},
			int64(0),
		},
	)
	// KEYS and ARGV
	test(
		"return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}",
		2,
		[]interface{}{
			"key1",
			"key2",
			"first",
			"second",
		},
		[]interface{}{
			"key1",
			"key2",
			"first",
			"second",
		},
	)

	{
		err = c.Eval(`return {err="broken"}`, []string{}).Err()
		mustFail(t, err, "broken")

		err = c.Eval(`return redis.error_reply("broken")`, []string{}).Err()
		mustFail(t, err, "broken")
	}

	var v interface{}

	{
		v, err = c.Eval(`return {ok="good"}`, []string{}).Result()
		ok(t, err)
		equals(t, "good", v)

		v, err = c.Eval(`return redis.status_reply("good")`, []string{}).Result()
		ok(t, err)
		equals(t, "good", v)
	}

	err = c.Eval(`return redis.error_reply()`, []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval(`return redis.error_reply(1)`, []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval(`return redis.status_reply()`, []string{}).Err()
	assert(t, err != nil, "no EVAL error")

	err = c.Eval(`return redis.status_reply(1)`, []string{}).Err()
	assert(t, err != nil, "no EVAL error")
}

func TestCmdEvalResponse(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()

	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})
	defer c.Close()

	var v interface{}

	{
		v, err = c.Eval("return redis.call('set', 'foo', 'bar')", []string{}).Result()
		ok(t, err)
		equals(t, "OK", v)
	}

	{
		v, err = c.Eval("return redis.call('get', 'foo')", []string{}).Result()
		ok(t, err)
		equals(t, "bar", v)
	}

	{
		v, err = c.Eval("return redis.call('HMSET', 'mkey', 'foo', 'bar', 'foo1', 'bar1')", []string{}).Result()
		ok(t, err)
		equals(t, "OK", v)
	}

	{
		v, err = c.Eval("return redis.call('HGETALL', 'mkey')", []string{}).Result()
		ok(t, err)
		equals(t, []interface{}{"foo", "bar", "foo1", "bar1"}, v)
	}

	{
		v, err = c.Eval("return redis.call('HMGET', 'mkey', 'foo1')", []string{}).Result()
		ok(t, err)
		equals(t, []interface{}{"bar1"}, v)
	}

	{
		v, err = c.Eval("return redis.call('HMGET', 'mkey', 'foo')", []string{}).Result()
		ok(t, err)
		equals(t, []interface{}{"bar"}, v)
	}

	{
		v, err = c.Eval("return redis.call('HMGET', 'mkey', 'bad', 'key')", []string{}).Result()
		ok(t, err)
		equals(t, []interface{}{nil, nil}, v)
	}
}

func TestCmdEvalAuth(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()

	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	eval := "return redis.call('set','foo','bar')"

	s.RequireAuth("123password")

	err = c.Eval(eval, []string{}).Err()
	mustFail(t, err, "NOAUTH Authentication required.")

	c.Close()
	c = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
		Password: "123password",
	})

	err = c.Eval(eval, []string{}).Err()
	ok(t, err)
}
