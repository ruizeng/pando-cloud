package mqtt

import (
	"testing"
)

func TestWildValid(t *testing.T) {
	topic := "/a/b/c"
	if isWildValid(topic) {
		t.Errorf("isWildValid judge error for: %s\n", topic)
	}

	topic = "/a/+"
	if !isWildValid(topic) {
		t.Errorf("isWildValid judge error for: %s\n", topic)
	}

	topic = "/a/+/b/c"
	if !isWildValid(topic) {
		t.Errorf("isWildValid judge error for: %s\n", topic)
	}

	topic = "/a/#/b/c"
	if isWildValid(topic) {
		t.Errorf("isWildValid judge error for: %s\n", topic)
	}

	topic = "/a/#"
	if !isWildValid(topic) {
		t.Errorf("isWildValid judge error for: %s\n", topic)
	}
}

func TestWildExist(t *testing.T) {
	topic := "/a/b/+/c"
	id := "itachili"

	w := NewWild(topic, id)
	
	if !w.isExist(topic, id) {
		t.Errorf("wild exist judeg error, topic: %s, id: %s\n", topic, id)
	}

	id = "uchiha"
	if w.isExist(topic, id) {
		t.Errorf("wild exist judeg error, topic: %s, id: %s\n", topic, id)
	}
}

func TestWildMatch(t *testing.T) {
	topic := "/a/b/+/c"
	id := "itachili"

	w := NewWild(topic, id)
	tpc := "/a/b/d/c"
	if !w.matches(tpc) {
		t.Errorf("wild match judeg error, topic: %s, id: %s, tpc:%s\n", topic, id, tpc)
	}

	tpc = "/a/b/d/d"
	if w.matches(tpc) {
		t.Errorf("wild match judeg error, topic: %s, id: %s, tpc:%s\n", topic, id, tpc)
	}

	topic = "/a/#"
	w = NewWild(topic, id)

	tpc = "/a/b/d/d"
	if !w.matches(tpc) {
		t.Errorf("wild match judeg error, topic: %s, id: %s, tpc:%s\n", topic, id, tpc)
	}

	tpc = "/b/b/d/d"
	if w.matches(tpc) {
		t.Errorf("wild match judeg error, topic: %s, id: %s, tpc:%s\n", topic, id, tpc)
	}
}
