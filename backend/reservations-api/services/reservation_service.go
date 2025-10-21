package services

import (
	"context"
	"time"

	repository "reservations/DAO"
	errorspkg "reservations/errors"
	"reservations/messaging"
	"reservations/models"
	"reservations/utils"
)

// ReservationService defines business operations
type ReservationService interface {
	Create(ctx context.Context, r *models.Reservation) (string, error)
	Get(ctx context.Context, id string) (*models.Reservation, error)
	Update(ctx context.Context, id string, r *models.Reservation) error
	Delete(ctx context.Context, id string) error
}

type reservationService struct {
	repo  repository.ReservationRepository
	pub   messaging.Publisher
	users UserService
	tasks int
}

func NewReservationService(repo repository.ReservationRepository, pub messaging.Publisher, users UserService, tasks int) ReservationService {
	return &reservationService{repo: repo, pub: pub, users: users, tasks: tasks}
}

func (s *reservationService) Create(ctx context.Context, r *models.Reservation) (string, error) {
	// validate owner
	ok, err := s.users.Exists(ctx, r.OwnerID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errorspkg.ErrOwnerNotValid
	}
	// concurrent calc
	calcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	score, err := utils.ComputeScore(calcCtx, s.tasks)
	if err != nil {
		return "", errorspkg.ErrConcurrentCalc
	}
	r.Score = score
	id, err := s.repo.Insert(ctx, r)
	if err != nil {
		return "", err
	}
	// notify async (best-effort)
	_ = s.pub.Publish("create", id)
	return id, nil
}

func (s *reservationService) Get(ctx context.Context, id string) (*models.Reservation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *reservationService) Update(ctx context.Context, id string, r *models.Reservation) error {
	// validate owner if present
	if r.OwnerID != "" {
		ok, err := s.users.Exists(ctx, r.OwnerID)
		if err != nil {
			return err
		}
		if !ok {
			return errorspkg.ErrOwnerNotValid
		}
	}
	// recalc
	calcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	score, err := utils.ComputeScore(calcCtx, s.tasks)
	if err != nil {
		return errorspkg.ErrConcurrentCalc
	}
	r.Score = score
	if err := s.repo.UpdateByID(ctx, id, r); err != nil {
		return err
	}
	_ = s.pub.Publish("update", id)
	return nil
}

func (s *reservationService) Delete(ctx context.Context, id string) error {
	// read to ensure exists and owner validation is done at repo/service level if required
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	ok, err := s.users.Exists(ctx, r.OwnerID)
	if err != nil {
		return err
	}
	if !ok {
		return errorspkg.ErrOwnerNotValid
	}
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return err
	}
	_ = s.pub.Publish("delete", id)
	return nil
}
