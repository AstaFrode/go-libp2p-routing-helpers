package routinghelpers

import (
	"context"
	"testing"

	"github.com/AstaFrode/go-libp2p/core/routing"
)

func TestLimitedValueStore(t *testing.T) {
	d := LimitedValueStore{
		ValueStore: new(dummyValueStore),
		Namespaces: []string{"allow"},
	}

	ctx := context.Background()

	for i, k := range []string{
		"/allow/hello",
		"/allow/foo",
		"/allow/foo/bar",
	} {
		if err := d.PutValue(ctx, k, []byte{byte(i)}); err != nil {
			t.Fatal(err)
		}
		v, err := d.GetValue(ctx, k)
		if err != nil {
			t.Fatal(err)
		}
		if len(v) != 1 || v[0] != byte(i) {
			t.Fatalf("expected value [%d], got %v", i, v)
		}
	}
	for i, k := range []string{
		"/deny/hello",
		"/allow",
		"allow",
		"deny",
		"",
		"/",
		"//",
		"///",
		"//allow",
	} {
		if err := d.PutValue(ctx, k, []byte{byte(i)}); err != routing.ErrNotSupported {
			t.Fatalf("expected put with key %s to fail", k)
		}
		_, err := d.GetValue(ctx, k)
		if err != routing.ErrNotFound {
			t.Fatalf("expected get with key %s to fail", k)
		}
		_, err = d.ValueStore.GetValue(ctx, k)
		if err != routing.ErrNotFound {
			t.Fatalf("expected get with key %s to fail", k)
		}
		err = d.ValueStore.PutValue(ctx, k, []byte{byte(i)})
		if err != nil {
			t.Fatal(err)
		}
		_, err = d.GetValue(ctx, k)
		if err == nil {
			t.Fatalf("expected get with key %s to fail", k)
		}
	}
}

func TestLimitedClose(t *testing.T) {
	closer := new(testCloser)
	d := LimitedValueStore{
		ValueStore: struct {
			*testCloser
			routing.Routing
		}{closer, Null{}},
	}
	d.Close()
	if closer.closed != 1 {
		t.Fatal("expected one close")
	}
}
