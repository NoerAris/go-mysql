package config

var constantValue map[string]interface{}

func SetConstantValue(key string, value interface{}) {
	if constantValue == nil {
		constantValue = make(map[string]interface{})
	}
	constantValue[key] = value
}
func GetConstantValue(key string) interface{} {
	if val, ok := constantValue[key]; ok {
		return val
	}
	return nil
}
