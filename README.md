# device-management
device-management
## Project directory
```
├─api/                              # Openapi specs for API interface
├─config/                           # Application configurations and properties
├─controller/                       # Application API endpoints handler
├─router/          
├─model/                            # Data model stucts
├─service/                          # Service where buisness logic written
├─repository/                       # Repositories to handle database related work
├─exciption/                        # Error and error handling logic
├─db/                               # Database schema and data related code
│   ├── migrations/                 # Sql migration files
│   └─ setup/                       # CI/Local/Test prepartion files
└─.github/workflows/build.yaml      # Github workflow CI/CD pipline
```

## Summary
1. API sample application for devices management

## Migration
- application uses [golang migrate](https://github.com/golang-migrate/migrate) to control database code schema as code.
- each time new changes to database schema need to be placed new file you can use [golang migrate cmd](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- the migrate cmd will create tow incremental versioned files with the following format `<version>_<patch_name>.<up/down>.sql`, with up suffix to apply 
  patch, down to rollback patch.
- each time we need to apply recently created patches we can execute `make db+`
## To Run locally 
In order to make running application locally easier execute the following script:
- [run-ci](Makefile)
## Application config override
You can override application config at run time by the following ways:
1. specify it as env variable, it will have higher precedence.
2. add config-<profile>.json config override and when run the app specify env:PROFILE for the profile to load it's config override

## CI/CD with github action runner
see [github action](https://docs.github.com/en/actions)

# TODO Next
Add  build docker image and deploy script to cloud