# Chameleon

A CLI tool that helps with copying DynamoDB data.

## Installation

### Homebrew

```bash
brew tap youjinp/brew
brew install chameleon
```

## Usage

### Download dynamodb table to a file

```bash
export AWS_ACCESS_KEY_ID='id'
export AWS_SECRET_ACCESS_KEY='key'
export AWS_SESSION_TOKEN='token'
export AWS_DEFAULT_REGION='us-east-1'

chameleon copy -t <tablename> -o <outputfile>
```

### Write data into dynamodb from the created file

```bash
export AWS_ACCESS_KEY_ID='id'
export AWS_SECRET_ACCESS_KEY='key'
export AWS_SESSION_TOKEN='token'
export AWS_DEFAULT_REGION='us-east-1'

chameleon paste -t <tablename> -f <inputfile>
```
