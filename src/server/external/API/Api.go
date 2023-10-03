package api

import (
	"docpreview/external/SQLite"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	Listenaddr string
	Server     *fiber.App
	Sql        *sqlite.SqliteStore
}

func NewAPI(listenaddr string) *API {
	api := &API{
		Listenaddr: listenaddr,
		Server: fiber.New(fiber.Config{
			StrictRouting: true,
			AppName:       "North API",
		}),
		Sql: sqlite.Init(),
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

	var (
		mt  int
		msg []byte
		err error
	)

	sessionId := a.Sql.CreateSession(c.Query("session"))

	log.Println(sessionId)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		log.Println("write:", err)
		if err = c.WriteMessage(mt, []byte(sessionId)); err != nil {
			break
		}
	}

}

func (a *API) Listen() {
	log.Fatal(a.Server.Listen(a.Listenaddr))
}
