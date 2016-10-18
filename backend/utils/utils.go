package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/minio/blake2b-simd"
)

// n - how many bytes of entropy, not length
// humanFriendly - if true, returns base32, otherwise base64url
func GenerateToken(n int, humanFriendly bool) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var token string
	if humanFriendly {
		token = base32.StdEncoding.EncodeToString(b)
	} else {
		token = strings.Replace( // base64url
			strings.Replace(
				base64.StdEncoding.EncodeToString(b), "+", "-", -1,
			), "/", "_", -1,
		)
	}

	return token, nil
}

func Blake2b256(s string) []byte {
	tmp := blake2b.Sum256([]byte(s))
	return tmp[:]
}

func EmailToGravatar(email string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(email))))
}

func TimeToDate(t time.Time) time.Time {
	_, z := t.Zone()
	return t.UTC().Add( // get rid of timezone without changing time
		time.Second * time.Duration(z),
	).Truncate( // round down to days
		time.Hour * 24,
	).Add( // set time to the middle of the day
		time.Hour * 12,
	)
}
