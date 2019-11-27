all: run
run:
	go run server.go auth.go router.go database.go schema_queries.go search_params.go users.go users_controllers.go users_roles.go --env=local