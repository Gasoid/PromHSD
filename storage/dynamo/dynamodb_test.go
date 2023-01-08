package dynamo

import (
	"promhsd/db"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

type testGetItem struct {
	err error
}

func (item *testGetItem) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, item.err
}

type testPutItem struct {
	err error
}

func (item *testPutItem) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, item.err
}

func (item *testPutItem) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return nil, item.err
}

type testDeleteItem struct {
	err error
}

func (item *testDeleteItem) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, item.err
}

type testTable struct {
	err error
}

func (item *testTable) DescribeTable(*dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return nil, item.err
}

func (item *testTable) CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	return nil, item.err
}

type testScan struct {
	err error
}

func (item *testScan) Scan(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return nil, item.err
}

func TestDynamoDB_Get(t *testing.T) {
	type fields struct {
		IGetItem  IGetItem
		tableName string
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
			name:    "NotFoundError",
			fields:  fields{IGetItem: &testGetItem{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)}},
			args:    args{target: db.NewTarget()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				IGetItem:  tt.fields.IGetItem,
				tableName: tt.fields.tableName,
			}
			if err := d.Get(tt.args.target); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestDynamoDB_Update(t *testing.T) {
	target := db.NewTarget()
	target.ID = db.ID("id")
	target.Name = "example"
	type fields struct {
		IPutItem  IPutItem
		tableName string
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
			name:    "NotFoundError",
			fields:  fields{IPutItem: &testPutItem{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)}},
			args:    args{target: target},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				IPutItem:  tt.fields.IPutItem,
				tableName: tt.fields.tableName,
			}
			if err := d.Update(tt.args.target); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestDynamoDB_Delete(t *testing.T) {
	type fields struct {
		IDeleteItem IDeleteItem
		tableName   string
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
			name: "NotFoundError",
			fields: fields{
				IDeleteItem: &testDeleteItem{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
			},
			args:    args{target: db.NewTarget()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				IDeleteItem: tt.fields.IDeleteItem,
				tableName:   tt.fields.tableName,
			}
			if err := d.Delete(tt.args.target); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestDynamoDB_Create(t *testing.T) {
	target := db.NewTarget()
	target.ID = db.ID("id")
	target.Name = "example"
	type fields struct {
		IPutItem  IPutItem
		IGetItem  IGetItem
		tableName string
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
			name: "NotFoundError",
			fields: fields{
				IPutItem: &testPutItem{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
				IGetItem: &testPutItem{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
			},
			args:    args{target: target},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				IPutItem:  tt.fields.IPutItem,
				IGetItem:  tt.fields.IGetItem,
				tableName: tt.fields.tableName,
			}
			if err := d.Create(tt.args.target); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestDynamoDB_IsHealthy(t *testing.T) {
	type fields struct {
		IDescribeTable IDescribeTable
		tableName      string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "NotFoundError",
			fields: fields{
				IDescribeTable: &testTable{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				IDescribeTable: tt.fields.IDescribeTable,
				tableName:      tt.fields.tableName,
			}
			assert.Equal(t, tt.want, d.IsHealthy())
		})
	}
}

func TestDynamoDB_GetAll(t *testing.T) {
	type fields struct {
		IScan     IScan
		tableName string
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
			name: "NotFoundError",
			fields: fields{
				IScan: &testScan{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
			},
			args:    args{list: &[]db.Target{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				IScan:     tt.fields.IScan,
				tableName: tt.fields.tableName,
			}
			if err := d.GetAll(tt.args.list); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestDynamoDB_createTable(t *testing.T) {
	type fields struct {
		ICreateTable ICreateTable
		tableName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Error",
			fields: fields{
				ICreateTable: &testTable{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DynamoDB{
				ICreateTable: tt.fields.ICreateTable,
				tableName:    tt.fields.tableName,
			}
			if err := d.createTable(); (err != nil) != tt.wantErr {
				t.Errorf("DynamoDB.createTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
