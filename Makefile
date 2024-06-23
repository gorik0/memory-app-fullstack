.PHONY: create-keypair
.PHONY: mi-up
.PHONY: create-keypair

PWD = $(shell pwd)
ACCTPATH = $(PWD)
PATH_TO_MIGRATIONS = $(ACCTPATH)/migrations
DATABASE_URL = gorik:123@localhost:5432/postgres?sslmode=disable

N=1
create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCTPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCTPATH)/rsa_private_$(ENV).pem -pubout -out $(ACCTPATH)/rsa_public_$(ENV).pem
mi-create:
	migrate create -ext sql -dir $(PATH_TO_MIGRATIONS) -digits 5 -seq $(NAME)
mi-up:
	migrate -source file://$(PATH_TO_MIGRATIONS) -database postgres://$(DATABASE_URL) up ${N}
mi-down:
	migrate -source file://$(PATH_TO_MIGRATIONS) -database postgres://$(DATABASE_URL) down $(N)
mi-force:
	migrate -source file://$(PATH_TO_MIGRATIONS) -database postgres://$(DATABASE_URL) force $(N)
