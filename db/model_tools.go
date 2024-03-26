package db

type Models struct {
	autoMigrateModelList []interface{}
}

func (ms *Models) Add(m ...interface{}) {
	ms.autoMigrateModelList = append(ms.autoMigrateModelList, m...)
}

func (ms *Models) modelListToAutoMigrate() []interface{} {
	return ms.autoMigrateModelList
}

var modelsToMigrate = Models{}
