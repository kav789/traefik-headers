package traefik_headers_test

import (
	"io"
	"context"
	"encoding/json"
	headersrwr "github.com/kav789/traefik-headers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"fmt"
)

type testdata struct {
	head  map[string][]string
}

func Test_Headers2(t *testing.T) {

	cases := []struct {
		name  string
		conf  string
		tests []testdata
	}{
		{
			name: "t1",
			conf: `{
  "Content-Security-Policy": [
   "connect-src *; frame-ancestors wildberries.ru *.wildberries.ru wildberries.am *.wildberries.am wildberries.kg *.wildberries.kg wildberries.by *.wildberries.by wildberries.kz *.wildberries.kz wildberries.ua *.wildberries.ua wildberries.eu *.wildberries.eu wildberries.ge *.wildberries.ge"
  ]

}`,
			tests: []testdata{
				testdata{
					head: map[string][]string{
						"cc-bb": []string{ "asdfgh" },
					},
				},
			},
		},
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("nnnnnnnnnnnnnn")
		io.WriteString(rw, "<html><body>Hello World!</body></html>")

	})

	cfg := headersrwr.CreateConfig()
	cfg.HeadersData = `{}`
	_, err := headersrwr.New(context.Background(), next, cfg, "headers")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var tst interface{}
			if err := json.Unmarshal([]byte(tc.conf), &tst); err != nil {
				t.Fatal("init json:", err)
			}
			cfg.HeadersData = tc.conf
			h, err := headersrwr.New(context.Background(), next, cfg, "headers")
			if err != nil {
				t.Fatal(err)
			}

			for _, d := range tc.tests {
				req, err := prepreq("http://aa.vv", d.head)
				if err != nil {
					panic(err)
				}

//				handler := func(w http.ResponseWriter, r *http.Request) {
//					io.WriteString(w, "<html><body>Hello World!</body></html>")
//				}

				rec := httptest.NewRecorder()
//				handler(rec, req)

				h.ServeHTTP(rec, req)
				for k,v := range rec.HeaderMap {
					fmt.Println("hhhhhh", k, v)
				}
/*
				if rec.Code != 200 {
					t.Errorf("first %s %v expected 200 but get %d", d.uri, d.head, rec.Code)
				}
*/
			}
		})
	}

}

func prepreq(uri string, head map[string][]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	if head != nil {
		for k, vv := range head {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}
	return req, nil
}
