# Postman to http CLI tool

This is a tool for quickly converting postman requests to .http format.
It converts basic features but not everything is handled yet.

## How to use

Convert a collection to .http
`postmanToHtpp -input 'somePostmanCollection.json' -output './output' generate-collection`

Convert an environment to .env file
`postmanToHtpp -input 'somePostmanEnvironment.json' -output './env' generate-env-file`
