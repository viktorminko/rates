# Moving average for currency rates

Collect moving averages for specific interval

## Usage

By default app returns data on 0.0.0.0:3512

### Binary

```
./rates:
  -api_endpoint string
        API endpoint to request average rates (default "/")
  -get_pairs_url string
        URL to get pairs information (default "https://wex.nz/api/3/info/")
  -pairs string
        comma separated currency pairs : btc_usd, eth_eur, xrp_btc
  -port int
        port for connections (default 3512)
  -timeframe duration
        time interval to use for moving averages (default 10m0s)
  -update_duration duration
        how often rates are updated (default 5s)
  -update_url string
        URL to update rates for pair (default "https://wex.nz/api/3/ticker/")
```

### Docker image

Run docker image with default config

```
docker run viktorminko/rates
```

### Docker compose

Customize configuration in .env file and run

```
docker-compose up
```

