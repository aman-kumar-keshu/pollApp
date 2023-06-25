func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	// Define the HTTP routes
	e.GET("/polls", func(c echo.Context) error { 
		return c.JSON(200, "GET Polls") 
	})

}

