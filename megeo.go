package megeo

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type GEOItems []*GEOItem
type GEOItem struct {
	Longitude float64
	Latitude  float64
	Member    uint64
}

type GEORadiusItems []*GEORadiusItem
type GEORadiusItem struct {
	Member   uint64
	Distance float64
}

type MeGEO struct {
	conn     redis.Conn
	resource string
}

func NewMeGEO(conn redis.Conn, resource string) *MeGEO {
	return &MeGEO{
		conn:     conn,
		resource: resource,
	}
}

func (geo *MeGEO) BatchAdd(items GEOItems) (num int64, err error) {
	args := make([]interface{}, 0)
	args = append(args, geo.resource)
	for _, item := range items {
		args = append(args, item.Longitude, item.Latitude, item.Member)
	}
	count, err := geo.conn.Do("GEOADD", args...)
	if err != nil {
		err = fmt.Errorf("megeo BatchAdd geo.conn.Do err = %v", err)
		return
	}
	num = count.(int64)
	return
}

func (geo *MeGEO) Add(item GEOItem) (num int64, err error) {
	args := make([]interface{}, 0)
	args = append(args, geo.resource, item.Longitude, item.Latitude, item.Member)
	count, err := geo.conn.Do("GEOADD", args...)
	if err != nil {
		err = fmt.Errorf("megeo Add geo.conn.Do err = %v", err)
		return
	}
	num = count.(int64)
	return
}

func (geo *MeGEO) Radius(longitude float64, latitude float64, radius uint) (items GEORadiusItems, err error) {
	//GEORADIUS stores:locations 113.3172 23.08331 100 km WITHDIST ASC
	args := make([]interface{}, 0)
	args = append(args, geo.resource, longitude, latitude, radius, "km", "WITHDIST", "ASC")
	values, err := redis.Values(geo.conn.Do("GEORADIUS", args...))
	if err != nil {
		err = fmt.Errorf("megeo Radius geo.conn.Do err = %v", err)
		return
	}
	//var item GEORadiusItem
	items = make(GEORadiusItems, 0)
	var item *GEORadiusItem
	for _, value := range values {
		vals := value.([]interface{})
		if len(vals) != 2 {
			err = errors.New("megeo Radius redis GEORADIUS return unexpected array ")
			return
		}
		item = new(GEORadiusItem)
		item.Member, err = redis.Uint64(vals[0], nil)
		if err != nil {
			err = fmt.Errorf("megeo Radius item.Member redis.Uint64 err = %v", err)
			return
		}
		item.Distance, err = redis.Float64(vals[1], nil)
		if err != nil {
			err = fmt.Errorf("megeo Radius item.Distance redis.Float64 err = %v", err)
			return
		}
		items = append(items, item)
	}
	return
}
