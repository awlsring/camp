package util

func PtrString(s string) *string {
	return &s
}

func PtrInt(i int) *int {
	return &i
}

func PtrInt64(i int64) *int64 {
	return &i
}

func UnwrapString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UnwrapInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func UnwrapInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
