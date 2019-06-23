all: run
run:
	go run server.go auth.go config.go database.go router.go schema_queries.go users.go utils.go