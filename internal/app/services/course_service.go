package services

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
)

type CourseService struct {
	repo *repositories.CourseRepository
}

func NewCourseService() *CourseService {
	return &CourseService{
		repo: repositories.NewCourseRepository(),
	}
}

func (s *CourseService) CreateCourse(course *models.Course) error {
	return s.repo.CreateCourse(course)
}
