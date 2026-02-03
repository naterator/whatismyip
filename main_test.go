package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type singleURLTestCase struct {
	name           string
	status         int
	responseBody   string
	wantErr        bool
	expectedResult string
}

func singleURLTests() []singleURLTestCase {
	return []singleURLTestCase{
		{
			name:           "success",
			status:         http.StatusOK,
			responseBody:   "1.2.3.4",
			wantErr:        false,
			expectedResult: "1.2.3.4",
		},
		{
			name:           "success with newline",
			status:         http.StatusOK,
			responseBody:   "1.2.3.4\n",
			wantErr:        false,
			expectedResult: "1.2.3.4",
		},
		{
			name:           "invalid IP body",
			status:         http.StatusOK,
			responseBody:   "not-an-ip",
			wantErr:        true,
			expectedResult: "",
		},
		{
			name:           "error status code",
			status:         http.StatusInternalServerError,
			responseBody:   "internal error",
			wantErr:        true,
			expectedResult: "",
		},
	}
}

func TestGetIp(t *testing.T) {
	for _, tt := range singleURLTests() {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := server.Client()
			ip, err := getIp(client, server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("getIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ip != tt.expectedResult {
				t.Errorf("getIp() = %v, want %v", ip, tt.expectedResult)
			}
		})
	}
}

func TestGetIpManySingleURL(t *testing.T) {
	for _, tt := range singleURLTests() {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := server.Client()
			ip, err := getIpMany(client, server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("getIpMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ip != tt.expectedResult {
				t.Errorf("getIpMany() = %v, want %v", ip, tt.expectedResult)
			}
		})
	}
}

func TestGetIpManyMultipleURLs(t *testing.T) {
	t.Run("first URL succeeds", func(t *testing.T) {
		server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "1.1.1.1")
		}))
		defer server1.Close()

		server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "2.2.2.2")
		}))
		defer server2.Close()

		ip, err := getIpMany(server1.Client(), server1.URL, server2.URL)
		if err != nil {
			t.Fatalf("getIpMany() unexpected error: %v", err)
		}
		if ip != "1.1.1.1" {
			t.Errorf("getIpMany() = %v, want %v", ip, "1.1.1.1")
		}
	})

	t.Run("first URL fails, second succeeds", func(t *testing.T) {
		server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server1.Close()

		server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "2.2.2.2")
		}))
		defer server2.Close()

		ip, err := getIpMany(server1.Client(), server1.URL, server2.URL)
		if err != nil {
			t.Fatalf("getIpMany() unexpected error: %v", err)
		}
		if ip != "2.2.2.2" {
			t.Errorf("getIpMany() = %v, want %v", ip, "2.2.2.2")
		}
	})

	t.Run("first URL returns invalid IP, second succeeds", func(t *testing.T) {
		server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "not-an-ip")
		}))
		defer server1.Close()

		server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "2.2.2.2")
		}))
		defer server2.Close()

		ip, err := getIpMany(server1.Client(), server1.URL, server2.URL)
		if err != nil {
			t.Fatalf("getIpMany() unexpected error: %v", err)
		}
		if ip != "2.2.2.2" {
			t.Errorf("getIpMany() = %v, want %v", ip, "2.2.2.2")
		}
	})

	t.Run("all URLs fail", func(t *testing.T) {
		server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server1.Close()

		server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server2.Close()

		_, err := getIpMany(server1.Client(), server1.URL, server2.URL)
		if err == nil {
			t.Fatal("getIpMany() expected error, got nil")
		}
	})
}
