package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const Base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encode(num int64) (string, error) {
	if num < 0 {
		return "", fmt.Errorf("Base62 编码失败: 不能编码负数 %d", num)
	}
	if num == 0 {
		return Base62Alphabet[:1], nil
	}
	base := int64(len(Base62Alphabet))
	encoded := make([]byte, 0)
	for num > 0 {
		rem := num % base
		num /= base
		encoded = append(encoded, Base62Alphabet[rem])
	}
	reverseBytes(encoded)
	return string(encoded), nil
}

func Base62Decode(s string) (int64, error) {
	base := int64(len(Base62Alphabet))
	var num int64
	for i := 0; i < len(s); i++ {
		idx := strings.IndexByte(Base62Alphabet, s[i])
		if idx < 0 {
			return 0, fmt.Errorf("Base62 解码失败: 非法字符 %q", s[i])
		}
		power := len(s) - (i + 1)
		place := int64(1)
		for range power {
			if place > (1<<63-1)/base {
				return 0, fmt.Errorf("Base62 解码失败: 数值超出 int64 范围")
			}
			place *= base
		}
		value := int64(idx)
		if value > 0 && place > (1<<63-1)/value {
			return 0, fmt.Errorf("Base62 解码失败: 数值超出 int64 范围")
		}
		if num > (1<<63-1)-value*place {
			return 0, fmt.Errorf("Base62 解码失败: 数值超出 int64 范围")
		}
		num += value * place
	}
	return num, nil
}

func URLToMid(url string) (int64, error) {
	if url == "" {
		return 0, fmt.Errorf("短链转 mid 失败: 短链为空")
	}
	reversed := reverseString(url)
	size := ceilDiv(len(reversed), 4)
	parts := make([]string, 0, size)
	for i := range size {
		start := i * 4
		end := min(start+4, len(reversed))
		chunk := reverseString(reversed[start:end])
		n, err := Base62Decode(chunk)
		if err != nil {
			return 0, fmt.Errorf("短链转 mid 失败: %w", err)
		}
		part := strconv.FormatInt(n, 10)
		if i < size-1 && len(part) < 7 {
			part = strings.Repeat("0", 7-len(part)) + part
		}
		parts = append(parts, part)
	}
	reverseStrings(parts)
	mid, err := strconv.ParseInt(strings.Join(parts, ""), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("短链转 mid 失败: %w", err)
	}
	return mid, nil
}

func MidToURL(mid int64) (string, error) {
	if mid < 0 {
		return "", fmt.Errorf("mid 转短链失败: mid 不能为负数 %d", mid)
	}
	reversed := reverseString(strconv.FormatInt(mid, 10))
	size := ceilDiv(len(reversed), 7)
	parts := make([]string, 0, size)
	for i := range size {
		start := i * 7
		end := min(start+7, len(reversed))
		chunk := reverseString(reversed[start:end])
		n, err := strconv.ParseInt(chunk, 10, 64)
		if err != nil {
			return "", fmt.Errorf("mid 转短链失败: %w", err)
		}
		part, err := Base62Encode(n)
		if err != nil {
			return "", fmt.Errorf("mid 转短链失败: %w", err)
		}
		if i < size-1 && len(part) < 4 {
			part = strings.Repeat("0", 4-len(part)) + part
		}
		parts = append(parts, part)
	}
	reverseStrings(parts)
	return strings.Join(parts, ""), nil
}

func ceilDiv(n, div int) int {
	size := n / div
	if n%div != 0 {
		size++
	}
	return size
}

func reverseString(s string) string {
	b := []byte(s)
	reverseBytes(b)
	return string(b)
}

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func reverseStrings(parts []string) {
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
}
