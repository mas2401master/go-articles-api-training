# Api Rest con Golang
## Ejercicio de aplicación de SIDESA GROUP
Api rest básica para dar soporte a un carrito de compras.

## Herramientas de construccion
- Go version go1.19.3 windows/amd64
- Gin framework github.com/gin-gonic/gin v1.8.1
- Jwt github.com/golang-jwt/jwt/v4 v4.4.3
- Driver postgresql para Golang github.com/lib/pq v1.10.7
- Viper github.com/spf13/viper v1.14.0
- swaggerFiles "github.com/swaggo/files"
- ginSwagger "github.com/swaggo/gin-swagger"

## Estructura del proyecto
El proyecto se estructuro de la siguente forma:
- /cmd: es la carpeta donde se encuentra el archivo principal del proyecto main.go 
- /data: contiene un archivo .yml donde se configuran las variables de entorno usadas en el proyecto y un archivo .sql donde se encuentra la estructura de la base de datos propuesta.
- /docs: contiene los archivos generados por swagger. El archivo docs.go que instancia  el swagger y los archivos swagger.json y swagger.yaml generados por la herramienta para generar la documentacion de la api. Disponible en http://localhost:3000/docs/index.html
- /internal: contiene la carpeta api, donde se encuentra los controladores, middlewares, router de la api rest. La carpeta pkg, contiene los modelos, repositorios, configuracion de base de datos y variables de entorno.
- /log: contiene el archivo log generado por gin 
- /pkg:contiene paquetes reutilizables en la aplicacion.
- Makefile:  comandos para la ejecucion de la aplicacion mediante comando make.

## Construccion del binario:
```sh
go build -o bin/restapi cmd/main.go
```

## Corrida de la api rest
```sh
go run cmd/main.go
```

## Acceso a la api rest
- Existe el usuario de rol ADMIN registrado con los datos: username:admin, password:1234

## Logeo
- El login comprueba datos de autenticacion enviados por el cliente. 
- La ruta del recurso es http://localhost:3000/api/v1/login
- Se envia por body el json:
```sh
{
    "username":"admin",
    "password":"1234"
}
```
- Genera como respuesta un token de autenticacion del usuario. La respuesta es:
```sh
{
    "X-USERID": "uMYb+qE81Kblo7xujUFq2Ck2bKoUPkrY9ht46qXyUXYgWqAi",
    "firsname": "pedro",
    "lastname": "perez",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZHVzZXIiOiIxMyIsInJvbGUiOiIyIn0.1PqMEY9I1YNmY4AP4Djp034T7IFeC5tLDeNZ42b60G4"
}
```
## Gestión de usuarios
- Permite la creacion y actualizacion de usuarios del api rest.
- La ruta para la creacion de usuarios: localhost:3000/api/v1/users por metodo post
- Para su acceso se debe enviar el Headers el token generado en el logeo del usuario
- Se valida campo username (no exista)
```sh
Authorization eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZHVzZXIiOiI0Iiwicm9sZSI6IjEifQ.kNH4QXz-ial4isH9EVlz4DBtVHhqe-rVR0OI60x3CVE
```
- Se envia por body el json:
```sh
{
  "email": "string",
  "firstname": "string",
  "lastname": "string",
  "password": "string",
  "role_id": 0,
  "status": true,
  "username": "string"
}
```
- La respuesta: 
```sh
{
  "message": "string",
  "status": 0
}
```
- Para la actualizacion localhost:3000/api/v1/users/{id} por metodo put
- Se valida campo username (no exista)
- Se valida autorizacion de usuario
- se valida existencia del id de usuario
- se valida usuario a editar asociado a usuario de session para el rol cliente
- Se envia por body el json:
```sh
{
  "email": "string",
  "firstname": "string",
  "lastname": "string",
  "password": "string",
  "role_id": 0,
  "status": true,
  "username": "string"
}
```

- la respuesta sera los datos del usuario editado:
```sh
{
  "created_at": "string",
  "email": "string",
  "firstname": "string",
  "id": 0,
  "lastname": "string",
  "password": "string",
  "role": {
    "id": 0,
    "name": "string"
  },
  "role_id": 0,
  "status": true,
  "updated_at": "string",
  "username": "string"
}
```
- Para la eliminacion localhost:3000/api/v1/users/{id} por metodo delete
- Se valida campo username (no exista)
- Se valida autorizacion de usuario
- se valida existencia del id de usuario
- se valida que el id de usuario no este asociado a una orden
- la respuesta:
```sh
{
  "message": "string",
  "status": 0
}
```
- Para la Listar localhost:3000/api/v1/users por metodo get
- Se valida autorizacion de usuario enviado en el header.
- Los paramertros de envio pueden ser status, rol (nombre de rol)
- por ejemplo:
```sh
localhost:3000/api/v1/users?role=2&status=false
```
- la respuesta:
```sh
[
  {
    "created_at": "string",
    "email": "string",
    "firstname": "string",
    "id": 0,
    "lastname": "string",
    "password": "string",
    "role": {
      "id": 0,
      "name": "string"
    },
    "role_id": 0,
    "status": true,
    "updated_at": "string",
    "username": "string"
  }
]
```
- Para mostrar informacion de un usuario: localhost:3000/api/v1/users/{id}
- Se valida autorizacion de usuario
- se valida existencia del id de usuario
- se valida usuario a editar asociado a usuario de session para el rol cliente
- La respuesta sera:

```sh
{
  "created_at": "string",
  "email": "string",
  "firstname": "string",
  "id": 0,
  "lastname": "string",
  "password": "string",
  "role": {
    "id": 0,
    "name": "string"
  },
  "role_id": 0,
  "status": true,
  "updated_at": "string",
  "username": "string"
}
```
## Gestion de Ordenes
- Para la creacion localhost:3000/api/v1/order por metodo post
- Se valida autorizacion de usuario
- Se valida que el usuario no tenga una orden registrada en estatus R=Revision
- Se valida la existencia del item a registra para crear la orden.
- se valida la disponibilidad del item
- en caso de incluir promocion, se valida la existencia del codigo enviado, y si esta usado o no ese codigo.
- Esta funcionalidad se realiza para generar la orden del usuario, en caso de querer agregar detalle se usa el recurso correspondiente. Es decir, se crea la orden siempre que el usuario no tenga una orden en estatus Revision.
- Se envia por body el json:
```sh
{
  "code": "string",
  "item_id": 0,
  "price": 0,
  "quantity": 1
}
```
- La respuesta sera:
```sh
{
  "message": "string",
  "status": 0
}
```
- Para la modificacion localhost:3000/api/v1/order/{id} por metodo put
- Se valida autorizacion de usuario
- Se valida la existencia del id de la orden
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- se valida la existencia y la disponibilidad del codigo de promocion, en caso de enviar codigo de la promocion. 
- Se envia por body el json:
```sh
{
  "code": "string",
  "status": "string"
}
```
- la respuesta sera:
```sh
{
  "created_at": "string",
  "detail_order": [
    {
      "created_at": "string",
      "id": 0,
      "item_id": 0,
      "name_item": "string",
      "order_id": 0,
      "price": 0,
      "quantity": 0,
      "status": "string",
      "total": 0,
      "updated_at": "string"
    }
  ],
  "id": 0,
  "order_number": 0,
  "promotion": {
    "code": "string",
    "discount": 1,
    "name": "string"
  },
  "promotion_id": 0,
  "quantity": 0,
  "status": "string",
  "subtotal": 0,
  "total": 0,
  "total_discount": 0,
  "updated_at": "string",
  "user": {
    "firstname": "string",
    "lastname": "string",
    "username": "string"
  },
  "user_id": 0
}
```
- Para la eliminacion localhost:3000/api/v1/order/{id} por metodo delete
- Se valida autorizacion de usuario
- Se valida la existencia del id de la orden
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- se libera items, y promocion asociadas a la orden
- la respuesta sera 
```sh
{
  "message": "string",
  "status": 0
}
```

- Para mostrar una orden localhost:3000/api/v1/order/{id} por el metodo get
- Se valida autorizacion de usuario
- Se valida la existencia del id de la orden
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- la respuesta:
```sh
{
  "created_at": "string",
  "detail_order": [
    {
      "created_at": "string",
      "id": 0,
      "item_id": 0,
      "name_item": "string",
      "order_id": 0,
      "price": 0,
      "quantity": 0,
      "status": "string",
      "total": 0,
      "updated_at": "string"
    }
  ],
  "id": 0,
  "order_number": 0,
  "promotion": {
    "code": "string",
    "discount": 1,
    "name": "string"
  },
  "promotion_id": 0,
  "quantity": 0,
  "status": "string",
  "subtotal": 0,
  "total": 0,
  "total_discount": 0,
  "updated_at": "string",
  "user": {
    "firstname": "string",
    "lastname": "string",
    "username": "string"
  },
  "user_id": 0
}
```
- Para listar las orden localhost:3000/api/v1/order por el metodo get
- Se valida autorizacion de usuario
- los filtros son code, status por ejemplo: localhost:3000/api/v1/orders?code=cyber&status=R
- produce como respuesta listado de ordenes.

# Orden Detalle
- Para la creacion localhost:3000/api/v1/orders/detail/:id  por metodo post
- Se valida autorizacion de usuario
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- Se valida que el usuario no tenga una orden registrada en estatus R=Revision
- Se valida la existencia del id ordena agregar detalle
- se valida la existencia y disponibilidad del item
- Se envia por body el json:
```sh
{
  "item_id": 0,
  "quantity": 1
}
```
- La respuesta sera:
```sh
{
  "message": "string",
  "status": 0
}
```
- Para la modificacion del detalle localhost:3000/api/v1/orders/detail/{id}  por metodo put
- Se valida autorizacion de usuario
- Se valida la existencia del id del detalle a ser agregado a la orden
- se valida la existencia y disponibilidad del item
- Se envia por body el json:
```sh
{
  "item_id": 0,
  "quantity": 1
}
```
- la respuesta sera: 
```sh
{
  "created_at": "string",
  "id": 0,
  "item_id": 0,
  "name_item": "string",
  "order_id": 0,
  "price": 0,
  "quantity": 0,
  "status": "string",
  "total": 0,
  "updated_at": "string"
}
```
- Para mostrar un detalle de la orden  localhost:3000/api/v1/orders/detail/{id}  por metodo get
- Se valida autorizacion de usuario
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- Se valida la existencia del id del detalle a ser mostrado de la orden
- la respuesta sera: 
```sh
{
  "created_at": "string",
  "id": 0,
  "item_id": 0,
  "name_item": "string",
  "order_id": 0,
  "price": 0,
  "quantity": 0,
  "status": "string",
  "total": 0,
  "updated_at": "string"
}
```
- Para la eliminacion de la orden detalle localhost:3000/api/v1/orders/detail/{id}  por metodo get
- Se valida autorizacion de usuario
- Se valida la existencia del id del detalle a ser mostrado de la orden 
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- La respuesta sera:
```sh
{
  "message": "string",
  "status": 0
}
```
- Para listar el detalle de la orden localhost:3000/api/v1/orders/details/{id}  por metodo get
- se recibe el id de la orden a listar su detalle
- Se valida autorizacion de usuario
- se valida existencia de la orden a listar su detalle
- se valida que la orden este asociada al id de usuario, en caso de rol CLIENTE.
- la respuesta sera:
```sh
[
  {
    "created_at": "string",
    "id": 0,
    "item_id": 0,
    "name_item": "string",
    "order_id": 0,
    "price": 0,
    "quantity": 0,
    "status": "string",
    "total": 0,
    "updated_at": "string"
  }
]
```