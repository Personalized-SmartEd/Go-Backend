package utils

type TutorBotInput struct {
	Subject struct {
		Subject          string `json:"subject"`
		Chapter          string `json:"chapter"`
		TopicDescription string `json:"topic_description"`
	} `json:"subject"`
	Student struct {
		StudentClass                 int    `json:"student_class"`
		StudentPerformanceFrom1To100 int    `json:"student_performance_from_1_to_100"`
		StudentLearningStyle         string `json:"student_learning_style"`
		StudentPerformanceLevel      string `json:"student_performance_level"`
		StudyPace                    string `json:"study_pace"`
	} `json:"student"`
	ChatId     string `json:"chat_id"`
	NewMessage string `json:"new_message"`
}

type ChatRecord struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

type ReturnResponse struct {
	Explanation        string       `json:"explanation"`
	UpdatedChatHistory []ChatRecord `json:"updated_chat_history"`
}
