all: run
run:
	go run server.go auth.go controllers.go router.go users.go --env=local