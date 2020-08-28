package domain

type InvalidLeadtimeError struct{}

func (e InvalidLeadtimeError) Error() string {
	return "Thời gian không hợp lệ, cần lớn hơn 0"
}
