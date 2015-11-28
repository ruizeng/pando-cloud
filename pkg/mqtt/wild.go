package mqtt

import (
	"reflect"
	"strings"
)

type Wild struct {
	wild     []string
	identify string
}

func NewWild(topic string, id string) *Wild {
	return &Wild{
		wild:     strings.Split(topic, "/"),
		identify: id,
	}
}

func isWild(topic string) bool {
	if strings.Contains(topic, "#") || strings.Contains(topic, "+") {
		return true
	}

	return false
}

func isWildValid(topic string) bool {
	if !isWild(topic) {
		return false
	}

	wilds := strings.Split(topic, "/")
	for i, part := range wilds {
		if isWild(part) && len(part) != 1 {
			return false
		}

		if part == "#" && i != len(wilds)-1 {
			return false
		}
	}

	return true
}

func (w *Wild) isExist(topic string, id string) bool {
	parts := strings.Split(topic, "/")
	if reflect.DeepEqual(parts, w.wild) && id == w.identify {
		return true
	}

	return false
}

func (w *Wild) matches(topic string) bool {
	parts := strings.Split(topic, "/")
	i := 0
	for i = 0; i < len(parts); i++ {
		if i >= len(w.wild) {
			return false
		}

		if w.wild[i] == "#" {
			return true
		}

		if parts[i] != w.wild[i] && w.wild[i] != "+" {
			return false
		}
	}

	if i == len(w.wild)-1 && w.wild[len(w.wild)-1] == "#" {
		return true
	}

	if i == len(w.wild) {
		return true
	}

	return false
}
