package model

import "github.com/jangozw/gin-smart/pkg/app"

func init()  {
	 app.LoadServices(app.DbService)
	// app.LoadDb()

}
