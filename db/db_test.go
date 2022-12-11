package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStorage struct {
	returnError error
	//returnItem  *Target
}

func (s *testStorage) Create(*Target) error {

	return s.returnError
}

func (s *testStorage) Update(*Target) error {
	return s.returnError
}

func (s *testStorage) Delete(*Target) error {
	return s.returnError
}

func (s *testStorage) Get(*Target) error {
	return s.returnError
}

func (s *testStorage) GetAll(*[]Target) error {
	return s.returnError
}

func (s *testStorage) IsHealthy() bool {
	return true
}

type testStorageService struct {
	storage Storage
}

func (s *testStorageService) ServiceID() string {
	return "testdb"
}

func (s *testStorageService) New(opt string) (Storage, error) {
	return s.storage, nil
}

func TestService_Get(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		target *Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{&testStorage{returnError: ErrValidation}},
			args:    args{&Target{ID: "key1"}},
			wantErr: true,
		},
		{
			name:    "test2",
			fields:  fields{&testStorage{}},
			args:    args{&Target{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			err := s.Get(tt.args.target)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	entryNoLabels := Entry{
		Targets: []string{"asd"},
	}
	entryNoTargets := Entry{
		Labels: map[string]string{"label1": "value"},
	}
	entry := Entry{
		Labels:  map[string]string{"label1": "value"},
		Targets: []string{"asd"},
	}
	type fields struct {
		storage Storage
	}
	type args struct {
		target *Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "NoEntries",
			fields:  fields{&testStorage{}},
			args:    args{&Target{Name: "test"}},
			wantErr: true,
		},
		{
			name:    "NoFields",
			fields:  fields{&testStorage{}},
			args:    args{&Target{}},
			wantErr: true,
		},
		{
			name:    "NoLabels",
			fields:  fields{&testStorage{}},
			args:    args{&Target{Name: "test", Entries: []Entry{entryNoLabels}}},
			wantErr: true,
		},
		{
			name:    "NoTargets",
			fields:  fields{&testStorage{}},
			args:    args{&Target{Name: "test", Entries: []Entry{entryNoTargets}}},
			wantErr: true,
		},
		{
			name:    "ValidationSucceded",
			fields:  fields{&testStorage{}},
			args:    args{&Target{Name: "test", Entries: []Entry{entry}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			err := s.Create(tt.args.target)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	entry := Entry{
		Labels:  map[string]string{"label1": "value"},
		Targets: []string{"asd"},
	}
	type fields struct {
		storage Storage
	}
	type args struct {
		target *Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "NoEntries",
			fields:  fields{&testStorage{}},
			args:    args{&Target{ID: "test", Name: "test"}},
			wantErr: true,
		},
		{
			name:    "NoID",
			fields:  fields{&testStorage{}},
			args:    args{&Target{Name: "test"}},
			wantErr: true,
		},
		{
			name:    "ValidationSucceded",
			fields:  fields{&testStorage{}},
			args:    args{&Target{ID: "test", Name: "test", Entries: []Entry{entry}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			err := s.Update(tt.args.target)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		target *Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "NoID",
			fields:  fields{&testStorage{}},
			args:    args{&Target{}},
			wantErr: true,
		},
		{
			name:    "ValidationSucceded",
			fields:  fields{&testStorage{}},
			args:    args{&Target{ID: "test"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			err := s.Delete(tt.args.target)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	targets := []Target{}
	type fields struct {
		storage Storage
	}
	type args struct {
		targets *[]Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "ListOK",
			fields:  fields{&testStorage{}},
			args:    args{&targets},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			err := s.List(tt.args.targets)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		storage Storage
	}
	tests := []struct {
		name    string
		args    args
		want    *Service
		wantErr bool
	}{
		// {
		// 	name:    "NilStorage",
		// 	args:    args{nil},
		// 	want:    nil,
		// 	wantErr: true,
		// },
		{
			name:    "DumbStorage",
			args:    args{&testStorage{}},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterStorage(&testStorageService{storage: tt.args.storage})
			got, err := New("testdb", "")
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			}
			assert.IsType(t, &Service{}, got)
		})
	}
}

func TestNewTarget(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "DumbTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTarget()
			assert.IsType(t, &Target{}, got)
		})
	}
}

func TestNewEntry(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "DumbTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEntry()
			assert.IsType(t, &Entry{}, got)
		})
	}
}

func TestID_String(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "string",
			id:   ID("1"),
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.id.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
