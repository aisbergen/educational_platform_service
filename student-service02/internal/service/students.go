package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"student/internal/domain"
	"student/internal/repository"
	"student/pkg/auth"
	"student/pkg/cache"
	"student/pkg/hash"
	"time"
)

type StudentService struct {
	repo            repository.Students
	hash            hash.PasswordHasher
	cache           cache.Cache
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewStudentService(repo repository.Students, hash hash.PasswordHasher, manager auth.TokenManager, cache cache.Cache, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *StudentService {
	return &StudentService{
		repo:            repo,
		hash:            hash,
		cache:           cache,
		tokenManager:    manager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *StudentService) Create(ctx context.Context, student domain.Student) error {
	var err error
	student.Password, err = s.hash.GenerateFromPassword(student.Password)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, student)

}

func (s *StudentService) GetByEmail(ctx context.Context, email string, password string) (domain.Token, error) {
	student, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("error service: %s", err)
		return domain.Token{}, domain.ErrNotFound
	}

	if err = s.hash.CompareHashAndPassword(student.Password, password); err != nil {
		log.Printf("error service: %s", err)
		return domain.Token{}, domain.ErrNotFound
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentService) GetStudentByID(ctx context.Context, id int) (domain.Student, error) {
	cachedStudent, err := s.cache.Get("student:" + strconv.Itoa(id))
	if err == nil {
		cachedData, ok := cachedStudent.(map[string]interface{})
		if ok {
			student := domain.Student{
				ID:      int(cachedData["id"].(float64)),
				Email:   cachedData["email"].(string),
				Name:    cachedData["name"].(string),
				GPA:     cachedData["gpa"].(float64),
				Courses: convertToStringSlice(cachedData["courses"].([]interface{})),
			}

			return student, nil
		}

	}

	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Student{}, err
	}
	log.Println(student)
	err = s.cache.Set("student:"+strconv.Itoa(id), student, s.accessTokenTTL)
	if err != nil {
		log.Printf("error caching student with ID %d: %s", id, err)
	}

	return student, nil
}

func (s *StudentService) Update(ctx context.Context, student domain.Student) error {
	stud, err := s.repo.GetByID(ctx, student.ID)
	if err != nil {
		return err
	}

	if stud.ID != student.ID {
		return errors.New("ID mismatch: retrieved student ID does not match the expected ID")
	}

	if len(student.Courses) == 0 {
		student.Courses = stud.Courses
	}

	if student.Email == "" {
		student.Email = stud.Email
	}

	if student.Name == "" {
		student.Name = stud.Name
	}

	if student.GPA == 0 {
		student.GPA = stud.GPA
	}

	if len(student.Password) == 0 {
		student.Password = stud.Password
	} else {
		student.Password, err = s.hash.GenerateFromPassword(student.Password)
		if err != nil {
			return err
		}
	}

	log.Println(student)

	return s.repo.Update(ctx, student)

}

func (s *StudentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *StudentService) GetStudentsByCoursesID(ctx context.Context, id string) ([]domain.Student, error) {
	return s.repo.GetStudentsByCoursesID(ctx, id)
}
func (s *StudentService) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.Token, error) {

	student, err := s.repo.GetByRefresh(ctx, refreshToken)
	if err != nil {
		return domain.Token{}, domain.ErrNotFound
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentService) createSession(ctx context.Context, studentID int) (domain.Token, error) {
	var (
		res domain.Token
		err error
	)
	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(studentID), s.accessTokenTTL)
	if err != nil {
		return res, err
	}
	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}
	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}
	err = s.repo.SetSession(ctx, session, studentID)

	return res, err
}

func (s *StudentService) DeleteStudentFromCache(id int) error {
	err := s.cache.Delete("student:" + strconv.Itoa(id))
	if err != nil {
		log.Printf("error deleting student with ID %d from cache: %s", id, err)
		return err
	}
	return nil
}

func convertToStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, val := range slice {
		result[i] = val.(string)
	}
	return result
}
