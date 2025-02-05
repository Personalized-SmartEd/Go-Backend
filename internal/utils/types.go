package utils

import "github.com/golang-jwt/jwt"

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

type ChatRecordInput struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

type RecommentdationInput struct {
	LearningStyle      string   `json:"learning_style"`
	CurrentLevel       string   `json:"current_level"`
	WeakAreas          []string `json:"weak_areas"`
	PerformanceHistory []int    `json:"performance_history"`
	PreferredPace      string   `json:"preferred_pace"`
	AvailableHours     int      `json:"available_hours"`
}

type QuizInput struct {
	SubjectInfo struct {
		Subject          string `json:"subject"`
		Chapter          string `json:"chapter"`
		TopicDescription string `json:"topic_description"`
	} `json:"subject_info"`
	QuizInfo struct {
		QuizDifficultyFrom1To10 int `json:"quiz_difficulty_from_1_to_10"`
		QuizDurationMinutes     int `json:"quiz_duration_minutes"`
		NumberOfQuestions       int `json:"number_of_questions"`
	} `json:"quiz_info"`
}

type DoubtBotInput struct {
	Student struct {
		StudentClass                 int    `json:"student_class"`
		StudentPerformanceFrom1To100 int    `json:"student_performance_from_1_to_100"`
		StudentLearningStyle         string `json:"student_learning_style"`
		StudentPerformanceLevel      string `json:"student_performance_level"`
		StudyPace                    string `json:"study_pace"`
	} `json:"student"`
	Doubt struct {
		Question         string `json:"question"`
		ImageURL         string `json:"image_url"`
		ImageDescription string `json:"image_description"`
	} `json:"doubt"`
	Subject string `json:"subject"`
}

type StaticAssessmentInput struct {
	Responses []int `json:"responses"`
}

type DynamicAssessmentInput struct {
	Subject string `json:"subject"`
	Scores  []int  `json:"scores"`
}

type LearningStyle struct {
	Style string `json:"style"`
}

type StudentLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type TeacherLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type StudentSignUp struct {
	Name        string   `json:"name" validate:"required"`
	Age         int      `json:"age" validate:"required,gte=5"`
	Password    string   `json:"password" validate:"required,min=8"`
	Email       string   `json:"email" validate:"required,email"`
	Image       string   `json:"image"`
	SchoolName  string   `bson:"school_name" validate:"required"`
	SchoolCode  string   `bson:"school_code" validate:"required"`
	Subjects    []string `json:"subjects" validate:"required"`
	ClassNumber int      `json:"class_number" validate:"required,oneof=6 7 8"`
}

type InputClassroom struct {
	TeacherID   string   `bson:"teacher_id" validate:"required"`
	SchoolCode  string   `bson:"school_code" validate:"required"`
	Students    []string `bson:"students"`
	ClassNumber string   `bson:"class_number" validate:"required"`
	ClassCode   string   `bson:"class_code" validate:"required"`
}

type SignedDetailsStudent struct {
	StudentID string
	Name      string
	Email     string
	Class     string
	jwt.StandardClaims
}

type SignedDetailsTeacher struct {
	TeacherID  string
	Name       string
	Email      string
	SchoolCode string
	jwt.StandardClaims
}
