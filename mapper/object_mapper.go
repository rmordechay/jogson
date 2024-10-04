package mapper

import "encoding/json"

func CreateObjectMapperFromStr(jsonStr string) Json {
	if jsonStr[0] == '[' {
		var ja JsonArray
		ja.value = []byte(jsonStr)
		_ = json.Unmarshal(ja.value, &ja.object)
		return ja
	} else {
		var jo JsonObject
		jo.value = []byte(jsonStr)
		_ = json.Unmarshal(jo.value, &jo.object)
		return jo
	}
}

func createObjectMapper(jsonAny interface{}) Json {
	var om JsonObject
	bytes, _ := json.Marshal(jsonAny)
	om.value = bytes
	_ = json.Unmarshal(om.value, &om.object)
	return om
}
