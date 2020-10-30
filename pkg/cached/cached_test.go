package cached_test

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andy2046/gopie/pkg/cached"
)

const noncePrefix = "NONCE"

var (
	truthTeller = func(key string) []byte {
		return []byte(fmt.Sprintf("%s-from-truth", key))
	}

	afterFromTruth = make(chan struct{})
	beforeCAS      = make(chan struct{})
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type (
	fakeCache struct {
		sync.RWMutex
		store map[string][]byte
	}

	fakeLease struct {
		nonce string
	}
)

func newCache() *fakeCache {
	return &fakeCache{
		store: make(map[string][]byte),
	}
}

func newLeaseLessor() *fakeLease {
	return &fakeLease{
		nonce: fmt.Sprintf("%s-%s-%s", noncePrefix, randomString(8), randomString(8)),
	}
}

func TestStaleSetProtection(t *testing.T) {
	k := "key1"
	expected := []byte(k + "-from-truth")
	expected2 := []byte(k + "-from-write")
	c := newCache()
	l := newLeaseLessor()
	cd := cached.New(100*time.Millisecond, c, l, truthTeller)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		<-afterFromTruth // reader already retrive value from truth
		c.Set(k, expected2)
		beforeCAS <- struct{}{}

		values := c.Get(k)
		if !bytes.Equal(values[0], expected2) {
			t.Fatalf("write expected %s got %s", expected2, values[0])
		}

		t.Log("write ok")
		wg.Done()
	}()

	go func() {
		value := cd.Read(k)
		if !bytes.Equal(value, expected) {
			t.Fatalf("read expected %s got %s", expected, value)
		}

		t.Log("read ok")
		wg.Done()
	}()

	wg.Wait()
}

func TestThunderingHerdProtection(t *testing.T) {
	k := "key2"
	expected := []byte(k + "-from-truth")
	expected2 := []byte(k + "-from-write")
	c := newCache()
	l := newLeaseLessor()
	cd := cached.New(100*time.Millisecond, c, l, truthTeller)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		<-afterFromTruth // reader already retrive value from truth
		time.Sleep(1 * time.Millisecond)
		c.Set(k, expected2)
		beforeCAS <- struct{}{}

		values := c.Get(k)
		if !bytes.Equal(values[0], expected2) {
			t.Fatalf("write expected %s got %s", expected2, values[0])
		}

		t.Log("write ok")
		wg.Done()
	}()

	go func() {
		value := cd.Read(k)
		if !bytes.Equal(value, expected) {
			t.Fatalf("read1 expected %s got %s", expected, value)
		}

		t.Log("read1 ok")
		wg.Done()
	}()

	go func() {
		value := cd.Read(k)
		if !bytes.Equal(value, expected) {
			t.Fatalf("read2 expected %s got %s", expected, value)
		}

		t.Log("read2 ok")
		wg.Done()
	}()

	wg.Wait()
}

func (fl *fakeLease) Nonce() string {
	return fl.nonce
}

func (fl *fakeLease) NewLease() cached.Lease {
	return newLeaseLessor()
}

func (fl *fakeLease) FromValue(nonce []byte) (cached.Lease, error) {
	v := string(nonce)
	if strings.HasPrefix(v, noncePrefix) {
		return &fakeLease{nonce: v}, nil
	}
	return nil, errors.New("not a Lease")
}

func (fl *fakeLease) MustFromValue(nonce []byte) cached.Lease {
	l, err := fl.FromValue(nonce)
	if err != nil {
		panic(err)
	}
	return l
}

func (fl *fakeLease) IsLease(value []byte) bool {
	if _, err := fl.FromValue(value); err != nil {
		return false
	}
	return true
}

func (fc *fakeCache) Get(keys ...string) [][]byte {
	values := make([][]byte, len(keys))
	fc.RLock()
	defer fc.RUnlock()

	for i, k := range keys {
		values[i] = fc.store[k]
	}
	return values
}

func (fc *fakeCache) Set(key string, value []byte) error {
	fc.Lock()
	defer fc.Unlock()

	fc.store[key] = value
	return nil
}

func (fc *fakeCache) AtomicAdd(key string, value []byte) bool {
	fc.Lock()
	defer fc.Unlock()

	if _, ok := fc.store[key]; !ok {
		fc.store[key] = value
		return true
	}
	return false
}

func (fc *fakeCache) AtomicCheckAndSet(key string, expectedValue, valueToSet []byte) bool {
	// to simulate cache poisoning with stale value
	afterFromTruth <- struct{}{}
	<-beforeCAS

	fc.Lock()
	defer fc.Unlock()

	if v, ok := fc.store[key]; ok && bytes.Equal(v, expectedValue) {
		fc.store[key] = valueToSet
		return true
	}
	return false
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
