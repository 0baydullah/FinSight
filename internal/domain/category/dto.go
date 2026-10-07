package category

type CreateCategoryRequestDto struct {
	Name string `json:"name"`
}

type UpdateCategoryRequestDto struct {
	Name string `json:"name"`
}

type CategoryResponseDto struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsSystem  bool   `json:"is_system"`
	CreatedAt string `json:"created_at"`
}
