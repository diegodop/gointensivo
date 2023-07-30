package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/diegodop/gointensivo/internal/infra/database"
	"github.com/diegodop/gointensivo/internal/usecase"
	"github.com/diegodop/gointensivo/pkg/rabbitmq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Car struct {
	Model string
	Color string
}

// metodo
func (c Car) Start() {
	println(c.Model + " has been started")
}

func (c Car) ChangeColor(color string) {
	c.Color = color //duplicando o valor de c.color na memoria - cópia do color original
	println("New color: " + c.Color)
}

func (c *Car) ChangeRealColor(color string) {
	c.Color = color
	println("New color: " + c.Color)
}

// funcao
func soma(x, y int) int {
	return x + y
}

var x string

func main() {

	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	sts := `CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id));`
	_, err = db.Exec(sts)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	orderRepository := database.NewOrderRepository(db)

	uc := usecase.NewCalculateFinalPrice(orderRepository)
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgRabbitMQChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitMQChannel)
	rabbitmqWorker(msgRabbitMQChannel, uc)

}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Starting rabbitmq")
	for msg := range msgChan {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			panic(err)
		}
		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}
		msg.Ack(false)
		fmt.Println("Mensagem processada e salva no banco: ", output)
	}
}
