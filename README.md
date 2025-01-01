# EnvBox

EnvBox is a cloud solution to launch environments directly onto the web.

EnvBox provides you web based environments for development. Think Gitpod or Codespaces but 0.1.0 version. It allows you to bring your IDEs on a browser to develop your applications in Golang or Python. 

This is an experimental project I did purely for learning and writing web services in Golang. Not to be used for production. Tested on Google Chrome and Mozilla Firefox.

## Tech Stack

- Backend: Golang, Gin Framework, Docker SDK
- Frontend: HTML, CSS, Bootstrap, Javascript
- Database: SQLite
- Reverse Proxy: Nginx

## Features

- Quick launch of pre-configured environments in a single click
- Write your code in the browser
- Access your environments terminal via web console
- Cleanup environments in a single-click

## Running the application

1. Clone the repository:

   `git clone https://github.com/justsushant/envbox.git`

2. Set the required environment variables in .env file.
   See .env-example for reference.

3. Run the command below to load the docker images for environments and respective database inserts.
   ```
   make migrate
   ```
4. Run the command below to run the application.
   ```
   make run
   ```
5. Connect to the server on public host specified in the .env file using Chrome or Firefox.

## API Endpoints

[Download Postman Collection](docs/envbox.postman_collection.json)

## Testing

- Run the below command to run the tests

  `make test-app`

## Improvements

- See TODOs.md for immediate improvements
- See Upcoming Features section on the Home page
