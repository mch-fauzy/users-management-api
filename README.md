# User Management System API

This is an API for managing user profiles and authentication, built using Golang and MySQL. It provides endpoints for user registration, login, updating profiles, and retrieving profiles with role-based access control. JWT is used for authentication and authorization.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Endpoints](#endpoints)
- [Postman Collection and Testing](#postman-collection-and-testing)

## Features

- User registration with username uniqueness validation.
- User login with JWT generation and validation.
- Update user profiles with validation and formatting.
- Retrieve user profiles with role-based access control.
- Delete a user by ID with role-based access control.
- Retrieve users with filtering and pagination.

## Installation

1. Clone the repository and navigate to the root folder:

```
git clone https://github.com/mucha-fauzy/users-management-crud-api.git
cd users-management-crud-api
```

2. Install the required dependencies.

3. Create `.env` and set your MySQL or other DB configurations. Refer to `./infras/mysql.go` for the required parameters.

4. Set up the database tables by running the migration scripts `go run ./migrations/users_table.go`.

5. Seed the Admin ID for testing `go run ./seeders/auth_seed.go`.

```
username = "admin_fauzy"
password = "passwordkuat"
```

6. Generate the necessary wire code:

```
go generate ./...
```

7. Build the application:

```
go build
```

8. Run the application:

```
go run .
```

The API will be accessible at http://localhost:8080.

## Endpoints

### Register

Send a POST request to `/v1/auth/register` with a JSON payload containing the registration details.

### Login

Send a POST request to `/v1/auth/login` with a JSON payload containing the login credentials. The response will contain a JWT token.

### Update Own Profile

Send a PATCH request to `/v1/profiles` with a JSON payload containing the profile fields to update. Requires a valid JWT for authentication.

### Retrieve Own Profile

Send a GET request to `/v1/profiles` to retrieve the profile of the authenticated user. Requires a valid JWT for authentication.

### Admin Get Users 

Send a GET request to `/v1/users` to retrieve the users data. Requires a valid JWT with admin role for authentication. You can use the following query parameters for filtering and pagination:

* name: Filter users by name (optional).
* city: Filter users by city (optional).
* province: Filter users by province (optional).
* jobRole: Filter users by job role (optional).
* status: Filter users by status (optional).
* page: Page number for pagination (optional, default: 1).
* size: Number of items per page (optional, default: 5).

The response will include a paginated list of users. The pagination information will be provided in the response body. Here's what the pagination information means:

Response Headers:
* `Total Data`: The total number of users that match the filter criteria.
* `Total Pages`: The total number of pages based on the provided `size`.
* `Current Page`: The current page number.
* `Next Page`: The page number for the next page (if available).
* `Previous Page`: The page number for the previous page (if available).


### Admin Delete User

Send a DELETE request to `/v1/users/{user_id}` to delete a user. Requires a valid JWT with admin role for authentication. Users with the 'admin' role cannot be deleted.

## Postman Collection and Testing

To facilitate testing and interacting with the Users Management API, I provide a Postman collection named `Users-Management-API.postman_collection.json`. This collection includes a set of pre-configured requests that you can use to test various API endpoints easily.

### Import Postman Collection

1. Download the `Users-Management-API.postman_collection.json` file from this repository.
2. Open Postman and click on the "Import" button in the top-left corner.
3. Select the downloaded JSON file and import it into Postman.

### Running Test Scripts

For each request in the Postman collection, I included test scripts to validate the response and ensure that the API endpoints are functioning correctly.
To run the test scripts:

1. Open the imported Postman collection.
2. Select the request you want to test.
3. Click the "Send" button to make the API request.
4. Postman will automatically execute the test scripts and display the results.

Feel free to use the Postman collection to explore the API's capabilities and verify its functionality. The included test scripts help ensure that your API is working as expected.

## Note

This README provides a basic overview of the API and its features. Please refer to the source code for more detailed information and implementation details.