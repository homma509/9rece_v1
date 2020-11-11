package registry

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/homma509/9rece/server/controller"
	"github.com/homma509/9rece/server/infra/db"
	"github.com/homma509/9rece/server/infra/file"
	"github.com/homma509/9rece/server/usecase"
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

// DailyClientPointFile 日別患者行為点数ファイルを生成します
func (f *Factory) DailyClientPointFile() controller.DailyClientPointFile {
	return f.container("DailyClientPointFile", func() interface{} {
		config := &aws.Config{
			Region: aws.String(f.envs.RegionName()),
		}
		return file.NewFile(config)
	}).(controller.DailyClientPointFile)
}

// UkeFile UKEファイルを生成します
func (f *Factory) UkeFile() controller.UkeFile {
	return f.container("UkeFile", func() interface{} {
		config := &aws.Config{
			Region: aws.String(f.envs.RegionName()),
		}
		return file.NewFile(config)
	}).(controller.UkeFile)
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

// DailyClientPointController 日別患者行為点数ハンドラを生成します
func (f *Factory) DailyClientPointController() controller.DailyClientPointController {
	return f.container("DailyClientPointController", func() interface{} {
		return controller.NewDailyClientPointController(
			f.DailyClientPointUsecase(),
			f.DailyClientPointFile(),
		)
	}).(controller.DailyClientPointController)
}

// UkeController UKEハンドラを生成します
func (f *Factory) UkeController() controller.UkeController {
	return f.container("UkeController", func() interface{} {
		return controller.NewUkeController(
			f.UkeFile(),
			f.envs.ServerBucketName(),
		)
	}).(controller.UkeController)
}

// FacilityUsecase 施設ユースケースを生成します
func (f *Factory) FacilityUsecase() usecase.FacilityUsecase {
	return f.container("FacilityUsecase", func() interface{} {
		return usecase.NewFacilityUsecase(
			f.FacilityRepository(),
		)
	}).(usecase.FacilityUsecase)
}

// DailyClientPointUsecase 日別患者行為点数ユースケースを生成します
func (f *Factory) DailyClientPointUsecase() usecase.DailyClientPointUsecase {
	return f.container("DailyClientPointUsecase", func() interface{} {
		return usecase.NewDailyClientPointUsecase(
			f.DailyClientPointRepository(),
		)
	}).(usecase.DailyClientPointUsecase)
}

// FacilityRepository 施設リポジトリを生成します
func (f *Factory) FacilityRepository() *db.FacilityRepository {
	return f.container("FacilityRepository", func() interface{} {
		return db.NewFacilityRepository(
			f.Session(),
		)
	}).(*db.FacilityRepository)
}

// DailyClientPointRepository 日別患者行為点数リポジトリを生成します
func (f *Factory) DailyClientPointRepository() *db.DailyClientPointRepository {
	return f.container("DailyClientPointRepository", func() interface{} {
		return db.NewDailyClientPointRepository(
			f.Session(),
		)
	}).(*db.DailyClientPointRepository)
}
