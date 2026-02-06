package qqwry

import (
	"fmt"
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	pwd, _ := os.Getwd()
	if err := Download(fmt.Sprintf("%s%s", pwd, "/qqwry.dat")); err != nil {
		t.Fatal(err)
	}
	t.Logf("download success")
}

func BenchmarkFind_DifferentIPs(b *testing.B) {
	datFile := "./qqwry.dat"
	qqWry, err := NewQQWry(datFile)
	if err != nil {
		b.Fatal(err)
	}

	benchCases := []struct {
		name string
		ip   string
	}{
		{"Local", "127.0.0.1"},
		{"China_DNS", "114.114.114.114"},
		{"Ali_DNS", "223.5.5.5"},
		{"Google", "8.8.8.8"},
		{"Cloudflare", "1.1.1.1"},
		{"Foreign_Cloud", "199.21.172.33"},
		{"Broadcast", "255.255.255.255"},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ip, country, area := qqWry.Find(bc.ip)
				b.Logf("ip: %s, country: %s, area: %s", ip, country, area)
			}
		})
	}
}
