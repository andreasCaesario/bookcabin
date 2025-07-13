# Airline Voucher Seat Assignment

This project is a full-stack web application for generating airline crew seat assignment vouchers. It features a Go (Golang 1.23.0) backend and a React frontend, with persistent storage using SQLite. The project is easy to run locally without Docker.

## ⚠️ Important: Clone with Correct Directory Name

Before running the project, make sure the root directory is named `bookcabin-test`.

If you clone the repository and the folder is not named `bookcabin-test`, rename it:

```sh
mv <cloned-folder-name> bookcabin-test
```

Or, when cloning, specify the directory name:

```sh
git clone https://github.com/andreasCaesario/bookcabin.git bookcabin-test
```

This ensures all Go module imports and scripts work correctly.

## Project Structure

```
bookcabin-test/
├── backend/
│   ├── cmd/                # Main entry point for backend server
│   ├── config/             # Database configuration
│   ├── internal/
│   │   ├── domain/         # Domain models (Voucher, Aircraft, etc.)
│   │   ├── handler/        # HTTP handlers (API endpoints)
│   │   ├── repository/     # Data access layer (GORM)
│   │   ├── usecase/        # Business logic
│   │   └── utils/          # Utility functions
│   ├── go.mod, go.sum      # Go dependencies
├── frontend/
│   ├── public/             # Static files
│   ├── src/                # React source code
│   ├── package.json        # Frontend dependencies
├── data/                   # SQLite database volume (persisted)
└── README.md               # Project documentation
```

## General Logic

- **Voucher Generation:**
  - Crew members can generate a voucher for a specific flight and date.
  - The backend ensures only one voucher per (flight number, date) pair (unique constraint).
  - Seats are randomly assigned using aircraft configuration (rows and seat letters).
  - Both backend and frontend validate input and prevent duplicate or invalid requests.

- **API Endpoints:**
  - `POST /api/check` — Check if a voucher exists for a flight/date.
  - `POST /api/generate` — Generate a new voucher (with random seats) if not already assigned.
  - `GET /api/aircraft-list` — List available aircraft types. (use for aircraft list dropdown in frontend)

- **Frontend:**
  - Simple React UI for entering crew info, flight, date, and aircraft.
  - Displays errors and success messages from backend.

## First-time Setup & Running Locally

1. **Backend:**
   ```sh
   cd backend
   go mod tidy
   go run ./cmd/main.go
   ```
   The backend will start on http://localhost:8080

2. **Frontend (open new terminal):**
   ```sh
   cd frontend
   npm install
   npm start
   ```
   The frontend will start on http://localhost:3000

3. **Database:**
   - SQLite file is persisted in the `data/` directory on your host (default: `backend/vouchers.db`).

## Backend Architecture

The backend is developed using the principles of Clean Architecture:

- **Domain Layer:** Core business models and logic (`internal/domain`).
- **Usecase Layer:** Application-specific business rules (`internal/usecase`).
- **Repository Layer:** Data access and persistence, abstracted via interfaces (`internal/repository`).
- **Handler Layer:** HTTP API endpoints and request/response handling (`internal/handler`).
- **Config/Utils:** Configuration and utility helpers.

This separation ensures the codebase is modular, testable, and easy to maintain or extend.

---
