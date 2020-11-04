package main

import (
	"encoding/json" //PODREMOS CONVERTIR GO TYPES A JSON LLAMADO(MARSHALL/UNMARSHALL)
	"fmt"           //FUNCIONES I/O SIMILAR A PRINTF
	"log"           //LOG PARA VER ERRORES
	"net/http"      //METODO PARA OPERACIONES SOBRE HTTP
	"strconv"       //METODO PARA CONVERSION DE STRING A DISTINTO TIPO DE DATOS
	"time"          //METODO PARA EL VALOR DEL TIEMPO(MANIPULAR)

	"github.com/gorilla/mux"                  //METODO PARA SOLICITUDES HTTP ENTRANTES A SUS RESPECTIVOS HANDLER
	"github.com/jinzhu/gorm"                  //MAPEO DE OBJETO RELACIONAL
	_ "github.com/jinzhu/gorm/dialects/mysql" //DRIVER PARA RECORDAR PATH MYSQL
)

// REPRESENTA EL MODELOr
// NOMBRE POR DEFECTO ORDERS
type Order struct {
	// gorm.Model
	OrderID      uint      `json:"orderId" gorm:"primary_key"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"orderedAt"`
	Items        []Item    `json:"items" gorm:"foreignkey:OrderID"`
}

// Item represents the model for an item in the order
type Item struct {
	// gorm.Model
	LineItemID  uint   `json:"lineItemId" gorm:"primary_key"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"-"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:@tcp(127.0.0.1:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// CREAR LA BD
	db.Exec("CREATE DATABASE orders_db")
	db.Exec("USE orders_db")

	// MIGRACION PARA CREAR TABLA ORDER E ITEM
	db.AutoMigrate(&Order{}, &Item{})
}

func createOrder(w http.ResponseWriter, r *http.Request) { //TODOS LOS HANDLER TIENEN ESTE MISMO FORMATO
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	// CREAR ORDEN INSERTANDO EN LAS TABLAS ORDERS Y ITEM
	db.Create(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func getOrders(w http.ResponseWriter, r *http.Request) { //TODOS LOS HANDLER TIENEN ESTE MISMO FORMATO
	w.Header().Set("Content-Type", "application/json")
	var orders []Order
	db.Preload("Items").Find(&orders)
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) { //TODOS LOS HANDLER TIENEN ESTE MISMO FORMATO
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderID := params["orderId"]

	var order Order
	db.Preload("Items").First(&order, inputOrderID)
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) { //TODOS LOS HANDLER TIENEN ESTE MISMO FORMATO
	var updatedOrder Order
	json.NewDecoder(r.Body).Decode(&updatedOrder)
	db.Save(&updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) { //TODOS LOS HANDLER TIENEN ESTE MISMO FORMATO
	params := mux.Vars(r)
	inputOrderID := params["orderId"]
	// CONVERTIR `orderId` string param to uint64
	id64, _ := strconv.ParseUint(inputOrderID, 10, 64)
	// CONVERTIR UINT64 EN UINT
	idToDelete := uint(id64)

	db.Where("order_id = ?", idToDelete).Delete(&Item{})
	db.Where("order_id = ?", idToDelete).Delete(&Order{})
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter() //UNA RUTA ES UNA FORMA DE ESPECIFICAR QUE FUNCION MANEJA UN API REQUEST
	/*CON EL GORILLAX MUX LAS RUTAS SON DEFINIDAS USANDO EL HANDLEFUNC EL PRIMER ARGUMENTO ES EL API PATH
	Y EL SEGUNDO ARGUMENTO ES EL NOMBRE DEL METODO LA FUNCION DEL METODO HTTP (POST,GET...)*/
	// SE REGISTRAN SUS ROUTES(CREAR,LEER,ACTUALIZAR) ASIGNANDO  RUTAS DE URL A SUS METODOS
	// CREAR
	router.HandleFunc("/orders", createOrder).Methods("POST")
	// LEER
	router.HandleFunc("/orders/{orderId}", getOrder).Methods("GET")
	// LEER TODO
	router.HandleFunc("/orders", getOrders).Methods("GET")
	// UPDATE
	router.HandleFunc("/orders/{orderId}", updateOrder).Methods("PUT")
	// BORRAR
	router.HandleFunc("/orders/{orderId}", deleteOrder).Methods("DELETE")
	// INICIAR CONEXION
	initDB()

	log.Fatal(http.ListenAndServe(":8080", router))
}
