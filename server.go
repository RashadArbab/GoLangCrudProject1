package main

import ("github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log") 

	type Todo struct { 
		Id int `json: "id"`
		Name string `json: "name"` 
		Completed bool `json: "is_complete"` 
	}

	var todos = []*Todo {
		{ 1, "task 1" , false}, 
		{ 2, "task 2" , false},
	}


	func main () {
		app := fiber.New() 
		app.Use(logger.New())

		app.Get("/" , func(c *fiber.Ctx) error {
			return c.SendString("Hello world this is air")
		})

		app.Get("/todos" , getTodos)
		app.Post("/todos", createTodos) 

		app.Get("/todo", getSingleTodo)
		app.Delete("/todos" , deleteTodo)
		app.Patch("/todos", updateTodo) 

		app.Use( func(c *fiber.Ctx) error {
			return c.Status(404).SendString("nothing to see here folks")
		})

		err:= app.Listen(":3000")
		if err != nil {
			log.Fatal(err)
		}



	}

	func getTodos(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	}

	func getSingleTodo (c *fiber.Ctx) error {
		type request struct {
			Name string `json: "name"` 
		}
		var body request
		err:= c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannor parse id body" , 
			})
		}
		var todosWithName []Todo 
		for _, todo := range todos {
			if todo.Name == body.Name {
				todosWithName = append(todosWithName , *todo)
			}
		}

		if len(todosWithName) > 0 {
			return c.Status(200).JSON(todosWithName)
		}else {
			return c.Status(fiber.StatusNotFound).SendString("no task found with that name")
		}
	}


	func deleteTodo (c *fiber.Ctx) error { 
		type request struct {
			Name string `json: "name"`
		}

		var body request
		err:= c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Could not parse request")
		}

		for i, todo:= range todos {
			if todo.Name == body.Name {
				todos = append(todos[0:i] , todos[i+1 :]...)
				return c.Status(fiber.StatusNoContent).SendString("todo deleted")	
			}
		}

		return c.Status(fiber.StatusNotFound).SendString("No todo found with that name")
	}


	func createTodos (c *fiber.Ctx) error {

		type request struct {
			Name string `json: "name"`
		}
		var body request

		err := c.BodyParser(&body)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}

		todo := &Todo {
			Id: len(todos) +1 , 
			Name: body.Name, 
			Completed: false,
		}

		todos = append(todos, todo)

		return c.Status(200).SendString("created new todo" + body.Name); 
	}


	func updateTodo (c *fiber.Ctx) error {
		type request struct {
			Name string `json: "name"`
			Completed bool `json: "is_complete"`
		}

		var body request 

		err:= c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request json")
		}
		var selected *Todo 
		for _, todo := range todos {
			if (body.Name == todo.Name){
				selected = todo
				break 
			}
		}

		if selected.Id != 0 {
			selected.Completed = body.Completed
			return c.Status(200).JSON(selected)
		}
		

		return c.Status(fiber.StatusNotFound).SendString("no todo found to update")
	}