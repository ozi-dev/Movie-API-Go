# Movie Details API

This project is a simple API that retrieves movie details from [The Movie Database (TMDb)](https://www.themoviedb.org/) and stores them in a MongoDB database. It provides endpoints to fetch and store movie details based on their ID.

## Setup

To run this project, you need to have the following prerequisites installed:

- Go programming language
- MongoDB

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/your-username/movie-details-api.git
   ```

2. Navigate to the project directory:

   ```
   cd movie-details-api
   ```

3. Install the dependencies:

   ```
   go mod download
   ```

4. Update the MongoDB connection URI and API key:

   Open the `main.go` file and locate the following constants:

   ```go
   const (
   	uri       = "mongodb+srv://<your info>@movie-details.attlleo.mongodb.net/?retryWrites=true&w=majority"
   	apiKey    = "<apikey>"
   )
   ```

   Replace the `uri` and `apiKey` values with your own MongoDB connection URI and TMDb API key.

5. Run the application:

   ```
   go run main.go
   ```

   The application will start the server on `http://localhost:8000`.

## Usage

### Fetch Movie Details

To fetch the details of a movie, send a GET request to the `/movies` endpoint with the `id` parameter set to the movie ID.

**Example Request:**

```
[GET /movies?id=12345
](http://localhost:8000/movies?id=12345)```

**Example Response:**

```
Movie found in database:
ID: 12345
Title: Avengers: Endgame
Genres: [Action, Adventure, Science Fiction]
Year: 2019
```

If the movie details are not found in the database, the API will fetch the details from TMDb, store them in the database, and return the response.

### Error Handling

The API handles the following error scenarios:

- If the `id` parameter is missing or empty, it will return a `400 Bad Request` error.
- If there is an error in fetching the movie details from TMDb or reading the response, it will return a `500 Internal Server Error` with the corresponding error message.
- If there is an error in connecting to the MongoDB database or performing database operations, it will return a `500 Internal Server Error` with the corresponding error message.

## Conclusion

This API provides a convenient way to fetch and store movie details using the TMDb API and MongoDB. Feel free to explore and modify the code according to your requirements. If you have any questions or need further assistance, please don't hesitate to reach out.
