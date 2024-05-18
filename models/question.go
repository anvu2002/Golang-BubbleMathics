package models

type Question struct {
	ID            string   `json:"id" bson:"_id,omitempty"`
	QuestionText  string   `json:"question" bson:"question"`
	Options       []string `json:"options" bson:"options"`
	CorrectAnswer int      `json:"correctAnswer" bson:"correctAnswer"`
}
