[![LICENCE](https://img.shields.io/github/license/HDRUK/gateway-web)](https://github.com/HDRUK/gateway-web/blob/dev/LICENSE)
[![Support](https://img.shields.io/badge/Supported%20By-HDR%20UK-blue)](https://hdruk.ac.uk)

# gateway-metadata-integrations

Welcome to the HDR UK Gateway Metadata Integrations, a **Go** application that allows integration into the Gateway datasets **HDR UK Gateway**. This application facilitates seamless interaction with data custodian endpoints, whether authorized or unauthorized, and passes the data into the Gateway API to be translated into the latest schema version before being added as a dataset in the Gateway.

- [Gateway-Web](https://github.com/HDRUK/gateway-web)
- [Gateway-API](https://github.com/HDRUK/gateway-api)
- [Traser](https://github.com/HDRUK/traser)

This project utilizes **Go** as the language of choice for its robust capabilities and memory-safe stability.

## 🚀 Getting Started

Follow these steps to run the project on your local machine.

### Prerequisites
Ensure you have the following installed:
- **Go** (Latest version) – More info [here](https://go.dev/doc/install)

### Installation & Setup

#### 1️⃣ Clone the repository
Open a terminal and execute:
```bash
git clone https://github.com/HDRUK/gateway-metadata-integrations
```

#### 2️⃣ Navigate to the project directory
```bash
cd gateway-metadata-integrations
```

#### 3️⃣ Set up environment variables
Create a `.env` file and populate it with appropriate values, using `.env.example` as a reference.

#### 4️⃣ Install dependencies
```bash
go mod tidy
```

#### 5️⃣ Start the application
Run the application using:
```bash
go run main.go
```

## 🛠 Available Commands

In the project directory, you can run the following commands:

- **`go run main.go`** – Starts the application.
- **`go build`** – Builds the application for production.
- **`go test ./...`** – Runs the Go test suite.

## 📂 Project Structure
A brief overview of the project's folder structure:
```

├── pkg/pull/          # Pull methods
├── pkg/push/          # Push methods
├── pkg/routes/        # Routing methods
├── pkg/secrets/       # Secret methods   ...shhh..
├── pkg/utils/         # Common utils and mocks
├── pkg/validator/     # Validation methods
├── tests/             # Unit tests
├── .env.example       # Sample environment variables file
├── go.mod             # Go module dependencies
├── go.sum             # Checksums for dependencies
├── main.go            # Application entry point
└── README.md          # Project documentation
```

## 🧪 Testing

We use Go's built-in `testing` package for unit and integration testing.

To run tests:
```bash
go test ./tests
```

## 📖 Additional Resources
- [Go Documentation](https://golang.org/doc/)
- [HDR UK Gateway API](https://github.com/HDRUK/gateway-api-2)

---

For further support, please reach out via [HDR UK](https://healthdatagateway.org/en) or raise a [bug](https://hdruk.atlassian.net/servicedesk/customer/portal/7/group/14/create/34) or even better, submit a PR!

