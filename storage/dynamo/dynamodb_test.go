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

type testDescribeTable struct {
	err error
}

func (item *testDescribeTable) DescribeTable(*dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
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
				IDescribeTable: &testDescribeTable{err: awserr.New(dynamodb.ErrCodeResourceNotFoundException, "notFound", nil)},
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
