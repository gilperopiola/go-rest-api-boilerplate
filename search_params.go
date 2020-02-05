package main

import (
	"github.com/gilperopiola/frutils"
	"github.com/gin-gonic/gin"
)

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

func (params *UserSearchParameters) Fill(c *gin.Context) *UserSearchParameters {

	searchParams := &UserSearchParameters{
		FilterID:        c.Query("id"),
		FilterEmail:     c.Query("email"),
		FilterFirstName: c.Query("firstName"),
		FilterLastName:  c.Query("lastName"),

		SortField:     c.Query("sortField"),
		SortDirection: c.Query("sortDirection"),

		Limit:  frutils.ToInt(c.Query("limit")),
		Offset: frutils.ToInt(c.Query("offset")),
	}

	if searchParams.Limit == 0 {
		searchParams.Limit = 99999
	}

	return searchParams
}

func (params *UserSearchParameters) getQueryFormat() []interface{} {
	params.FilterID = "%" + params.FilterID + "%"
	params.FilterEmail = "%" + params.FilterEmail + "%"
	params.FilterFirstName = "%" + params.FilterFirstName + "%"
	params.FilterLastName = "%" + params.FilterLastName + "%"

	return frutils.WrapMultipleValues(params.FilterID, params.FilterEmail, params.FilterFirstName, params.FilterLastName, params.Limit, params.Offset)
}
