package password

import "golang.org/x/crypto/bcrypt"

// Hash แปลงรหัสผ่านเป็น bcrypt hash (bcrypt รับได้ไม่เกิน 72 bytes)
func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
