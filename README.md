# url-shortener

Another url shortener

## Usage

Uses postgressql for the database

```
go run main.go -help
  -dev
        hides dev routes (default true)
  -host string
        speicify postgres host (default "localhost")
  -name string
        database name (default "longer")
  -password string
        db user password (default "postgres")
  -port int
        specify postgres port (default 5432)
  -username string
        db user username (default "postgres")
```

## Running

```bash
git clone https://github.com/filipwtf/url-shortener.git
cd /url-shortener
go build
./main
```

Pass flags if needed
