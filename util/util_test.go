package util

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIsEven(t *testing.T) {
	res := IsEven(5)
	if res {
		t.Errorf("isEven(5) = true; wanted false")
	}

	res = IsEven(12)
	if !res {
		t.Errorf("isEven(12) = false; wanted true")
	}
}

func TestGetInteger(t *testing.T) {
	os.Setenv("ENV", "TEST")
	server := httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.String() == "/integers/5" {
				rw.Write([]byte(`{"value":5}`))
			}
		},
	))
	defer server.Close()

	api := NewApi(server.URL)
	ctx := context.Background()
	externalResp, err := api.GetInteger(ctx, 5)
	if err != nil {
		t.Errorf("expected resp.Val to be 5, got err %s", err)
		return
	}
	if externalResp == nil {
		t.Errorf("expected resp.Val to be 5, got no resp")
		return
	}
	if externalResp.Value != 5 {
		t.Errorf("expected resp.Val to be 5, got %d", externalResp.Value)
	}
}
