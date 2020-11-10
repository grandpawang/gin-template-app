package redis

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleScript() {
	c, err := Dial("tcp", ":6379")
	if err != nil {
		// handle error
	}
	defer c.Close()
	// Initialize a package-level variable with a script.
	var getScript = NewScript(1, `return call('get', KEYS[1])`)

	// In a function, use the script Do method to evaluate the script. The Do
	// method optimistically uses the EVALSHA command. If the script is not
	// loaded, then the Do method falls back to the EVAL command.
	if _, err = getScript.Do(c, "foo"); err != nil {
		// handle error
	}
}

func TestScript(t *testing.T) {
	c, err := DialDefaultServer()
	if err != nil {
		t.Fatalf("error connection to database, %v", err)
	}
	defer c.Close()

	// To test fall back in Do, we make script unique by adding comment with current time.
	script := fmt.Sprintf("--%d\nreturn {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}", time.Now().UnixNano())
	s := NewScript(2, script)
	reply := []interface{}{[]byte("key1"), []byte("key2"), []byte("arg1"), []byte("arg2")}

	v, err := s.Do(c, "key1", "key2", "arg1", "arg2")
	if err != nil {
		t.Errorf("s.Do(c, ...) returned %v", err)
	}

	if !reflect.DeepEqual(v, reply) {
		t.Errorf("s.Do(c, ..); = %v, want %v", v, reply)
	}

	err = s.Load(c)
	if err != nil {
		t.Errorf("s.Load(c) returned %v", err)
	}

	err = s.SendHash(c, "key1", "key2", "arg1", "arg2")
	if err != nil {
		t.Errorf("s.SendHash(c, ...) returned %v", err)
	}

	err = c.Flush()
	if err != nil {
		t.Errorf("c.Flush() returned %v", err)
	}

	v, err = c.Receive()
	if !reflect.DeepEqual(v, reply) {
		t.Errorf("s.SendHash(c, ..); c.Receive() = %v, want %v", v, reply)
	}

	err = s.Send(c, "key1", "key2", "arg1", "arg2")
	if err != nil {
		t.Errorf("s.Send(c, ...) returned %v", err)
	}

	err = c.Flush()
	if err != nil {
		t.Errorf("c.Flush() returned %v", err)
	}

	v, err = c.Receive()
	if !reflect.DeepEqual(v, reply) {
		t.Errorf("s.Send(c, ..); c.Receive() = %v, want %v", v, reply)
	}

}