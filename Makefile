SHELL := /bin/bash

# Directory where migration files live
MIGRATE_DIR := cmd/migrate/migrations
# Default DB URL (override by setting DB_URL environment variable)
DB_URL ?= postgres://admin:admin123@localhost:5432/social?sslmode=disable

.PHONY: migration
migration:
	@name=$(firstword $(filter-out $@,$(MAKECMDGOALS))); \
	if [ -z "$$name" ]; then \
		echo "Usage: make migration <name>"; exit 1; \
	fi; \
	mkdir -p $(MIGRATE_DIR); \
	# Find highest existing 6-digit migration prefix (files like 000001_name.up.sql)
	max=$$(ls $(MIGRATE_DIR) 2>/dev/null | sed -n 's/^\([0-9]\{6\}\)_.*/\1/p' | sort -n | tail -n1); \
	if [ -z "$$max" ]; then \
		next=1; \
	else \
		next=$$((10#$$max + 1)); \
	fi; \
	prefix=$$(printf "%06d" $$next); \
	upfile=$(MIGRATE_DIR)/$${prefix}_$$name.up.sql; \
	downfile=$(MIGRATE_DIR)/$${prefix}_$$name.down.sql; \
	echo "-- +migrate Up" > $$upfile; \
	echo "-- +migrate Down" > $$downfile; \
	echo "Created migration files: $$upfile $$downfile"

# Silence unknown targets (so passing the migration name as a target doesn't error)
%:
	@:
