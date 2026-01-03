package questions

type PostQuestionBody struct {
	ProductID uint   `json:"product_id"`
	Content   string `json:"content"`
}

type PutQuestionBody struct {
	Answer string `json:"answer"`
}
