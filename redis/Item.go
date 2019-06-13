/*
   @Time : 2019-05-28 15:22
   @Author : frozenchen
   @File : Item
   @Software: studio
*/
package redis

import "encoding/json"

type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key string

	// Value is the Item's value.
	Value []byte

	// Object is the Item's object for use codec.
	Object interface{}

	// Flags are server-opaque flags whose semantics are entirely
	// up to the app.
	Flags uint32

	// Expiration is the cache expiration time, in seconds: either a relative
	// time from now (up to 1 month), or an absolute Unix epoch time.
	// Zero means the Item has no expiration time.
	Expiration int32

	// Compare and swap ID.
	cas uint64
}

// Reply is the result of Get
type Reply struct {
	err    error
	item   *Item
	closed bool
}

func (r *Reply) Scan(item interface{}) (err error) {
	if r.err != nil {
		return r.err
	}
	if len(r.item.Value) == 0 {
		return ErrNotFound
	}

	switch item.(type) {
	case *[]byte:
		*(item.(*[]byte)) = r.item.Value
	case *string:
		*(item.(*string)) = string(r.item.Value)
	case interface{}:
		err = json.Unmarshal(r.item.Value, item)
	}

	return err
}

// Replies is the result of GetMulti
type Replies struct {
	err       error
	items     map[string]*Item
	usedItems map[string]struct{}
	closed    bool
}

func (r *Replies) Scan(key string, item interface{}) (err error) {
	if r.err != nil {
		err = r.err
		return
	}

	var temp *Item
	var ok bool

	if temp, ok = r.items[key]; !ok {
		err = ErrNotFound
		return
	}

	switch item.(type) {
	case *[]byte:
		*(item.(*[]byte)) = temp.Value
	case *string:
		*(item.(*string)) = string(temp.Value)
	case interface{}:
		err = json.Unmarshal(temp.Value, item)
	}

	return

}
