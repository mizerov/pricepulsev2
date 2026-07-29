module github.com/mizerov/pricepulse/receipts

go 1.26.5

require (
	github.com/caarlos0/env/v11 v11.4.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/joho/godotenv v1.5.1
	github.com/mizerov/pricepulse/api v0.0.0-00010101000000-000000000000
	github.com/mizerov/router v0.0.2
	github.com/samber/do/v2 v2.1.0
	google.golang.org/grpc v1.82.1
)

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/mizerov/pricepulse/api => ../api

replace github.com/mizerov/pricepulse/catalog => ../catalog
