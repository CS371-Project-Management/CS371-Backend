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

func (s *CourseService) CreateCourse(course *models.CreateCourseRequest) error {
	return s.repo.Create(course)
}

func (s *CourseService) UpdateCourse(course *models.UpdateCourseRequest) error {
	return s.repo.Update(course)
}

func (s *CourseService) GetCoursesByClassID(classID string) ([]models.Course, error) {
	return s.repo.FindByClassId(classID)
}
