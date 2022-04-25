# Flag or Env

A simple go module to manage flags and env variables, and get them as needed.

**Supports Generics**

## Features

 - Get flags and environment variables
 - Specify default values
 - Specify a custom struct to fill
 - Specify preference, whichever value will be used first

# Basic Usage

See the example folder and test files for usage examples

## Dependencies

 - [go-strcase](github.com/stoewer/go-strcase) For converting strings interally to either SNAKE_CASE (env vars) or kebab-case (flags)