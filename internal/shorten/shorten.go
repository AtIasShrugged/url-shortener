package shorten

import (
	"slices"
	"strings"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

func Shorten(id uint32) string {
	var (
		nums    []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		nums = append(nums, num%alphabetLen)
		num /= alphabetLen
	}

	slices.Reverse(nums)

	for _, num := range nums {
		builder.WriteString(string(alphabet[num]))
	}

	return builder.String()
}
