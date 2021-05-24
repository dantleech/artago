package main

type Config struct {
    address string;
    workspacePath string;
}

func LoadConfig() Config {
    return Config{
        address: ":8080",
        workspacePath: "workspace",
    }
}


