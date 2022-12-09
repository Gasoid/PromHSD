package db

import (
	"log"
	"time"
)

var (
	nilID    ID
	storages = map[string]func(string) (Storage, error){}
)

type StorageService interface {
	ServiceID() string
	New(string) (Storage, error)
}

func RegisterStorage(storage StorageService) {
	storages[storage.ServiceID()] = storage.New
}

type ID string

func (id ID) String() string {
	return string(id)
}

type Storage interface {
	Create(*Target) error
	Update(*Target) error
	Delete(*Target) error
	Get(*Target) error
	GetAll(*[]Target) error
}

type Target struct {
	ID      ID        `json:"id"`
	Name    string    `json:"name"`
	Time    time.Time `json:"time"`
	Entries []Entry   `json:"entries"`
}

type Entry struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type Service struct {
	storage Storage
}

func (s *Service) Create(target *Target) error {
	if err := target.validate(); err != nil {
		return err
	}
	target.Time = time.Now()
	return s.storage.Create(target)
}

func (s *Service) Update(target *Target) error {
	if target.ID == nilID {
		return ErrValidation
	}
	if err := target.validate(); err != nil {
		return err
	}
	target.Time = time.Now()
	return s.storage.Update(target)
}

func (s *Service) Delete(target *Target) error {
	if target.ID == nilID {
		return ErrValidation
	}
	return s.storage.Delete(target)
}

func (s *Service) Get(target *Target) error {
	if target.ID == nilID {
		return ErrValidation
	}
	err := s.storage.Get(target)
	if err != nil {
		log.Println("(Get) Storage returned error: ", err.Error())
	}
	return err
}

func (s *Service) List(targets *[]Target) error {
	err := s.storage.GetAll(targets)
	if err != nil {
		log.Println("(GetAll) Storage returned error: ", err.Error())
	}
	return err
}

func (t *Target) validate() error {
	if t.Name == "" {
		return ErrValidation
	}
	if t.Entries == nil || len(t.Entries) == 0 {
		return ErrValidation
	}
	for _, e := range t.Entries {
		if err := e.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Entry) validate() error {
	if e.Targets == nil || len(e.Targets) == 0 {
		return ErrValidation
	}
	if e.Labels == nil || len(e.Labels) == 0 {
		return ErrValidation
	}
	return nil
}

func New(storageID, opt string) (*Service, error) {
	if newFunc, ok := storages[storageID]; ok {
		storage, err := newFunc(opt)
		if err != nil {
			return nil, &StorageError{Text: "storage returned error", Err: err}
		}
		return &Service{storage}, nil
	}
	return nil, &StorageError{Text: "storage is not implemented"}
}

func NewTarget() *Target {
	return &Target{Entries: []Entry{}}
}

func NewEntry() *Entry {
	return &Entry{Labels: map[string]string{}, Targets: []string{}}
}
