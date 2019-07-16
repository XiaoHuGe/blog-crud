package tag_service

type AddTagServer struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	State     int    `json:"state"`
}

type EditTagServer struct {
	//ID          int    `json:"id"`
	Name        string `json:"name"`
	ModifiedBy string `json:"modified_by"`
	State       int    `json:"state"`
}
