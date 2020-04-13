# SampleApp for AWS X-Ray Go SDK

This repository contains sample app to show the tracing use case of aws-xray-sdk-go. The SampleApp contains example of tracing aws sdk calls like list all SQS queues and list all s3 buckets. Moreover, it contains tracing SQL request (creating, deleting table and populating data inside that table) and tracing upstream HTTP request. 

## Prerequirements

* Should have a mysql database setup since SampleApp will be querying to the local database 
* Should have XRay daemon installed and running in order to see traces on the AWS XRay console

The following environment variable is expected to set by the customer
```
DSN_STRING - The connection string of the database (set username, password and dbname)
```
NOTE: example of recommended dsn string: username:password@tcp(127.0.0.1:3306)/dbname

## Setup

The SampleApp for AWS X-Ray Go SDK is compatible with Go 1.9 and above. Go modules (go.mod) will fetch all the dependency this sample app requires.

To clone the SampleApp in your environment,
```
git clone https://github.com/aws-samples/aws-xray-sdk-go-sample.git
```
To run the SampleApp with environment variable set up,
```
DSN_STRING="username:password@tcp(127.0.0.1:3306)/dbname" go run src/main.go
```
## Opening Issues

If you encounter a bug specifically with the SampleApp for AWS X-Ray Go SDK should be reported to this repository whereas bugs with the SDK should be reported [here](https://github.com/aws/aws-xray-sdk-go/issues). Search the [existing issues](https://github.com/aws/aws-xray-sdk-go/issues) and see if others are also experiencing the issue before opening a new issue. The GitHub issues are intended for bug reports and feature requests.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.