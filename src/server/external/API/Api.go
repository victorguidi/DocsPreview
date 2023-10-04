package api

import (
	sqlite "docpreview/external/SQLite"
	shared "docpreview/external/shared"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	Listenaddr string
	Server     *fiber.App
	Sql        *sqlite.SqliteStore
	PrivateLLM *shared.PrivateLLM
}

func NewAPI(listenaddr string) *API {
	api := &API{
		Listenaddr: listenaddr,
		Server: fiber.New(fiber.Config{
			StrictRouting: true,
			AppName:       "North API",
		}),
		Sql:        sqlite.Init(),
		PrivateLLM: shared.NewPrivateLLM(),
	}
	return api
}

func (a *API) UpgradeWss(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (a *API) HandleWss(c *websocket.Conn) {

	type Request struct {
		Prompt string `json:"prompt"`
		Type   string `json:"type"`
		Gpt    bool   `json:"gpt"`
	}

	var (
		mt  int
		msg []byte
		err error
		req Request
	)

	sessionId := a.Sql.CreateSession(c.Query("session"))

	for {

		err = c.ReadJSON(&req)
		if err != nil {
			log.Fatal(err)
		}

		if req.Gpt {
			go a.PrivateLLM.PromptForBulletPoint(req.Prompt, req.Type, c)
		}

		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		fmt.Println(msg)

		if err = c.WriteMessage(mt, []byte(sessionId)); err != nil {
			break
		}
	}
}

func (a *API) Listen() {
	log.Fatal(a.Server.Listen(a.Listenaddr))
}
