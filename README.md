Product Inventory Management System
=======================================


## Preqs
1) go version >= 1.22
```bash
☁  product-inventory-management-system [master] ⚡  go version
go version go1.22.0 linux/amd64
☁  product-inventory-management-system [master] ⚡
```
2) docker-compose
```bash
☁  product-inventory-management-system [master] ⚡  docker-compose --version
docker-compose version 1.29.2, build unknown
☁  product-inventory-management-system [master] ⚡  
```

## Running the application

Clone the repo
```bash
git clone git@github.com:happilymarrieddad/product-inventory-management-system.git
```

cd into project
```bash
cd product-inventory-management-system
```

Ensure deps are working
```bash
go mod tidy
```

Now we need to startup the database.
```bash
cd db
docker-compose up --build
```

Let's install the tools we need and migrate the database in a new terminal
```bash
// in root project directory
make install.deps
make db.migrate.up
```

It should look like this
```bash
☁  product-inventory-management-system [master] ⚡  make db.migrate.up  
goose -allow-missing -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/product_inventory_management_system?connect_timeout=180&sslmode=disable" up
2024/05/16 10:34:35 OK   20240515232947_add_products.sql (81.64ms)
2024/05/16 10:34:35 goose: successfully migrated database to version: 20240515232947
☁  product-inventory-management-system [master] ⚡  
```

Copy the config sample file to a config file in cmd and fill it out (I could have done this many ways but I felt this was the simpliest)
```bash
cp cmd/config.sample.yaml cmd/config.yaml
```

Time to start the api
```bash
☁  product-inventory-management-system [master] ⚡  go run cmd/main.go 
INFO[05-16|10:35:32] Server running                           api=nil LOG15_ERROR="Normalized odd number of arguments by adding nil" port=9090 versions=[v1]
```

## Using the API
The API has 5 endpoints. Here are some examples of using the endpoints

#### Create
```bash
☁  product-inventory-management-system [master] ⚡  curl -X POST -d '{"name":"nick-test-1","sku":"123532545","qty":50}' -H 'ContextType:application/json' localhost:9090/v1/products
{"id":1,"name":"nick-test-1","sku":"123532545","qty":50,"createdAt":"2024-05-16T10:36:42.100677338-06:00","updatedAt":null}%                         ☁  product-inventory-management-system [master] ⚡  
```

#### Get
```bash
☁  product-inventory-management-system [master] ⚡  curl localhost:9090/v1/products/1
{"id":1,"name":"nick-test-1","sku":"123532545","qty":50,"createdAt":"2024-05-16T04:36:42-06:00","updatedAt":null}%                                   ☁  product-inventory-management-system [master] ⚡  
```

#### Find
```bash
☁  product-inventory-management-system [master] ⚡  curl localhost:9090/v1/products  
{"data":[{"id":1,"name":"nick-test-1","sku":"123532545","qty":50,"createdAt":"2024-05-16T04:36:42-06:00","updatedAt":null}],"count":1}%              ☁  product-inventory-management-system [master] ⚡ 
```

#### Update
```bash
☁  product-inventory-management-system [master] ⚡  curl -X PUT -d '{"name":"new-name","sku":"999","qty":25}' -H 'ContextType:application/json' localhost:9090/v1/products/1  
{"id":1,"name":"new-name","sku":"999","qty":25,"createdAt":"2024-05-16T04:36:42-06:00","updatedAt":"2024-05-16T10:38:34.241864654-06:00"}%           ☁  product-inventory-management-system [master] ⚡  curl localhost:9090/v1/products/1
{"id":1,"name":"new-name","sku":"999","qty":25,"createdAt":"2024-05-15T22:36:42-06:00","updatedAt":"2024-05-16T04:38:34-06:00"}%                     ☁  product-inventory-management-system [master] ⚡  
```

#### Delete
```bash
☁  product-inventory-management-system [master] ⚡  curl -X DELETE localhost:9090/v1/products/1    
☁  product-inventory-management-system [master] ⚡  curl localhost:9090/v1/products            
{"data":[],"count":0}%                                                                                                                     ☁  product-inventory-management-system [master] ⚡  
```

## Testing
I have included integration and unit tests. You can run them by doing the following
```bash
☁  product-inventory-management-system [master] ⚡  ginkgo -r --fail-fast --cover
[1715877615] Products Suite - 23/23 specs DBUG[05-16|10:40:16] unable to get global repo from context   /v1/products=nil LOG15_ERROR="Normalized odd number of arguments by adding nil" requestId=2794a69c-642c-46fd-a375-fe6ee9c56304
...
• SUCCESS! 423.184929ms PASS
coverage: 78.2% of statements
composite coverage: 81.1% of statements

Ginkgo ran 2 suites in 2.676793751s
Test Suite Passed
☁  product-inventory-management-system [master] ⚡  
```
