package main

import (
	"database/sql"
	"fmt"

	"github.com/diegodop/gointensivo/internal/infra/database"
	"github.com/diegodop/gointensivo/internal/usecase"
	_ "github.com/mattn/go-sqlite3"
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

	input := usecase.OrderInput{
		ID:    "1234",
		Price: 10.0,
		Tax:   1.0,
	}

	output, err := uc.Execute(input)
	if err != nil {
		panic(err)
	}
	println(output)
	fmt.Println(output)
	/* x = "Hello World"
	car := Car{ //declarando e atribuindo a variavel car
		Model: "Ferrari",
		Color: "Red",
	}
	car.Model = "Fiat" //atribuindo o valor do atributo Model
	println(car.Model)

	car.Start()

	car.ChangeColor("blue")
	println(car.Color)
	car.ChangeRealColor("blue")
	println(car.Color)

	a := 10
	//b := a //copia o valor de a, mas são independentes
	//b = 20

	println(a)
	//println(b)

	println(&a) //imprime o endereço de memoria onde a está armazenado

	b := &a
	*b = 20

	println(a)
	println(b)

	order, err := entity.NewOrder("1", -10, 1)
	if err != nil {
		println(err.Error())
	} else {
		println(order.ID)
	} */
}
