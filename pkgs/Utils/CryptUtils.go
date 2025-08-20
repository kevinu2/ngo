package Utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
)

const InvitationKey = "E5aFbCDcGd3HeQAf4Bg1NiOjPhIkJ2lRSnTUmV67MWX89KLYZ"

type crypt struct {
}

func Crypt() crypt {
	return crypt{}
}

func (crypt) Sha1(v string) string {
	sha := sha1.New()
	sha.Write([]byte(v))
	bs := sha.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (crypt) Sha3(v string) string {
	sha := sha3.NewLegacyKeccak512()
	sha.Write([]byte(v))
	bs := sha.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (crypt) Sha512(v string) string {
	sha := sha512.New()
	sha.Write([]byte(v))
	bs := sha.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (crypt) Hmac(v string) string {
	hash := hmac.New(sha1.New, []byte(v))
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (crypt) UUID() string {
	ui := uuid.Must(uuid.NewV4(), nil)
	return ui.String()
}

func (crypt) MD5(v string) string {
	md := md5.New()
	md.Write([]byte(v))
	bs := md.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (crypt) SessionId() string {
	return Crypt().Sha1(Crypt().UUID())
}

func (crypt) EncodeInvitation(userId uint64) string {
	length := uint64(len(InvitationKey))
	num := userId
	var code string
	for {
		if num <= 0 {
			break
		}
		mod := num % length
		num = (num - mod) / length

		code = fmt.Sprintf("%s%s", string(InvitationKey[mod]), code)
	}
	fmt.Println(code)
	if len(code) < 6 {
		code = fmt.Sprintf("%0*s", 6, code)
	}
	return code
}

func (crypt) DecodeInvitation(invitationCode string) int64 {
	length := int64(len(InvitationKey))
	startIndex := strings.LastIndex(invitationCode, "0") + 1
	code := invitationCode[startIndex:len(invitationCode)]
	code = Strings().Reverse(code)
	var num int64
	for i := 0; i < len(code); i++ {
		index := strings.Index(InvitationKey, string(code[i]))
		le := decimal.NewFromInt(length)
		m := decimal.NewFromInt(int64(i))
		ind := decimal.NewFromInt(int64(index))
		nu := decimal.NewFromInt(num)
		num = le.Pow(m).Mul(ind).Add(nu).IntPart()
	}
	return num
}

func (crypt) EncryptPwd(pwd string, salt string) string {
	first := Crypt().Sha3(pwd)
	second := Crypt().Sha3(first + salt)
	return Crypt().MD5(second)
}
