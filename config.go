package main

type Action struct {
    Type string;
    Params interface{};
}

type Rule struct {
    expression string;
    actions []Action
}

type Config struct {
    address string;
    workspacePath string;
    rules []Rule;
}

func LoadConfig(path string) Config {
    return Config{
        address: ":8080",
        workspacePath: "workspace",
    }
}


