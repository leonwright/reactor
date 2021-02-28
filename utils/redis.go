package utils

import (
	"context"
	"fmt"
	"time"

	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/leonwright/reactor/logger"
)

var ctx context.Context = context.Background()

const managementAPITokenKey = "AUTH0_MANAGEMENTAPI_TOKEN"
const githubAPIToken = "GITHUB_API_TOKEN"

func rClient() *redis.Client {
	var cfg Config
	ReadFile(&cfg)
	ReadEnv(&cfg)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func CheckForManagementAPIToken() bool {
	deb.Info("Enter method CheckForManagementAPIToken()")
	_, err := GetManagementAPIToken()
	if err != nil {
		deb.Info("Exit method CheckForManagementAPIToken() with value false")
		return false
	}
	deb.Info("Exit method CheckForManagementAPIToken() with value true")
	return true
}

func GetManagementAPIToken() (string, error) {
	deb.Info("Enter method GetManagementAPIToken()")
	client := rClient()
	val, err := client.Get(ctx, managementAPITokenKey).Result()
	if err == redis.Nil {
		deb.Error(err)
		deb.Info("Exit method GetManagementAPIToken() with error.")
		return "", errors.New("token doesn't exist")
	}
	if err != nil {
		deb.Error(err)
		deb.Info("Exit method GetManagementAPIToken() with error.")
		return "", err
	}

	deb.Infof("Successfully retrieved token %s", logger.TruncateString(val, 40))

	deb.Info("Exit method GetManagementAPIToken()")
	return val, nil
}

func UpdateManagementAPIToken(token string) {
	client := rClient()
	err := client.Set(ctx, managementAPITokenKey, token, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("There was a problem updating the API token.")
	}
}

func UpdateGithubToken(username string, token string) {
	client := rClient()
	err := client.Set(ctx, githubAPIToken+":"+username, token, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("There was a problem updating the API")
	}
}

func GetGithubToken(username string) (string, error) {
	client := rClient()
	val, err := client.Get(ctx, githubAPIToken+":"+username).Result()
	return val, err
}

// func ExampleClient() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "redis:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	err := rdb.Set(ctx, "key", "value", 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	val, err := rdb.Get(ctx, "key").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key", val)

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// }
