package bootstrap

import "github.com/amitshekhariitbhu/go-backend-clean-architecture/mongo"

type Application struct {
	Env   *Env
	Mongo mongo.Client
}

func App() (*Application, error) {
	env, err := NewEnv()
	if err != nil {
		return nil, err
	}

	mongo, err := NewMongoDatabase(env)
	if err != nil {
		return nil, err
	}

	app := &Application{
		Env:   env,
		Mongo: mongo,
	}

	return app, nil
}

func (app Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
