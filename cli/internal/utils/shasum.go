package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/defenseunicorns/zarf/cli/internal/message"
)

func ValidateSha256Sum(expectedChecksum string, path string) {
	actualChecksum, _ := GetSha256Sum(path)
	if expectedChecksum != actualChecksum {
		message.Fatalf("Invalid or mismatched file checksum for %s.  Expected %s, computed %s", path, expectedChecksum, actualChecksum)
	}
}

// GetSha256Sum returns the computed SHA256 Sum of a given file
func GetSha256Sum(path string) (string, error) {
	var data io.ReadCloser
	var err error

	if IsUrl(path) {
		// Handle download from URL
		message.Warn("This is a remote source. If a published checksum is available you should use that rather than calculating it directly from the remote link.")
		data = Fetch(path)
	} else {
		// Handle local file
		data, err = os.Open(path)
		if err != nil {
			return "", err
		}
	}

	defer data.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, data)

	if err != nil {
		return "", err
	} else {
		computedHash := hex.EncodeToString(hash.Sum(nil))
		return computedHash, nil
	}
}
