package main

import (
	"database/sql" 
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"os"
)

//inicializamos nuestras dependencia
 var db *sql.DB
 func initBD(){
	var err error
	connectionString:=os.Getenv("connectionString")
	db,err=sql.Open("sqlserver",connectionString)
	if err !=nil{
		log.Fatal("Error en la conexion en la conexion de base de datos",err.Error())
	}

    //realizamos si la conexion es activa o no
	err=db.Ping();
	//Enviamos el estado de error
	if err!=nil{
       log.Fatal(err)
	}
	fmt.Println("Se ha conectado correctamente a la base de datos")
 }

func main() {
	initBD()
    router:=mux.NewRouter() //inicializamos nuestro router
	router.HandleFunc("/cerveza",getCerveza).Methods("GET")
	router.HandleFunc("/cerveza",postCerveza).Methods("POST")
	log.Fatal(http.ListenAndServe(":6000",router))
}

func getCerveza(w http.ResponseWriter, r *http.Request){
    rows,err:= db.Query("SELECT Nombre_cerveza,Cantidad,Precio FROM productoCerveza")
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	defer rows.Close()//liberamos espacio mediante la funcion
	var cervezas []Cerveza//realizamos un slice de nuestro array para que aumente o disminuya durante el proceso de ejecucion
   for rows.Next(){
       var cerveza Cerveza
	   err:=rows.Scan(&cerveza.Nombre_cerveza,&cerveza.Cantidad,&cerveza.Precio)
	   if err !=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	   }
	   cervezas=append(cervezas,cerveza)
	}
      w.Header().Set("Content-Type","application/json")
	  json.NewEncoder(w).Encode(cervezas)
}
//Realizamos una peticion para enviar los datos
func postCerveza(w http.ResponseWriter,r *http.Request){
	var data Cerveza
	err:=json.NewDecoder(r.Body).Decode(&data)
	 if err !=nil{
	  
		 http.Error(w,err.Error(),http.StatusBadRequest)
		 return
       }
	   if data.Nombre_cerveza== "" || data.Cantidad==0||data.Precio==0{
          http.Error(w,"Datos de productos nulos",http.StatusBadRequest)
		  return
	   }
	   _,err=db.Query("INSERT INTO productoCerveza(Nombre_cerveza,Cantidad,Precio)VALUES(@Nombre_cerveza,@Cantidad,@Precio)",sql.Named("Nombre_cerveza",data.Nombre_cerveza),sql.Named("Cantidad",data.Cantidad),sql.Named("Precio",data.Precio))
	   if err !=nil{
		 http.Error(w,err.Error(),http.StatusInternalServerError)
		 return
	   }
	   w.WriteHeader(http.StatusCreated)
	   w.Write([]byte("Datos del producto enviados correctamente"))
	   }
	   
	  
type Cerveza struct{

     Cantidad int `json:"Cantidad"`
	Nombre_cerveza string `json:"Nombre_cerveza"`
	Precio float64 `json:"Precio"`
}