# api

[![Go Report Card](https://goreportcard.com/badge/github.com/indrasaputra/orvosi-api)](https://goreportcard.com/report/github.com/indrasaputra/orvosi-api)
[![Workflow](https://github.com/indrasaputra/orvosi-api/workflows/Test/badge.svg)](https://github.com/indrasaputra/orvosi-api/actions)
[![codecov](https://codecov.io/gh/indrasaputra/orvosi-api/branch/main/graph/badge.svg?token=HM45WCWOLW)](https://codecov.io/gh/indrasaputra/orvosi-api)
[![Maintainability](https://api.codeclimate.com/v1/badges/2bf28f86e8cecde2563c/maintainability)](https://codeclimate.com/github/indrasaputra/orvosi-api/maintainability)
[![Go Reference](https://pkg.go.dev/badge/github.com/indrasaputra/orvosi-api.svg)](https://pkg.go.dev/github.com/indrasaputra/orvosi-api)

## Description

API provides HTTP REST API for Orvosi application.

## SLI and SLO

- Availability: TBD
- Average response time
    - `POST /sign-in`: TBD
    - `POST /medical-records`: TBD
    - `GET /medical-records`: TBD
    - `GET /medical-records/:id`: TBD
    - `PUT /medical-records/:id`: TBD

## Architecture Diagram

![orvosi-api](https://user-images.githubusercontent.com/4661221/106680454-43908300-65f1-11eb-9f60-c92e900d99f9.png)

## Owner

[Indra Saputra](https://github.com/indrasaputra)

## Onboarding and Development Guide

### How to Run

- Install Go

    We use version 1.15. Follow [Golang installation guideline](https://golang.org/doc/install).

- Install PostgreSQL

    Follow [PostgreSQL installation guideline](https://www.postgresql.org/download/).

- Go to project folder

    Usually, it would be `cd go/src/github.com/indrasaputra/orvosi-api`.

- Run the database

    - Make sure to run PostgreSQL.

- Fill in the environment variables

    Copy the sample env file.
    ```
    cp env.sample .env
    ```
    Then, fill the values according to your setting in `.env` file.

- Download the dependencies

    ```
    make dep-download
    ```
    or run this command if you don't have `make` installed in your local.
    ```
    go mod download 
    ```

- Run the database migration

    Install `golang-migrate`. Follow [Golang Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

    Run the migration command.
    ```
    make migrate url=<database url>
    ```

    e.g:
    ```
    make migrate url=postgres://user:password@localhost:5432/orvosi
    ```

- Run the application

    ```
    go run app/api/main.go
    ```

### Development Guide

- Fork the project

- Create a meaningful branch

    ```
    git checkout -b <your-goal>
    ```
    e.g:
    ```
    git checkout -b strenghten-security-on-sign-in-process
    ```

- Create some changes and their tests (unit test and any test if any).

- Make sure to have unit test coverage at least 90%. There will be times when the code is quite hard to test. Please, explain it in your Pull Request.

- Push the changes to repository.

- Create Pull Request (PR) for your branch. In your PR's description, please explain the goal of the PR and its changes.

- Ask the other contributors to review.

- Once your PR is approved and its pipeline status is green, ask the owner to merge your PR.

## Request Flows

See [Request Flows](https://github.com/indrasaputra/orvosi-api/blob/main/doc/REQUEST-FLOWS.md)

## Endpoints

See [Endpoints](https://github.com/indrasaputra/orvosi-api/blob/main/doc/ENDPOINTS.md)

## Code Map

See [Code Map](https://github.com/indrasaputra/orvosi-api/blob/main/doc/CODE-MAP.md)

## Dependencies

- PostgreSQL