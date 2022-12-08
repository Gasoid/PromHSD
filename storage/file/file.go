package file

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"promhsd/db"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

const (
	StorageID = "filedb"
)

type FileDB struct {
	filepath string
	mu       sync.Mutex
	filelock *flock.Flock
}

func (f *FileDB) readFile() (map[string]db.Target, error) {
	jsonFile, err := os.Open(f.filepath)
	if err != nil {
		log.Println(err)
		return nil, &db.StorageError{Text: "Couldn't open a file", Err: err}
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var targets map[string]db.Target
	err = json.Unmarshal(byteValue, &targets)
	if err != nil {
		log.Println(err)
		return nil, &db.StorageError{Text: "Couldn't decode json", Err: err}
	}
	return targets, nil
}

func (f *FileDB) Lock() error {
	f.mu.Lock()
	f.filelock = flock.New(f.filepath + ".lock")
	for i := 0; i < 10; i++ {
		locked, err := f.filelock.TryLock()
		if err != nil {
			return err
		}

		if locked {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("couldn't lock file")
}

func (f *FileDB) Unlock() error {
	f.mu.Unlock()
	return f.filelock.Unlock()
}

func (f *FileDB) writeToFile(targets map[string]db.Target) error {
	file, _ := json.Marshal(targets)
	// if err != nil {
	// 	log.Println(err)
	// 	return &db.StorageError{Text: "Couldn't encode to json", Err: err}
	// }
	err := os.WriteFile(f.filepath, file, 0644)
	if err != nil {
		log.Println(err)
		return &db.StorageError{Text: "Couldn't write file", Err: err}
	}
	return nil
}

func (f *FileDB) Create(target *db.Target) error {
	err := f.Lock()
	if err != nil {
		return &db.StorageError{Text: "Couldn't lock file", Err: err}
	}
	defer f.Unlock()
	targets, err := f.readFile()
	if err != nil {
		return err
	}
	target.ID = db.ID(target.Name)
	target.Time = time.Now()
	targets[target.ID.String()] = *target
	err = f.writeToFile(targets)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileDB) Update(target *db.Target) error {
	err := f.Lock()
	if err != nil {
		return &db.StorageError{Text: "Couldn't lock file", Err: err}
	}
	defer f.Unlock()
	targets, err := f.readFile()
	if err != nil {
		return err
	}
	if _, ok := targets[target.ID.String()]; !ok {
		return db.ErrNotFound
	}
	targets[target.ID.String()] = *target
	err = f.writeToFile(targets)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileDB) Delete(target *db.Target) error {
	err := f.Lock()
	if err != nil {
		return &db.StorageError{Text: "Couldn't lock file", Err: err}
	}
	defer f.Unlock()
	targets, err := f.readFile()
	if err != nil {
		return err
	}
	if _, ok := targets[target.ID.String()]; !ok {
		return db.ErrNotFound
	}
	delete(targets, target.ID.String())
	err = f.writeToFile(targets)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileDB) Get(target *db.Target) error {
	targets, err := f.readFile()
	if err != nil {
		return err
	}
	if _, ok := targets[target.ID.String()]; !ok {
		return db.ErrNotFound
	}
	target.Name = targets[target.ID.String()].Name
	target.Entries = targets[target.ID.String()].Entries
	target.Time = targets[target.ID.String()].Time
	return nil
}

func (f *FileDB) GetAll(list *[]db.Target) error {
	targets, err := f.readFile()
	if err != nil {
		return err
	}
	targetList := make([]db.Target, 0, len(targets)-1)
	for _, v := range targets {
		targetList = append(targetList, v)
	}
	*list = targetList
	return nil
}

type StorageService struct{}

func (s *StorageService) ServiceID() string {
	return StorageID
}

func (s *StorageService) New(path string) (db.Storage, error) {
	return &FileDB{filepath: path}, nil
}

func init() {
	db.RegisterStorage(&StorageService{})
}

// binary Search
// func findTarget(name string, targets []db.Target) (int, error) {
// 	left := 0
// 	right := len(targets)
// 	for left <= right {
// 		midPos := (left + right) / 2
// 		ret := strings.Compare(name, targets[midPos].Name)
// 		if ret == 0 {
// 			return midPos, nil
// 		}
// 		if ret < 0 {
// 			right = midPos
// 		} else {
// 			left = midPos + 1
// 		}

// 	}
// 	return -1, errors.New("item doesn't exist")
// }

// func findTargetByNameTime(name string, timestamp time.Time, targets []db.Target) (int, error) {
// 	i, err := findTarget(name, targets)
// 	if err != nil {
// 		return -1, errors.New("item doesn't exist")
// 	}
// 	for j := range targets[i:] {
// 		if targets[j+i].Time == timestamp {
// 			return j + i, nil
// 		}
// 	}
// 	return -1, errors.New("item doesn't exist")
// }

// func sortTargets(targets []db.Target) {
// 	sort.SliceStable(targets, func(i, j int) bool {
// 		if strings.Compare(targets[i].Name, targets[j].Name) == -1 {
// 			return true
// 		}
// 		if strings.Compare(targets[i].Name, targets[j].Name) == 0 {
// 			return true
// 		}
// 		return false
// 	})
// 	for i := range targets {
// 		targets[i].ID = db.ID(strconv.Itoa(i))
// 	}
// }
