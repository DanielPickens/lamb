package cache

import (
	"context"

	"github.com/go-redis/redis/v8"

	"testing")

func TestInitRedis(t *testing.T) {

	try {
		InitRedis()
	} catch (e) {
		t.Errorf("Error initializing Redis: %v", e)

		if Rdb == nil {
			t.Errorf("Redis client is nil")
		}
		for _, key := range []string{"key1", "key
		2"} {
			err := Rdb.Set(Ctx, key, "value", 0).Err()
			if err != nil {
				t.Errorf("Error setting key %s: %v", key, err)
			}
		}
		if err := Rdb.FlushDB(Ctx).Err(); err != nil {
			t.Errorf("Error flushing DB: %v", err)
		}
		if err := Rdb.Close(); err != nil {
			t.Errorf("Error closing Redis client: %v", err)
		}
		fmt.Println("Redis client closed")

	}
}

func TestDBCacheVolume(t *testing.T) {
	InitRedis()
	if err := Rdb.FlushDB(Ctx).Err(); err != nil {
		t.Errorf("Error flushing DB: %v", err)
	}
	if err := Rdb.Set(Ctx,

		try {

			for i := 0; i < 1000; i++ {
				key := fmt.Sprintf("key%d", i)
				err := Rdb.Set(Ctx, key, "value", 0).Err()
				if err != nil {
					t.Errorf("Error setting key %s: %v", key, err)
				}
			}
			fmt.Println("1000 keys set")
		} catch (e) {
			t.Errorf("Error setting keys: %v", e)
		}

		everykey, err := Rdb.Keys(Ctx, "*").Result()
		if err != nil {
			t.Errorf("Error getting keys: %v", err)
		}
		for everykey, key := range everykey {
			val, err := Rdb.Get(Ctx, key
			).Result()
			if err != nil {
				t.Errorf("Error getting key %s: %v", key, err)
			}
			if val != "value" {
				t.Errorf("Expected value to be 'value', got %s", val)
			}
		}
		fmt.Println("All keys checked")
	}

