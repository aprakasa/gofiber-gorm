package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Store struct {
		DB *gorm.DB
	}

	Projects struct {
		ID          uint   `gorm:"primaryKey" json:"id"`
		Name        string `gorm:"size:255" json:"name"`
		Description string `gorm:"size:255" json:"description"`
		CreatedAt   *time.Time
		UpdatedAt   *time.Time
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func (s Store) createProject(c *fiber.Ctx) error {

	project := new(Projects)
	if err := c.BodyParser(&project); err != nil {
		return c.Status(500).JSON(Response{"unable to bind data", nil})
	}

	if err := s.DB.Create(&project).Error; err != nil {
		return c.Status(500).JSON(Response{"unable to save project", nil})
	}

	return c.Status(201).JSON(Response{"project created", project})

}

func (s Store) getProjects(c *fiber.Ctx) error {

	var projects []Projects

	if err := s.DB.Find(&projects).Error; err != nil {
		return c.Status(500).JSON(Response{"unable to display projects", nil})
	}

	return c.Status(200).JSON(Response{"project displayed", projects})

}

func (s Store) getProject(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	var project Projects

	err := s.DB.First(&project, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(Response{"project not found", nil})
	} else if err != nil {
		return c.Status(500).JSON(Response{"unable to display project", nil})
	}

	return c.Status(200).JSON(Response{"project created", project})

}

func (s Store) updateProject(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	project := new(Projects)
	if err := c.BodyParser(&project); err != nil {
		return c.JSON(Response{"unable to parse data", nil})
	}

	err := s.DB.First(&project, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(Response{"project not found", id})
	} else if err != nil {
		return c.Status(500).JSON(Response{"unable to display project", nil})
	}

	if err := s.DB.Save(&project).Error; err != nil {
		return c.Status(500).JSON(Response{"unable to update project", nil})
	}

	return c.Status(200).JSON(Response{"project updated", project})

}

func (s Store) deleteProject(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	var project Projects

	if err := s.DB.Delete(&project, id).Error; err != nil {
		return c.Status(500).JSON(Response{"unable to delete project", nil})
	}

	return c.Status(204).JSON(Response{"project deleted", nil})

}

func main() {

	db, err := gorm.Open(sqlite.Open("crud.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Projects{})

	f := fiber.New()
	f.Use(logger.New())
	f.Use(cors.New())

	f.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(Response{"service is up and running", nil})
	})

	app := Store{
		DB: db,
	}

	f.Post("/projects", app.createProject)
	f.Get("/projects", app.getProjects)
	f.Get("/projects/:id", app.getProject)
	f.Patch("/projects/:id", app.updateProject)
	f.Delete("/projects/:id", app.deleteProject)

	log.Fatal(f.Listen("localhost:3333"))

}
