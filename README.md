# Flag or Env

A simple go module to manage flags and env variables, and get them as needed.

**Supports Generics**

## Features

 - Get flags and environment variables
 - Specify default values
 - Specify a custom struct to fill
 - Specify preference, whichever value will be used first

# Basic Usage

See the example folder and test files for usage examples.

## Installation

    go get -u github.com/maxkruse/flagorenv

## Behaviour
If you pass both the Environment-Variable and the Flag, the config's `PreferFlag` determines which value will be used.

If you pass only one of the two, the value will be used.

If you omit the value, the default value will be used. Either using [go-default's](https://github.com/mcuadros/go-defaults) tags, or Zero-initialized.

## Dependencies

 - [go-strcase](github.com/stoewer/go-strcase) For converting strings interally to either SNAKE_CASE (env vars) or kebab-case (flags)