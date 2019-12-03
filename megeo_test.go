package megeo

import (
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
)

var conn redis.Conn

func setup() {
	var err error
	conn, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
}

func teardown() {
	conn.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func BenchmarkMeGEO_Radius(b *testing.B) {
	setup()
	geo := NewMeGEO(conn, "stores:loc")
	for i := 0; i < b.N; i++ {
		_, err := geo.Radius(113.3172, 23.08331, 10)
		if err != nil {
			b.Error(err)
		}
	}
	teardown()
}

func TestMeGEO_Radius(t *testing.T) {
	type fields struct {
		conn     redis.Conn
		resource string
	}
	type args struct {
		longitude float64
		latitude  float64
		radius    uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "radius",
			fields: fields{
				conn:     conn,
				resource: "stores:loc",
			},
			args: args{
				longitude: 113.3172,
				latitude:  23.08331,
				radius:    10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geo := &MeGEO{
				conn:     tt.fields.conn,
				resource: tt.fields.resource,
			}
			items, err := geo.Radius(tt.args.longitude, tt.args.latitude, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("MeGEO.Radius() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, item := range items {
				t.Log(*item)
			}
		})
	}
}

func TestMeGEO_BatchAdd(t *testing.T) {
	type fields struct {
		conn     redis.Conn
		resource string
	}
	type args struct {
		items GEOItems
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNum int64
		wantErr bool
	}{
		{
			name: "batch add",
			fields: fields{
				conn:     conn,
				resource: "stores:loc",
			},
			args: args{
				items: GEOItems{
					{
						Longitude: 113.33752,
						Latitude:  23.09787,
						Member:    441,
					},
					{
						Longitude: 113.321316,
						Latitude:  23.095648,
						Member:    453,
					},
				},
			},
			wantNum: 2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geo := &MeGEO{
				conn:     tt.fields.conn,
				resource: tt.fields.resource,
			}
			gotNum, err := geo.BatchAdd(tt.args.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("MeGEO.BatchAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNum != tt.wantNum {
				t.Errorf("MeGEO.BatchAdd() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestMeGEO_Add(t *testing.T) {
	type fields struct {
		conn     redis.Conn
		resource string
	}
	type args struct {
		item GEOItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNum int64
		wantErr bool
	}{
		{
			name: "add",
			fields: fields{
				conn:     conn,
				resource: "stores:loc",
			},
			args: args{
				item: GEOItem{
					Longitude: 113.274065,
					Latitude:  23.095387,
					Member:    404,
				},
			},
			wantNum: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geo := &MeGEO{
				conn:     tt.fields.conn,
				resource: tt.fields.resource,
			}
			gotNum, err := geo.Add(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("MeGEO.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNum != tt.wantNum {
				t.Errorf("MeGEO.Add() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}
