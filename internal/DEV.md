# Add new REST command

1. register command in `.../pkr/proto/shared__def.go Commands
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

1. register command in `.../internal/cli/model/meta/meta::namespace:.go`
2. create entry for command in adm.ArgumentsForCommand()
3. create bashcompletion in `.../internal/cli/cmpl/`
