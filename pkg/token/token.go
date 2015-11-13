package token

import (
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/generator"
	"github.com/PandoCloud/pando-cloud/pkg/redispool"
	"reflect"
	"strconv"
)

const (
	DeviceTokenKeyPrefix = "device:token:"
	DeviceTokenExpires   = 7200
)

type Helper struct {
	redishost string
}

func NewHelper(host string) *Helper {
	helper := &Helper{
		redishost: host,
	}
	return helper
}

func (helper *Helper) GenerateToken(id uint64) ([]byte, error) {
	token, err := generator.GenRandomToken()
	if err != nil {
		return nil, err
	}

	conn, err := redispool.GetClient(helper.redishost)
	if err != nil {
		return nil, err
	}

	key := DeviceTokenKeyPrefix + strconv.FormatUint(id, 10)

	_, err = conn.Do("SET", key, token)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("EXPIRE", key, DeviceTokenExpires)
	if err != nil {
		return nil, err
	}

	return token, nil

}

func (helper *Helper) ValidateToken(id uint64, token []byte) error {
	key := DeviceTokenKeyPrefix + strconv.FormatUint(id, 10)

	conn, err := redispool.GetClient(helper.redishost)
	if err != nil {
		return err
	}

	readToken, err := conn.Do("GET", key)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(readToken, token) {
		return errors.New("token not match.")
	}

	_, err = conn.Do("EXPIRE", key, DeviceTokenExpires)
	if err != nil {
		return err
	}

	return nil
}

func (helper *Helper) ClearToken(id uint64) error {
	key := DeviceTokenKeyPrefix + strconv.FormatUint(id, 10)

	conn, err := redispool.GetClient(helper.redishost)
	if err != nil {
		return err
	}

	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}
