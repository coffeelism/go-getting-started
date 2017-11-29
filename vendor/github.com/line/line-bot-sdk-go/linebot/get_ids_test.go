// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package linebot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGetGroupMemberIDs(t *testing.T) {
	type want struct {
		URLPath           string
		ContinuationToken string
		RequestBody       []byte
		Response          *MemberIDsResponse
		Error             error
	}
	var testCases = []struct {
		GroupID           string
		ContinuationToken string
		ResponseCode      int
		Response          []byte
		Want              want
	}{
		{
			GroupID:           "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ContinuationToken: "xxxxx",
			ResponseCode:      200,
			Response:          []byte(`{"memberIds": ["U0047556f2e40dba2456887320ba7c76d", "U0047556f2e40dba2456887320ba7c76e"]}`),
			Want: want{
				URLPath:           fmt.Sprintf(APIEndpointGetGroupMemberIDs, "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				ContinuationToken: "xxxxx",
				RequestBody:       []byte(""),
				Response: &MemberIDsResponse{
					MemberIDs: []string{
						"U0047556f2e40dba2456887320ba7c76d",
						"U0047556f2e40dba2456887320ba7c76e",
					},
				},
			},
		},
		{
			GroupID:      "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ResponseCode: 200,
			Response:     []byte(`{"memberIds": ["U0047556f2e40dba2456887320ba7c76d", "U0047556f2e40dba2456887320ba7c76e"], "next": "xxxxx"}`),
			Want: want{
				URLPath:     fmt.Sprintf(APIEndpointGetGroupMemberIDs, "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				RequestBody: []byte(""),
				Response: &MemberIDsResponse{
					MemberIDs: []string{
						"U0047556f2e40dba2456887320ba7c76d",
						"U0047556f2e40dba2456887320ba7c76e",
					},
					Next: "xxxxx",
				},
			},
		},
		{
			// Internal server error
			GroupID:           "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ContinuationToken: "xxxxx",
			ResponseCode:      500,
			Response:          []byte("500 Internal server error"),
			Want: want{
				URLPath:           fmt.Sprintf(APIEndpointGetGroupMemberIDs, "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				ContinuationToken: "xxxxx",
				RequestBody:       []byte(""),
				Error: &APIError{
					Code: 500,
				},
			},
		},
	}

	var currentTestIdx int
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tc := testCases[currentTestIdx]
		if r.Method != http.MethodGet {
			t.Errorf("Method %s; want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != tc.Want.URLPath {
			t.Errorf("URLPath %s; want %s", r.URL.Path, tc.Want.URLPath)
		}
		q := r.URL.Query()
		if start, want := q.Get("start"), tc.Want.ContinuationToken; start != want {
			t.Errorf("ContinuationToken: %s; want %s", start, want)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(body, tc.Want.RequestBody) {
			t.Errorf("RequestBody %s; want %s", body, tc.Want.RequestBody)
		}
		w.WriteHeader(tc.ResponseCode)
		w.Write(tc.Response)
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Fatal(err)
	}
	for i, tc := range testCases {
		currentTestIdx = i
		res, err := client.GetGroupMemberIDs(tc.GroupID, tc.ContinuationToken).Do()
		if tc.Want.Error != nil {
			if !reflect.DeepEqual(err, tc.Want.Error) {
				t.Errorf("Error %d %q; want %q", i, err, tc.Want.Error)
			}
		} else {
			if err != nil {
				t.Error(err)
			}
		}
		if !reflect.DeepEqual(res, tc.Want.Response) {
			t.Errorf("Response %d %q; want %q", i, res, tc.Want.Response)
		}
	}
}

func TestGetGroupMemberIDsWithContext(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("{}"))
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	_, err = client.GetGroupMemberIDs("cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "xxxxx").WithContext(ctx).Do()
	if err != context.DeadlineExceeded {
		t.Errorf("err %v; want %v", err, context.DeadlineExceeded)
	}
}

func BenchmarkGetGroupMemberIDs(b *testing.B) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte(`{"memberIds": ["U0047556f2e40dba2456887320ba7c76d", "U0047556f2e40dba2456887320ba7c76e"], "next": "xxxxx"}`))
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = client.GetGroupMemberIDs("cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "xxxxx").Do()
	}
}

func TestGetRoomMemberIDs(t *testing.T) {
	type want struct {
		URLPath           string
		ContinuationToken string
		RequestBody       []byte
		Response          *MemberIDsResponse
		Error             error
	}
	var testCases = []struct {
		RoomID            string
		ContinuationToken string
		ResponseCode      int
		Response          []byte
		Want              want
	}{
		{
			RoomID:            "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ContinuationToken: "xxxxx",
			ResponseCode:      200,
			Response:          []byte(`{"memberIds": ["U0047556f2e40dba2456887320ba7c76d", "U0047556f2e40dba2456887320ba7c76e"]}`),
			Want: want{
				URLPath:           fmt.Sprintf(APIEndpointGetRoomMemberIDs, "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				ContinuationToken: "xxxxx",
				RequestBody:       []byte(""),
				Response: &MemberIDsResponse{
					MemberIDs: []string{
						"U0047556f2e40dba2456887320ba7c76d",
						"U0047556f2e40dba2456887320ba7c76e",
					},
				},
			},
		},
		{
			RoomID:       "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ResponseCode: 200,
			Response:     []byte(`{"memberIds": ["U0047556f2e40dba2456887320ba7c76d", "U0047556f2e40dba2456887320ba7c76e"], "next": "xxxxx"}`),
			Want: want{
				URLPath:     fmt.Sprintf(APIEndpointGetRoomMemberIDs, "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				RequestBody: []byte(""),
				Response: &MemberIDsResponse{
					MemberIDs: []string{
						"U0047556f2e40dba2456887320ba7c76d",
						"U0047556f2e40dba2456887320ba7c76e",
					},
					Next: "xxxxx",
				},
			},
		},
		{
			// Internal server error
			RoomID:            "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ContinuationToken: "xxxxx",
			ResponseCode:      500,
			Response:          []byte("500 Internal server error"),
			Want: want{
				URLPath:           fmt.Sprintf(APIEndpointGetRoomMemberIDs, "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				ContinuationToken: "xxxxx",
				RequestBody:       []byte(""),
				Error: &APIError{
					Code: 500,
				},
			},
		},
	}

	var currentTestIdx int
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tc := testCases[currentTestIdx]
		if r.Method != http.MethodGet {
			t.Errorf("Method %s; want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != tc.Want.URLPath {
			t.Errorf("URLPath %s; want %s", r.URL.Path, tc.Want.URLPath)
		}
		q := r.URL.Query()
		if start, want := q.Get("start"), tc.Want.ContinuationToken; start != want {
			t.Errorf("ContinuationToken: %s; want %s", start, want)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(body, tc.Want.RequestBody) {
			t.Errorf("RequestBody %s; want %s", body, tc.Want.RequestBody)
		}
		w.WriteHeader(tc.ResponseCode)
		w.Write(tc.Response)
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Fatal(err)
	}
	for i, tc := range testCases {
		currentTestIdx = i
		res, err := client.GetRoomMemberIDs(tc.RoomID, tc.ContinuationToken).Do()
		if tc.Want.Error != nil {
			if !reflect.DeepEqual(err, tc.Want.Error) {
				t.Errorf("Error %d %q; want %q", i, err, tc.Want.Error)
			}
		} else {
			if err != nil {
				t.Error(err)
			}
		}
		if !reflect.DeepEqual(res, tc.Want.Response) {
			t.Errorf("Response %d %q; want %q", i, res, tc.Want.Response)
		}
	}
}

func TestGetRoomMemberIDsWithContext(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("{}"))
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	_, err = client.GetRoomMemberIDs("cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "xxxxx").WithContext(ctx).Do()
	if err != context.DeadlineExceeded {
		t.Errorf("err %v; want %v", err, context.DeadlineExceeded)
	}
}

func BenchmarkGetRoomMemberIDs(b *testing.B) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte(`{"memberIds": ["U0047556f2e40dba2456887320ba7c76d", "U0047556f2e40dba2456887320ba7c76e"], "next": "xxxxx"}`))
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = client.GetRoomMemberIDs("cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "xxxxx").Do()
	}
}

func TestGetGroupMemberIDsScanner(t *testing.T) {
	res1 := &MemberIDsResponse{
		MemberIDs: []string{},
		Next:      "xxxxx",
	}

	res2 := &MemberIDsResponse{
		MemberIDs: []string{},
	}

	for i := 0; i < 150; i++ {
		if i < 100 {
			res1.MemberIDs = append(res1.MemberIDs, fmt.Sprintf("U%032d", i))
		} else {
			res2.MemberIDs = append(res2.MemberIDs, fmt.Sprintf("U%032d", i))
		}
	}

	under100Res := &MemberIDsResponse{
		MemberIDs: []string{},
	}

	for i := 0; i < 50; i++ {
		under100Res.MemberIDs = append(under100Res.MemberIDs, fmt.Sprintf("U%032d", i))
	}

	testCases := []struct {
		FirstResponse   *MemberIDsResponse
		SecoundResponse *MemberIDsResponse
	}{
		{
			FirstResponse:   res1,
			SecoundResponse: res2,
		},
		{
			FirstResponse: under100Res,
		},
	}

	var currentTestIdx int
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tc := testCases[currentTestIdx]
		if r.Method != http.MethodGet {
			t.Errorf("Method %s; want %s", r.Method, http.MethodGet)
		}
		q := r.URL.Query()
		w.WriteHeader(200)
		if q.Get("start") == res1.Next {
			if err := json.NewEncoder(w).Encode(tc.SecoundResponse); err != nil {
				t.Fatal(err)
			}
			return
		}
		if err := json.NewEncoder(w).Encode(tc.FirstResponse); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		currentTestIdx = i

		s := client.GetGroupMemberIDs("cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "").NewScanner()
		for k := 0; s.Scan(); k++ {
			if id, want := s.ID(), fmt.Sprintf("U%032d", k); id != want {
				t.Fatalf("case[%d] id = %s; want %s; scanner = %#v", currentTestIdx, id, want, s)
			}
		}
		if err := s.Err(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetRoomMemberIDsScanner(t *testing.T) {
	res1 := &MemberIDsResponse{
		MemberIDs: []string{},
		Next:      "xxxxx",
	}

	res2 := &MemberIDsResponse{
		MemberIDs: []string{},
	}

	for i := 0; i < 150; i++ {
		if i < 100 {
			res1.MemberIDs = append(res1.MemberIDs, fmt.Sprintf("U%032d", i))
		} else {
			res2.MemberIDs = append(res2.MemberIDs, fmt.Sprintf("U%032d", i))
		}
	}

	under100Res := &MemberIDsResponse{
		MemberIDs: []string{},
	}

	for i := 0; i < 50; i++ {
		under100Res.MemberIDs = append(under100Res.MemberIDs, fmt.Sprintf("U%032d", i))
	}

	testCases := []struct {
		FirstResponse   *MemberIDsResponse
		SecoundResponse *MemberIDsResponse
	}{
		{
			FirstResponse:   res1,
			SecoundResponse: res2,
		},
		{
			FirstResponse: under100Res,
		},
	}

	var currentTestIdx int
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tc := testCases[currentTestIdx]
		if r.Method != http.MethodGet {
			t.Errorf("Method %s; want %s", r.Method, http.MethodGet)
		}
		q := r.URL.Query()
		w.WriteHeader(200)
		if q.Get("start") == res1.Next {
			if err := json.NewEncoder(w).Encode(tc.SecoundResponse); err != nil {
				t.Fatal(err)
			}
			return
		}
		if err := json.NewEncoder(w).Encode(tc.FirstResponse); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()
	client, err := mockClient(server)
	if err != nil {
		t.Fatal(err)
	}

	for i := range testCases {
		currentTestIdx = i

		s := client.GetRoomMemberIDs("cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "").NewScanner()
		for k := 0; s.Scan(); k++ {
			if id, want := s.ID(), fmt.Sprintf("U%032d", k); id != want {
				t.Fatalf("case[%d] id = %s; want %s; scanner = %#v", currentTestIdx, id, want, s)
			}
		}
		if err := s.Err(); err != nil {
			t.Fatal(err)
		}
	}
}
