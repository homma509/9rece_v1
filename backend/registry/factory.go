package registry

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/homma509/9rece/backend/controller"
	"github.com/homma509/9rece/backend/infra/db"
	"github.com/homma509/9rece/backend/infra/file"
	"github.com/homma509/9rece/backend/usecase"
)

var factory *Factory

// Factory インターフェースの詳細を生成する構造体
type Factory struct {
	envs  *Envs
	cache map[string]interface{}
}

// Clear キャッシュをクリアします
func Clear() {
	factory = nil
}

// Creater Factoryのインスタンスを取得します
func Creater() *Factory {
	if factory == nil {
		factory = &Factory{
			envs: GetEnvs(),
		}
	}
	return factory
}

func (f *Factory) container(key string, builder func() interface{}) interface{} {
	if f.cache == nil {
		f.cache = map[string]interface{}{}
	}
	if _, ok := f.cache[key]; !ok {
		f.cache[key] = builder()
	}
	return f.cache[key]
}

// Session DB接続を生成します
func (f *Factory) Session() *db.Session {
	return f.container("Session", func() interface{} {
		config := &aws.Config{
			Region: aws.String(f.envs.RegionName()),
		}
		// if f.Envs.DynamoDBLocalEndpoint() != "" {
		// 	config.Credentials = credentials.NewStaticCredentials("dummy_id", "dummy_secret", "dymmy_token")
		// 	config.Endpoint = aws.String(f.Envs.DynamoDBLocalEndpoint())
		// }
		return db.NewSession(config, f.envs.DynamoTableName())
	}).(*db.Session)
}

// FacilityFile 施設ファイルを生成します
func (f *Factory) FacilityFile() controller.FacilityFile {
	return f.container("FacilityFile", func() interface{} {
		config := &aws.Config{
			Region: aws.String(f.envs.RegionName()),
		}
		return file.NewFile(config)
	}).(controller.FacilityFile)
}

// FacilityController 施設ハンドラを生成します
func (f *Factory) FacilityController() controller.FacilityController {
	return f.container("FacilityController", func() interface{} {
		return controller.NewFacilityController(
			f.FacilityUsecase(),
			f.FacilityFile(),
		)
	}).(controller.FacilityController)
}

// FacilityUsecase 施設ユースケースを生成します
func (f *Factory) FacilityUsecase() usecase.FacilityUsecase {
	return f.container("FacilityUsecase", func() interface{} {
		return usecase.NewFacilityUsecase(
			f.FacilityRepository(),
		)
	}).(usecase.FacilityUsecase)
}

// FacilityRepository 施設リポジトリを生成します
func (f *Factory) FacilityRepository() *db.FacilityRepository {
	return f.container("FacilityRepository", func() interface{} {
		return db.NewFacilityRepository(
			f.Session(),
		)
	}).(*db.FacilityRepository)
}
