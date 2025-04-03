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

type CourseService interface {
	CreateCourse(request *CreateCourseRequest) (*models.Course, error)
	UpdateCourse(course *models.UpdateCourseRequest) error
	GetCoursesByClassID(classID string) ([]models.Course, error)
	EnrollCourse(userCourse *models.UserCourse) error
	DeleteCourseByID(courseID string) error
	GetCourseByID(courseID string) (*models.Course, error)
}

type CourseServiceImpl struct {
	courseRepository     *repositories.CourseRepository
	classRepository      *repositories.ClassRepository
	userCourseRepository *repositories.UserCourseRepository
}

type courseServiceImpl struct {
	courseRepository     *repositories.CourseRepository
	classRepository      *repositories.ClassRepository
	userCourseRepository *repositories.UserCourseRepository
}

func NewCourseService() CourseService {
	return &courseServiceImpl{
		courseRepository:     repositories.NewCourseRepository(),
		classRepository:      repositories.NewClassRepository(),
		userCourseRepository: repositories.NewUserCourseRepository(),
	}
}

func (s *courseServiceImpl) CreateCourse(request *CreateCourseRequest) (*models.Course, error) {
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

	err = s.createUserCoursesForClass(course.ClassID, course.ID)
	if err != nil {
		return nil, fmt.Errorf("Error creating user courses: %w", err)
	}

	return course, nil
}

func (s *courseServiceImpl) createUserCoursesForClass(classID, courseID string) error {
	userIDs, err := s.classRepository.GetUsersByClassID(classID)
	if err != nil {
		return fmt.Errorf("Error fetching users from user_classes: %w", err)
	}

	for _, userID := range userIDs {
		userCourse := &models.UserCourse{
			UserID:   userID.ID,
			CourseID: courseID,
		}

		err := s.userCourseRepository.Create(userCourse)
		if err != nil {
			return fmt.Errorf("Error creating user_course for user %s: %w", userID, err)
		}
	}

	return nil
}

func (s *courseServiceImpl) UpdateCourse(course *models.UpdateCourseRequest) error {
	return s.courseRepository.Update(course)
}

func (s *courseServiceImpl) GetCoursesByClassID(classID string) ([]models.Course, error) {
	courses, err := s.courseRepository.FindByClassId(classID)
	if err != nil {
		return nil, err
	}
	if len(courses) == 0 {
		return nil, fmt.Errorf("no courses found for classID: %s", classID)
	}
	return courses, nil
}

func (s *courseServiceImpl) EnrollCourse(userCourse *models.UserCourse) error {
	return s.userCourseRepository.Create(userCourse)
}

func (s *courseServiceImpl) DeleteCourseByID(courseID string) error {
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

func (s *courseServiceImpl) GetCourseByID(courseID string) (*models.Course, error) {
	course, err := s.courseRepository.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, fmt.Errorf("course not found")
	}
	return course, nil
}