# Build the application
all: build test

build:
	@echo "Building..."


	@go build -o main cmd/app/main.go

# Run the application
run:
	@go run cmd/app/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi
create-patch-tag:
	@echo "Creating tag..."
	@last_tag=$$(git tag --list | sort --reverse | head -n 1); \
	echo "Last tag is: $$last_tag"; \
	new_tag=$$(echo $$last_tag | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}'); \
	echo "New tag is: $$new_tag"; \
	git tag $$new_tag \
	git push origin $$new_tag

create-minor-tag:
	@echo "Creating tag..."
	@last_tag=$$(git tag --list | sort --reverse | head -n 1); \
	echo "Last tag is: $$last_tag"; \
	new_tag=$$(echo $$last_tag | awk -F. '{printf "v%d.%d.%d", $$1, $$2+1, $$3}'); \
	echo "New tag is: $$new_tag"; \
	git tag $$new_tag \
	git push origin $$new_tag

create-major-tag:
	@echo "Creating tag..."
	@last_tag=$$(git tag --list | sort --reverse | head -n 1); \
	echo "Last tag is: $$last_tag"; \
	new_tag=$$(echo $$last_tag | awk -F. '{printf "v%d.%d.%d", $$1+1, $$2, $$3}'); \
	echo "New tag is: $$new_tag"; \
	git tag $$new_tag \
	git push origin $$new_tag

.PHONY: all build run test clean watch create-patch-tag create-minor-tag create-major-tag
