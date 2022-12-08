package file

import (
	"os"
	"path/filepath"
	"promhsd/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDB_readFile(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileOk.Name(), []byte("{}"), 0644)
	fileWrong, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileWrong.Name(), []byte("----"), 0644)
	defer os.Remove(fileOk.Name())
	defer os.Remove(fileWrong.Name())
	type fields struct {
		filepath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]db.Target
		wantErr bool
	}{
		{
			name: "FileIsOk",
			fields: fields{
				filepath: fileOk.Name(),
			},
			want:    map[string]db.Target{},
			wantErr: false,
		},
		{
			name: "FileNameIsWrong",
			fields: fields{
				filepath: "wrongName.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "FileDataIsInvalid",
			fields: fields{
				filepath: fileWrong.Name(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			got, err := f.readFile()
			if !tt.wantErr {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileDB_writeToFile(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	defer os.Remove(fileOk.Name())
	type fields struct {
		filepath string
	}
	type args struct {
		targets map[string]db.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "argsIsNil",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args:    args{targets: nil},
			wantErr: false,
		},
		{
			name: "fileIsWrong",
			fields: fields{
				filepath: filepath.Join("dirNotExist", "wrongfile.json"),
			},
			args:    args{targets: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			err := f.writeToFile(tt.args.targets)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestFileDB_Create(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileOk.Name(), []byte("{}"), 0644)
	defer os.Remove(fileOk.Name())
	type fields struct {
		filepath string
	}
	type args struct {
		target *db.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fileOk",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test"},
			},
			wantErr: false,
		},
		{
			name: "filenameWrong",
			fields: fields{
				filepath: filepath.Join("dirNotExist", "wrong.json"),
			},
			args: args{
				target: &db.Target{Name: "test"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			err := f.Create(tt.args.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileDB_Update(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileOk.Name(), []byte(`{"test": {"name": "test"}}`), 0644)
	defer os.Remove(fileOk.Name())
	type fields struct {
		filepath string
	}
	type args struct {
		target *db.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fileOk",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test", ID: "test"},
			},
			wantErr: false,
		},
		{
			name: "filenameWrong",
			fields: fields{
				filepath: filepath.Join("dirNotExist", "wrong.json"),
			},
			args: args{
				target: &db.Target{Name: "test"},
			},
			wantErr: true,
		},
		{
			name: "idNotFound",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test", ID: "test1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			err := f.Update(tt.args.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileDB_Delete(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileOk.Name(), []byte(`{"test": {"name": "test"}}`), 0644)
	defer os.Remove(fileOk.Name())
	type fields struct {
		filepath string
	}
	type args struct {
		target *db.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fileOk",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test", ID: "test"},
			},
			wantErr: false,
		},
		{
			name: "keyNotExists",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test", ID: "test1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			err := f.Delete(tt.args.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileDB_Get(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileOk.Name(), []byte(`{"test": {"name": "test"}}`), 0644)
	defer os.Remove(fileOk.Name())
	type fields struct {
		filepath string
	}
	type args struct {
		target *db.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fileOk",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test", ID: "test"},
			},
			wantErr: false,
		},
		{
			name: "keyNotExist",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				target: &db.Target{Name: "test", ID: "test1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			err := f.Get(tt.args.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileDB_GetAll(t *testing.T) {
	fileOk, err := os.CreateTemp("", "promhsd-*")
	assert.NoError(t, err)
	os.WriteFile(fileOk.Name(), []byte(`{"test": {"name": "test"}}`), 0644)
	defer os.Remove(fileOk.Name())
	list := []db.Target{}
	type fields struct {
		filepath string
	}
	type args struct {
		list *[]db.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fileOk",
			fields: fields{
				filepath: fileOk.Name(),
			},
			args: args{
				list: &list,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				filepath: tt.fields.filepath,
			}
			err := f.GetAll(tt.args.list)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *FileDB
	}{
		{
			name: "testFile",
			args: args{
				path: "test.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := StorageService{}
			got, err := service.New(tt.args.path)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}
