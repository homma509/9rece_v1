package registry

import "os"

var envs *Envs

// Envs 環境変数を表した構造体
type Envs struct {
	Cache map[string]string
}

// GetEnvs Envsのインスタンスを取得します
func GetEnvs() *Envs {
	if envs == nil {
		envs = &Envs{
			Cache: map[string]string{},
		}
	}
	return envs
}

func (e *Envs) env(key string) string {
	return os.Getenv(key)
}

// DynamoDBLocalEndpoint DynamoDBのローカルEndpoint
func (e *Envs) DynamoDBLocalEndpoint() string {
	return e.env("DYNAMODB_LOCAL_ENDPOINT")
}

// DynamoTableName DynamoDBのテーブル名
func (e *Envs) DynamoTableName() string {
	return e.env("DYNAMO_TABLE_NAME")
}

// RegionName AWSのリージョン名
func (e *Envs) RegionName() string {
	return e.env("REGION_NAME")
}
