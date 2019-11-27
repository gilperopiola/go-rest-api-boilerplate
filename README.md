# Go Rest API Boilerplate

## Project Architecture
### Main Files

**auth.go** -> Handles token hashing and validation

**database.go** -> Handles database connection and common functions

**schema_queries.go** -> Holds the actual table-creation queries that *database.go* uses

**router.go** -> Sets the API endpoints

**search_params.go** -> Holds the search-parameters struct for every entity

**server.go** -> This project's entrypoint, handles init

**config/config.go** -> This project's configuration, uses the .json files in the same folder

### Entity Dependant Files

**users.go** -> Handles users

**users_roles.go** -> Handles users roles

**users_controllers.go** -> Holds the controller functions for users

**users_tests.go** -> Holds the test functions for users

**users_controllers_tests.go** -> Holds the test functions for users controllers