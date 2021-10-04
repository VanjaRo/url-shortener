package redis

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"url-shortener/base62"
	"url-shortener/storage"

	redisClient "github.com/gomodule/redigo/redis"
)

type redis struct{ pool *redisClient.Pool }

func New(host, port string) (storage.Service, error) {
	pool := &redisClient.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redisClient.Conn, error) {
			return redisClient.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}
	return &redis{pool}, nil
}

func (r *redis) isUsed(id uint64) bool {
	conn := r.pool.Get()
	defer conn.Close()

	exists, err := redisClient.Bool(conn.Do("EXISTS", strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (r *redis) Save(url string, expires time.Time) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	var id uint64

	for used := true; used; used = r.isUsed(id) {
		id = rand.Uint64()
	}

	shortLink := storage.Item{
		id,
		url,
		expires.Format("1999-01-01 11:11:11.000000 +0300 EEST"),
		0,
	}
	_, err := conn.Do("HMSET", redisClient.Args{"Shotener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)
	if err != nil {
		return "", err
	}

	_, err = conn.Do("EXPIREAT", "Shotener:"+strconv.FormatUint(id, 10), expires.Unix())
	if err != nil {
		return "", err
	}

	return base62.Encode(id), nil

}

func (r *redis) Load(code string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil {
		return "", err
	}

	urlString, err := redisClient.String(conn.Do("HGET", "Shortener: "+strconv.FormatUint(decodedId, 10), "url"))
	if err != nil {
		return "", nil
	} else if len(urlString) == 0 {
		return "", storage.ErrNoLink
	}

	conn.Do("HINCBY", "Shortener: "+strconv.FormatUint(decodedId, 10), "visits", 1)

	return urlString, nil
}

func (r *redis) LoadInfo(code string) (*storage.Item, error) {
	conn := r.pool.Get()
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	values, err := redisClient.Values(conn.Do("HGETALL", "Shortener: "+strconv.FormatUint(decodedId, 10)))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, storage.ErrNoLink
	}

	var shortLink storage.Item
	if err := redisClient.ScanStruct(values, &shortLink); err != nil {
		return nil, err
	}
	return &shortLink, nil
}

func (r *redis) Close() error {
	return r.pool.Close()
}
