package token

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/XiaoMengXinX/FuckCoolapkTokenV3/bcrypt"
	"strconv"
	"strings"
)

// GetToken generate token with X-App-Device, X-APP-Code and timestamp
func GetToken(deviceID string, appVersion int, timeStamp int64) (string, error) {
	const sourceString = "TTBUOFQsQ0ElLUMkWDEjQSEsUyEmLDNAUy1DPFYuIy0iMUM5IzBUJFctMykhMUQ1JS1DOSUxMzhQLVMwV00sIyhSMTQsWTFDRFEwQ0RVMSQkUS0jKSMtJDhYLVQkUjEzKSQsM0BVLTQwUi00NSMsM2BXLVMlJC1DNFlNLCQ0WTEzLSQxQ0EhLEMlIjEzKFcwMz0lLiNEVzEjPSYuNDkkLFMwVixTKFktQzhXLFMsWC4jLSYtU2BTTTA0LFEuMzhTLjMhIjFDMFMsI0BVLCQkUCwzJSIsIzBSMEM8UCxELSItI2BQLUQoUy0jLFMsIyUkMFQwVk0tRCxYMTMlJC40NSEwU0RQLFMhIi0kNSMwU2BQLjNgUywjOFEsNCxTLEMwWCxULSEtIyRTLDNBIy00JFZNMUMsUDBEJFgtNC0kLDQkWC1DMFMxRCkjLVM0VixTMFgxIy0mLVMxJS0zJFQxM2BQLSQ4WTBDMFMuI0UhTSwzPFgtRDBRMDNEUS0jOFYuJDklMSQ0UzAzNSMwRDkkMTMkVSxUMFQtRCxYMTNFIyxELFYsIzUhLiQ0V00xNCxVLCM0Vy4zJFUxNCxTLFQ4Vy4jYFYxJDUjLUMwUzFDPFIsM0ElLDMsVDA0OSMuNCkmLEQkVixEMSJNLTMoUiwzLFctU0RQLiQwVTFDJSUwUykjMUMhIy0zISUsI0EjMTMwUS1UNFcxRDElLUQtJS4kNFIsIzhZTTBTYFEwQyRYLSMlIS40LFgsVDRZLUQoWTBDQFYtRCxWLjNEVDBDMFktNDUlLFQ4VS1UJFgwMzRWLDQ4UU0wQy0mLUQ4WS1TMSMsUyUiMUNAUSxULFYuI0RULiMsWDFELFMxI2BQLVM0WS4kOFcwM0RXLVMoWC0jKFExLCM0VixTNFAtRDRWLUQlJC1TJFAwNCxgYA"
	const packageName = "com.coolapk.market"

	magic := 4*((int(timeStamp)+appVersion)%100) + 128
	const sourceSize = 930
	offset := sourceSize - magic
	if offset >= 0x80 {
		offset = 128
	}

	subString := sourceString[magic : magic+offset]
	decodeBytes, err := base64.StdEncoding.DecodeString(subString)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode sub_string: %w", err)
	}
	decode := string(decodeBytes)

	md5DevBytes := md5.Sum([]byte(deviceID))
	md5Dev := hex.EncodeToString(md5DevBytes[:])

	var combineBuilder strings.Builder
	combineBuilder.WriteString(packageName)
	combineBuilder.WriteString("&")
	combineBuilder.WriteString(decode)
	combineBuilder.WriteString("&")
	combineBuilder.WriteString(md5Dev)
	combineBuilder.WriteString("&")
	combineBuilder.WriteString(strconv.FormatInt(timeStamp, 10))
	combineBuilder.WriteString("&")
	combineBuilder.WriteString(strconv.Itoa(appVersion))
	combine := combineBuilder.String()

	encode := base64.StdEncoding.EncodeToString([]byte(combine))
	md5EncBytes := md5.Sum([]byte(encode))
	md5Enc := hex.EncodeToString(md5EncBytes[:])
	md5CombineBytes := md5.Sum([]byte(combine))
	md5Combine := hex.EncodeToString(md5CombineBytes[:])

	hexTime := strconv.FormatInt(timeStamp, 16)
	hextimeMd5combine := hexTime + "/" + md5Combine
	b64HextimeMd5combine := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(hextimeMd5combine))

	saltPayload := b64HextimeMd5combine
	if len(saltPayload) > 22 {
		saltPayload = saltPayload[:22]
	}
	cost := 10
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(md5Enc), []byte(saltPayload), cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hashing failed: %w", err)
	}

	hashRunes := []rune(string(hashBytes))
	hashRunes[2] = 'y'
	hash := string(hashRunes)
	b64Hash := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(hash))
	return "v3" + b64Hash, nil
}
