# lions-club-borken-advent-calendar

The application lions-club-borken-advent-calendar sends a daily email to recipients with the winning numbers for the Lions Club Borken Advent calendar.

## Usage

1. Update the `config.yml` file to your needs.
2. Create a `.env` file with a content similar to the following or set the environment variables accordingly:
   ```
   SMTP_HOST=smtp.asdfasdfasdfasdfasdfadsfadsf.de
   SMTP_USERNAME=forename1.surname@asdfasdfasdfasdfasdfadsfadsf.de
   SMTP_PASSWORD=secret
   ```
3. Update the docker setup if necessary.
4. Start the application using docker compose (recommended)
   ```
   docker compose up -d
   ```
   or run it without docker.
   ```
   go run ./cmd/main.go
   ```
