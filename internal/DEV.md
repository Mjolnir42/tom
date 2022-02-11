# Add new REST command

1. register command
    1. action keywords in `.../pkg/proto/constants.go`
    2. command keywords in `.../pkg/proto/${model}::${entity}.go
    3. command method definitions in `.../pkg/proto/${model}::${${entity}.go`
2. ensure model file `.../internal/model/${model}/${model}::.go`
3. ensure entity file `.../internal/model/${model}/${model}::${entity}:.go`
4. ensure command file `.../internal/model/${model}/${model}::${entity}:${command}.go
    1. init function calling AssertCommandIsDefined($command)
    2. init function calling registry append
    3. ensure registry callback function
    4. ensure data export function
    5. ensure rest facade method
    6. ensure internal handler method

# Add new cli command

1. register command in `.../internal/cli/model/${model}/${model}::${entity}:.go`
2. create entry for command in `.../internal/cli/adm/arguments.go ArgumentsForCommand()`
3. create bashcompletion in `.../internal/cli/cmpl/`
4. create documentation in `.../docs/tom/cmd_ref/${model}::${entity}:${command}`
5. create implementation in `.../internal/cli/model/${model}/${model}::${entity}:${command}.go`
