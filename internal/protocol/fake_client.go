package protocol

import (
	"fmt"
	"math/rand"
)

// this function is used to generate random phone model information
// note that appVersion is always 1.8.2, which could lead to suspicion!!!
func GenerateFakeClient() ClientInfo {
	brands := []string{
		"Xiaomi", "HUAWEI", "HONOR", "OPPO", "vivo", "OnePlus", "Samsung",
	}

	// Common Android versions
	androidVersions := []string{
		"11.0", "12.0", "13.0", "14.0",
	}

	// Generate random device token (16 characters)
	// const letterBytes = "abcdef0123456789"
	// deviceToken := make([]byte, 16)
	// for i := range deviceToken {
	// 	deviceToken[i] = letterBytes[rand.Intn(len(letterBytes))]
	// }

	// Random model numbers
	modelNumbers := []string{"2201123C", "2207122C", "22081212C", "23046PNC9C"}

	brand := brands[rand.Intn(len(brands))]
	model := modelNumbers[rand.Intn(len(modelNumbers))]

	return ClientInfo{
		AppVersion:  "1.8.2",                                          // Fixed version
		Brand:       brand,                                            // Random brand
		DeviceToken: "",                                               // Random device token
		DeviceType:  fmt.Sprintf("%s_%s", brand, model),               // Brand_Model format
		MobileType:  "android",                                        // Fixed as android
		SysVersion:  androidVersions[rand.Intn(len(androidVersions))], // Random Android version
	}
}
