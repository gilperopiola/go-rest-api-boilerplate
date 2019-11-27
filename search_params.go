package main

//UserSearchParameters are used in the User.Search endpoint
type UserSearchParameters struct {
	FilterID        string
	FilterEmail     string
	FilterFirstName string
	FilterLastName  string

	SortField     string
	SortDirection string

	Limit  int
	Offset int
}
