package services

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
)

type CourseService struct {
	courseRepository     *repositories.CourseRepository
	userCourseRepository *repositories.UserCourseRepository
}

func NewCourseService() *CourseService {
	return &CourseService{
		courseRepository:     repositories.NewCourseRepository(),
		userCourseRepository: repositories.NewUserCourseRepository(),
	}
}

func (s *CourseService) CreateCourse(course *models.CreateCourseRequest) error {
	return s.courseRepository.Create(course)
}

func (s *CourseService) UpdateCourse(course *models.UpdateCourseRequest) error {
	return s.courseRepository.Update(course)
}

func (s *CourseService) GetCoursesByClassID(classID string) ([]models.Course, error) {
	return s.courseRepository.FindByClassId(classID)
}

func (s *CourseService) EnrollCourse(userCourse *models.UserCourse) error {
	return s.userCourseRepository.Create(userCourse)
}
