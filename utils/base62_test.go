package utils

import "testing"

func TestBase62(t *testing.T) {
	cases := []struct {
		num     int64
		encoded string
	}{
		{0, "0"},
		{61, "Z"},
		{62, "10"},
		{3844, "100"},
	}
	for _, tc := range cases {
		got, err := Base62Encode(tc.num)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.encoded {
			t.Fatalf("encode %d = %s, want %s", tc.num, got, tc.encoded)
		}
		decoded, err := Base62Decode(tc.encoded)
		if err != nil {
			t.Fatal(err)
		}
		if decoded != tc.num {
			t.Fatalf("decode %s = %d, want %d", tc.encoded, decoded, tc.num)
		}
	}
}

func TestMidURLRoundTrip(t *testing.T) {
	cases := []struct {
		mid int64
		url string
	}{
		{0, "0"},
		{1, "1"},
		{61, "Z"},
		{62, "10"},
		{5254998191509253, "Qn4KL6kCN"},
		{3501715734342932, "z0IDjidNi"},
	}
	for _, tc := range cases {
		url, err := MidToURL(tc.mid)
		if err != nil {
			t.Fatal(err)
		}
		if url != tc.url {
			t.Fatalf("MidToURL(%d) = %s, want %s", tc.mid, url, tc.url)
		}
		mid, err := URLToMid(tc.url)
		if err != nil {
			t.Fatal(err)
		}
		if mid != tc.mid {
			t.Fatalf("URLToMid(%s) = %d, want %d", tc.url, mid, tc.mid)
		}
	}
}

func TestBase62DecodeRejectsInvalidChar(t *testing.T) {
	if _, err := Base62Decode("!"); err == nil {
		t.Fatal("expected invalid character to fail")
	}
}
