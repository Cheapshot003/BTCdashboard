package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PriceResponse struct {
	Time int64 `json:"time"`
	USD  int64 `json:"USD"`
	EUR  int64 `json:"EUR"`
	GBP  int64 `json:"GBP"`
	CAD  int64 `json:"CAD"`
	CHF  int64 `json:"CHF"`
	AUD  int64 `json:"AUD"`
	JPY  int64 `json:"JPY"`
}

var i int

func get_blockheight() int {
	var apiUrl string = "https://mempool.space/api/blocks/tip/height"

	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	blockheight, err := strconv.Atoi(string(body))
	if err != nil {
		return 0
	}
	return blockheight
}

func get_price() int {
	var apiUrl string = "https://mempool.space/api/v1/prices"

	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	var priceResponse PriceResponse
	if err := json.Unmarshal(body, &priceResponse); err != nil {
		return 0
	}
	return (int(priceResponse.USD))
}

func blockheight_handler(c *fiber.Ctx) error {
	var blockheight int = get_blockheight()
	if blockheight == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Error lol")
	}
	return c.SendString(strconv.Itoa(get_blockheight()))

}

func price_handler(c *fiber.Ctx) error {
	var price int = get_price()
	if price == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Error lol")
	}
	return c.SendString(strconv.Itoa(price) + "$")
}

func satsperdollar_handler(c *fiber.Ctx) error {
	var price int = get_price()
	if price == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Error lol")
	}
	var satsperdollar int = 100000000 / price
	return c.SendString(strconv.Itoa(satsperdollar))
}

func blockstohalving_handler(c *fiber.Ctx) error {
	var blockheight int = get_blockheight()
	if blockheight == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Error lol")
	}
	var blockstohalving int = 840000 - blockheight
	return c.SendString(strconv.Itoa(blockstohalving))
}
func main() {
	app := fiber.New()
	app.Static("/", "./static")

	app.Get("/api/v1/blockheight", blockheight_handler)
	app.Get("/api/v1/price", price_handler)
	app.Get("/api/v1/satsperdollar", satsperdollar_handler)
	app.Get("/api/v1/blockstohalving", blockstohalving_handler)
	app.Listen(":3001")
}
