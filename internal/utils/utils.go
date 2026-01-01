package utils

import (
	"crypto/rand"
	"encoding/json"
	"io"
)

func Copy(src, dest any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Hàm tạo OTP chỉ chứa số (0-9)
func GenerateOTP(length int) (string, error) {
	// Bộ ký tự cho phép
	const otpChars = "0123456789"
	buffer := make([]byte, length)

	// Đọc byte ngẫu nhiên từ hệ thống
	_, err := io.ReadAtLeast(rand.Reader, buffer, length)
	if err != nil {
		return "", err
	}

	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		// Map byte ngẫu nhiên vào bộ ký tự otpChars
		otp[i] = otpChars[int(buffer[i])%len(otpChars)]
	}

	return string(otp), nil
}
