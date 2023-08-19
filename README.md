# Postman to http CLI tool

This is a tool for quickly converting postman requests to .http format.
It converts basic features but not everything is handled yet.

## How to use

Convert a collection to .http
`postmanToHtpp -output './collection' convert-collection ./path-to-json-collection`

Convert an environment to .env file
`postmanToHtpp -output './env' convert-environment ./path-to-json-environment`
