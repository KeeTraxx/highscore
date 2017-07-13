package main

import (
	"net/http"

	"strconv"

	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Base struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

type Game struct {
	Base
	PlatformID  uint       `sql:"index" json:"platform_id,omitempty"`
	Platform    Platform   `json:"platform"`
	UpdatedByID *uint      `sql:"index" json:"updated_by_id,omitempty"`
	UpdatedBy   *User      `json:"updated_by,omitempty"`
	Names       []GameName `json:"names"`
	Scores      []Score    `json:"scores,omitempty"`
	Score       *uint      `gorm:"-" json:"score,omitempty"`
	Rank        *uint      `gorm:"-" json:"rank,omitempty"`
}

type GameName struct {
	Base
	Game   *Game  `json:"game,omitempty"`
	GameID *uint  `gorm:"unique_index:idx_name_code" sql:"index" json:"game_id,omitempty"`
	Name   string `gorm:"unique_index:idx_name_code" sql:"index" json:"name"`
}

type Platform struct {
	Base
	Name  string `sql:"index" json:"name"`
	Games []Game `json:"games,omitempty"`
}

type User struct {
	Base       `json:"-"`
	ProviderID string  `sql:"index" json:"-"`
	Type       string  `json:"-"`
	Name       string  `sql:"index" json:"name"`
	Email      string  `sql:"index" json:"-"`
	Picture    *string `json:"picture"`
	Password   *string `json:"-"`
}

type Score struct {
	Base

	UserID uint `sql:"index" json:"user_id"`
	User   User `json:"user"`

	GameID uint  `json:"game_id"`
	Game   *Game `json:"game,omitempty"`
	Score  uint  `sql:"index" json:"score"`
}

func initAPI(e *echo.Group) {
	DB.AutoMigrate(&Game{}, &Platform{}, &User{}, &GameName{}, &Score{})

	e.GET("/games", getGames)
	e.GET("/games/:id", getGame)
	e.POST("/scores", postScore, authRequired())
	e.GET("/platforms", getPlatforms)
	e.POST("/games", postGame, authRequired())
	e.PATCH("/games/:id", postGame, authRequired())
	e.PATCH("/settings", patchSettings, authRequired())

	var p Platform

	DB.FirstOrCreate(&Platform{}, &Platform{Name: "NES / Famicom"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sega Master System / Sega Mark III"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Game Boy"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "GameGear"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Megadrive / Genesis"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Neo-Geo"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sega 32X"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sega CD"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "SNES / Super Famicom"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "PC-Engine / TurboGrafx-16"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "3DO"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Amiga CD32"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Jaguar"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo 64"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sony Playstation"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sega Saturn"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Virtual Boy"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "WonderSwan"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Dreamcast"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Game Body Advance"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo GameCube"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "NeoGeo Pocket"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sony Playstation 2"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sony Playstation 3"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sony Playstation 4"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Microsoft Xbox"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo DS"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Sony PSP"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo Wii"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Microsoft Xbox 360"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Microsoft Xbox One"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "PC"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo 3DS"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo Switch"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Nintendo WiiU"})
	DB.FirstOrCreate(&Platform{}, &Platform{Name: "Arcade"})

	DB.FirstOrCreate(&p, &Platform{Name: "Dreamcast"})

	/*g := new(Game)
	g.Names = append(g.Names, GameName{Name: "Testerli"})
	g.Names = append(g.Names, GameName{Name: "Alt Testgame"})
	r := rand.New(rand.NewSource(99))
	for i := 0; i < 20; i++ {
		g.Scores = append(g.Scores, Score{Score: uint(r.Intn(2500) + 500), UserID: 1})
	}
	g.Platform = p

	DB.Create(&g)*/
}

func patchSettings(c echo.Context) error {
	var userSettings User
	c.Bind(&userSettings)
	sess, _ := sessionStore.Get(c.Request(), sessionName)
	p := sess.Values["profile"].(*Profile)

	var user User

	DB.Find(&user, p.ID)

	user.Picture = userSettings.Picture
	user.Name = userSettings.Name

	DB.Save(&user)

	return c.NoContent(http.StatusOK)
}

func postGame(c echo.Context) error {

	sess, _ := sessionStore.Get(c.Request(), sessionName)

	var game Game
	err := c.Bind(&game)

	if err != nil {
		Error.Println(err)
	}

	Info.Printf("Got %+v", &game)

	//id := game.ID
	DB.FirstOrInit(&game, game.ID)

	p := sess.Values["profile"].(*Profile)

	game.UpdatedByID = &p.ID

	DB.Save(&game)

	if DB.Error != nil {
		Error.Println(DB.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &game)
}

func getGames(c echo.Context) error {
	var games []Game

	DB.LogMode(true)

	sess, errSess := sessionStore.Get(c.Request(), sessionName)
	if errSess != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if sess.Values["profile"] != nil {

		userID := sess.Values["profile"].(*Profile).ID

		DB.
			Preload("Names").
			Preload("Platform").
			Preload("Scores", func(db *gorm.DB) *gorm.DB {
				return db.
					Select("*, MAX(score)").
					Order("scores.score DESC").
					Group("scores.user_id, scores.game_id")
			}).
			Preload("Scores.User").
			Select(`
			games.*,
			(SELECT scores.score FROM scores WHERE game_id=games.id AND user_id = ? ORDER BY scores.score DESC LIMIT 1) as score,
			(SELECT COUNT(*) FROM scores WHERE game_id=games.id 
				AND scores.score > (SELECT scores.score FROM scores WHERE game_id=games.id AND user_id = ? ORDER BY scores.score DESC ) )+1 as rank`, userID, userID).
			Joins("INNER JOIN game_names ON game_names.game_id = games.id AND game_names.name LIKE ?", "%"+c.QueryParam("name")+"%").
			Group("games.id").
			Limit(100).
			Find(&games)

	} else {
		DB.
			Preload("Names").
			Preload("Platform").
			Preload("Scores", func(db *gorm.DB) *gorm.DB {
				return db.
					Select("*, MAX(score)").
					Order("scores.score DESC").
					Group("scores.user_id, scores.game_id")
			}).
			Preload("Scores.User").
			Joins("INNER JOIN game_names ON game_names.game_id = games.id AND game_names.name LIKE ?", "%"+c.QueryParam("name")+"%").
			Group("games.id").
			Limit(100).
			Find(&games)
	}

	Info.Printf("%+v", games)

	if DB.Error != nil {
		Error.Println(DB.Error.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, games)
}

func getGame(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var game Game
	DB.
		Preload("UpdatedBy").
		Preload("Platform").
		Preload("Names").
		Preload("Scores", func(db *gorm.DB) *gorm.DB {
			return db.Order("scores.score DESC").Limit(10)
		}).
		Preload("Scores.User").
		First(&game, uint(id))

	if DB.Error != nil {
		Error.Println(DB.Error.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &game)

}

func postScore(c echo.Context) error {
	var score Score
	err := c.Bind(&score)

	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	sess, errSess := sessionStore.Get(c.Request(), sessionName)
	if errSess != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	score.UserID = sess.Values["profile"].(*Profile).ID

	DB.Save(&score)

	var game Game

	DB.Select(`
		games.*,
		(SELECT scores.score FROM scores WHERE game_id=? AND user_id = ? ORDER BY scores.score DESC LIMIT 1) as score,
		(SELECT COUNT(*) FROM scores WHERE game_id = ? 
		AND scores.score > (SELECT scores.score FROM scores WHERE game_id = ? AND user_id = ? ORDER BY scores.score DESC ) )+1 as rank`,
		score.GameID,
		score.UserID,
		score.GameID,
		score.GameID,
		score.UserID,
	).First(&game)

	return c.JSON(http.StatusOK, &game)
}

func getPlatforms(c echo.Context) error {
	var platforms []Platform

	DB.
		Order("name ASC").
		Find(&platforms)

	if DB.Error != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &platforms)

}
