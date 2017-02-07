package sha

import(
  "crypto/sha1"
  "io"
)

func shaone(input string, n int) string {
  hash := sha1.New()
  io.WriteString(hash, input)
  bytes := hash.Sum(nil)
  output := string(bytes[:(n/8)])
  return output
}
