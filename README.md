# Go-Backend

This repository contains the backend code for the Personalized-SmartEd project, written in Go.

## Introduction

The Go-Backend project is designed to provide a robust and scalable backend for the Personalized-SmartEd platform. It leverages the Go programming language to deliver high performance and concurrent processing for educational services.

## Features

- High performance backend services
- Scalable and maintainable codebase
- RESTful API endpoints
- Integration with various data sources

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Project Substructure](#project-substructure)
- [Interaction with ML-Service](#interaction-with-ml-service)
- [Contributing](#contributing)

## Installation

To get started with the Go-Backend, follow these steps:

1. Clone the repository:
   ```sh
   git clone https://github.com/Personalized-SmartEd/Go-Backend.git
   ```
2. Navigate to the project directory:
   ```sh
   cd Go-Backend
   ```
3. Install the necessary dependencies:
   ```sh
   go mod tidy
   ```

## Usage

Run the application with the following command:
```sh
go run main.go
```

## Project Substructure

Below is the folder structure of the **Go-Backend** repository:

```
Go-Backend/
├── cmd/
│   └── main.go                  # Main application entry point
├── internal/
│   ├── config/
│   │   └── db.go                # Database configuration and initialization
│   ├── routes/
│   │   ├── assessment_routes.go # Routes related to assessments
│   │   ├── classroom_routes.go  # Routes related to classrooms
│   │   ├── doubt_routes.go      # Routes related to doubt solving
│   │   ├── quiz_routes.go       # Routes related to quizzes
│   │   ├── recommendation_routes.go # Routes related to recommendations
│   │   ├── student_routes.go    # Routes related to students
│   │   └── teacher_routes.go    # Routes related to teachers
│   └── services/
│       └── ml_service.go        # Service interacting with ML-Service repository
├── go.mod                       # Go module file
└── README.md                    # Project documentation
```

## Interaction with ML-Service

The Go-Backend interacts with the ML-Service repository to leverage machine learning models for educational services. Here is how the interaction works:

1. **Assessment Routes**: The backend sends student responses to the ML-Service to get assessments and predictions.
2. **Quiz Routes**: The backend requests personalized quizzes from the ML-Service based on student profiles and learning styles.
3. **Tutor Routes**: The backend utilizes the ML-Service to provide context-aware tutoring sessions.
4. **Doubt Routes**: The backend sends student doubts (text and images) to the ML-Service for resolution using multimodal AI agents.
5. **Recommendation Routes**: The backend requests study routines and learning resources from the ML-Service based on aggregated student data.

Example of a service file interacting with ML-Service (internal/services/ml_service.go):

```go
package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AssessmentRequest struct {
	Responses []int `json:"responses"`
}

type AssessmentResponse struct {
	StudyType   string `json:"study_type"`
	Description string `json:"description"`
}

func GetAssessment(responses []int) (*AssessmentResponse, error) {
	reqBody, _ := json.Marshal(AssessmentRequest{Responses: responses})
	resp, err := http.Post("https://ml-service/api/assessments/static", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var assessmentResp AssessmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&assessmentResp); err != nil {
		return nil, err
	}

	return &assessmentResp, nil
}
```

## Contributing

We welcome contributions from the community. To contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes to the branch.
4. Create a pull request with a detailed description of the changes.

---
