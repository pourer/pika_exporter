package exporter

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	defaultScanCount = 100
)

const (
	keyTypeNone   = "none"
	keyTypeString = "string"
	//keyTypeHyperLogLog = "hyperloglog"
	keyTypeList = "list"
	keyTypeSet  = "set"
	keyTypeZSet = "zset"
	keyTypeHash = "hash"
)

var (
	errNotFound = errors.New("key not found")
)

type keyInfo struct {
	size    float64
	keyType string
}

type client struct {
	addr, alias string
	conn        redis.Conn
}

func newClient(addr, password, alias string) (*client, error) {
	conn, err := redis.Dial("tcp", addr,
		redis.DialConnectTimeout(5*time.Second),
		redis.DialWriteTimeout(5*time.Second),
		redis.DialReadTimeout(5*time.Second),
		redis.DialPassword(password))
	if err != nil {
		return nil, err
	}

	return &client{
		addr:  addr,
		alias: alias,
		conn:  conn,
	}, nil
}

func (c *client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *client) Addr() string {
	return c.addr
}

func (c *client) Alias() string {
	return c.alias
}

func (c *client) Select(db string) error {
	_, err := c.conn.Do("SELECT", db)
	return err
}

func (c *client) Info() (string, error) {
	return redis.String(c.conn.Do("INFO", "ALL"))
}

func (c *client) InfoKeySpaceZero() (string, error) {
	return redis.String(c.conn.Do("INFO", "KEYSPACE", 0))
}

func (c *client) InfoKeySpaceOne() (string, error) {
	return redis.String(c.conn.Do("INFO", "KEYSPACE", 1))
}

func (c *client) Del(keys ...string) (int, error) {
	ikeys := make([]interface{}, 0, len(keys))
	for _, k := range keys {
		ikeys = append(ikeys, k)
	}
	return redis.Int(c.conn.Do("DEL", ikeys...))
}

func (c *client) Set(key, value string) (string, error) {
	return redis.String(c.conn.Do("SET", key, value))
}

func (c *client) Hset(key, field, value string) (int, error) {
	return redis.Int(c.conn.Do("HSET", key, field, value))
}

func (c *client) Hget(key, field string) (string, error) {
	return redis.String(c.conn.Do("HGET", key, field))
}

func (c *client) Lpush(key, element string) (int, error) {
	return redis.Int(c.conn.Do("LPUSH", key, element))
}

func (c *client) Lrange(key string, start, stop int) ([]interface{}, error) {
	return redis.Values(c.conn.Do("LRANGE", key, start, stop))
}

func (c *client) Zadd(key string, score float64, member string) (int, error) {
	return redis.Int(c.conn.Do("ZADD", key, score, member))
}

func (c *client) Zcard(key string) (int, error) {
	return redis.Int(c.conn.Do("ZCARD", key))
}

func (c *client) Sadd(key, member string) (int, error) {
	return redis.Int(c.conn.Do("SADD", key, member))
}

func (c *client) Scard(key string) (int, error) {
	return redis.Int(c.conn.Do("SCARD", key))
}

// Pika的SCAN命令，会顺序迭代当前db的快照，由于Pika允许重名五次，所以SCAN有优先输出顺序，依次为：string -> hash -> list -> zset -> set
func (c *client) Scan(keyPattern string, count int) ([]string, error) {
	if count == 0 {
		count = defaultScanCount
	}

	var (
		cursor int
		keys   []string
	)
	for {
		values, err := redis.Values(c.conn.Do("SCAN", cursor, "MATCH", keyPattern, "COUNT", count))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", keyPattern)
		}
		if len(values) != 2 {
			return keys, fmt.Errorf("invalid response from SCAN for pattern: %s", keyPattern)
		}

		ks, _ := redis.Strings(values[1], nil)
		keys = append(keys, ks...)

		if cursor, _ = redis.Int(values[0], nil); cursor == 0 {
			break
		}
	}

	return keys, nil
}

// Pikad的TYPE命令，由于Pika允许重名五次，所以TYPE有优先输出顺序，依次为：string -> hash -> list -> zset -> set，如果这个key在string中存在，那么只输出sting，如果不存在，那么则输出hash的，依次类推
func (c *client) Type(key string) (*keyInfo, error) {
	keyType, err := redis.String(c.conn.Do("TYPE", key))
	if err != nil {
		return nil, err
	}

	info := &keyInfo{keyType: keyType}
	switch keyType {
	case keyTypeNone:
		return nil, errNotFound
	case keyTypeString:
		if size, err := redis.Int64(c.conn.Do("STRLEN", key)); err == nil {
			info.size = float64(size)
		}
	case keyTypeList:
		if size, err := redis.Int64(c.conn.Do("LLEN", key)); err == nil {
			info.size = float64(size)
		}
	case keyTypeSet:
		if size, err := redis.Int64(c.conn.Do("SCARD", key)); err == nil {
			info.size = float64(size)
		}
	case keyTypeZSet:
		if size, err := redis.Int64(c.conn.Do("ZCARD", key)); err == nil {
			info.size = float64(size)
		}
	case keyTypeHash:
		if size, err := redis.Int64(c.conn.Do("HLEN", key)); err == nil {
			info.size = float64(size)
		}
	default:
		return nil, fmt.Errorf("unknown type: %v for key: %v", info.keyType, key)
	}

	return info, nil
}

func (c *client) Get(key string) (string, error) {
	return redis.String(c.conn.Do("GET", key))
}
