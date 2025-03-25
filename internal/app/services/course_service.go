package services

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
	"fmt"
	"log"
	"strings"
)

type CreateCourseRequest struct {
	ID              string `json:"id,omitempty"`
	ClassID         string `json:"class_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DifficultyLevel string `json:"difficulty_level"`
	Number          int    `json:"number"`
}

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

func (s *CourseService) CreateCourse(request *CreateCourseRequest) (*models.Course, error) {
	course := new(models.Course)
	course.ClassID = request.ClassID
	course.Title = request.Title
	course.Description = request.Description
	course.DifficultyLevel = request.DifficultyLevel
	course.Number = request.Number

	err := s.courseRepository.Create(course)
	if err != nil {
		return nil, fmt.Errorf("Error creating course: %w", err)
	}

	return course, nil
}

func (s *CourseService) UpdateCourse(course *models.UpdateCourseRequest) error {
	return s.courseRepository.Update(course)
}

func (s *CourseService) GetCoursesByClassID(classID string) ([]models.Course, error) {
	courses, err := s.courseRepository.FindByClassId(classID)
	if err != nil {
		return nil, err
	}
	if len(courses) == 0 {
		return nil, fmt.Errorf("no courses found for classID: %s", classID)
	}
	return courses, nil
}

func (s *CourseService) EnrollCourse(userCourse *models.UserCourse) error {
	return s.userCourseRepository.Create(userCourse)
}

func (s *CourseService) DeleteCourseByID(courseID string) error {
	log.Println("Service: Deleting course with id: ", courseID)

	err := s.courseRepository.DeleteCourseByID(courseID)
	if err != nil {
		if strings.Contains(err.Error(), "no course found with id") {
			return fmt.Errorf("Course not found: %w", err)
		}
		return fmt.Errorf("Error deleting course: %w", err)
	}

	return nil
}
