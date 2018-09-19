package util

func NewBool(v bool) *bool {
	return &v
}

func ToBool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}
