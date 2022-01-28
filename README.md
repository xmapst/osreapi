# OSREApi -- OS Remote Executor Api

## Help
```text
usage: osreapi [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
      --addr=":2376"      host:port for execution.
      --log_level="info"  Set log output level.
      --exec_timeout=30m  Set the default exec command expire time.
      --max-requests=0    Maximum number of concurrent requests. 0 to disable.
      --data="/tmp/osreapi/data"  
                          Save data with memory or custom path.
      --version           Show application version.

```

## Swagger
![osreapi](https://raw.githubusercontent.com/xmapst/osreapi/main/swagger.jpg)

### Linux
```shell
/your/path/osreapi --addr=:2376
```

## Windows
Open powershell execution in administrative mode

```powershell
New-Service -Name osreapi -BinaryPathName "C:\your\path\osreapi.exe --addr=:2376" -DisplayName  "OS Remote Executor " -StartupType Automatic
sc.exe failure osreapi reset= 0 actions= restart/0/restart/0/restart/0
sc.exe start osreapi
```

## Compile locally
### Windows
```powershell
git clone https://github.com/xmapst/osreapi.git
cd osreapi
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -ldflags="-linkmode internal" -o osreapi.exe cmd/windows/main.go
upx.exe --lzma osreapi.exe
```
### Linux
```shell
git clone https://github.com/xmapst/osreapi.git
cd osreapi
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -ldflags="-linkmode internal" -o osreapi cmd/linux/main.go
upx --lzma osreapi
```

## Test
### Get all current tasks
```shell
# Sort by start time by default
curl -XGET http://localhost:2376
# Sort by completion time
curl -XGET http://localhost:2376/?sort=completed
# Sort by expiration time
curl -XGET http://localhost:2376/?sort=expired
```
Result: 
```json
{
  "code":0,
  "message":"success",
  "data":{
    "running":0,
    "tasks":[
      {
        "id":"8e04b882-a55d-4d20-b70e-89201da6b214",
        "state":"Stopped",
        "times":{
          "start":"2022-01-04T20:35:14+08:00",
          "completed":"2022-01-04T20:35:21+08:00",
          "expired":"2022-01-04T21:25:14+08:00"
        }
      },
      {
        "id":"fd7b0a7a-ed05-4e37-b010-ba234300c525",
        "state":"Stopped",
        "times":{
          "start":"2022-01-04T20:35:14+08:00",
          "completed":"2022-01-04T20:35:21+08:00",
          "expired":"2022-01-04T21:25:14+08:00"
        }
      }
    ],
    "total": 2
  }
}
```
### Execute a script or command
#### Execute in order
Request:
```shell
curl -XPOST http://localhost:2376 -d '[
  {
    "name": "step0",
    "command_type": "cmd",
    "command_content": "ping baidu.com",
    "timeout": "5m",
    "envs": {
      "env1":"a",
      "env2":"b",
      "env3":"c"
    }
  },
  {
    "name": "step1",
    "command_type": "cmd",
    "command_content": "curl https://www.baidu.com",
    "timeout": "30m",
    "envs": {
      "env1":"a",
      "env2":"b",
      "env3":"c"
    },
    "depends_on": [
    "step0"
    ]
  }
  ,
  {
    "name": "step2",
    "command_type": "cmd",
    "command_content": "set",
    "envs": {
      "env1":"a",
      "env2":"b",
      "env3":"c"
    },
    "depends_on": [
    "step1"
    ]
  }
]'
```

#### Concurrent execution
Request:
```shell
curl -XPOST http://localhost:2376 -d '[
  {
    "name": "step0",
    "command_type": "cmd",
    "command_content": "ping baidu.com",
    "timeout": "5m",
    "envs": {
      "env1":"a",
      "env2":"b",
      "env3":"c"
    }
  },
  {
    "name": "step1",
    "command_type": "cmd",
    "command_content": "curl https://www.baidu.com",
    "timeout": "30m",
    "envs": {
      "env1":"a",
      "env2":"b",
      "env3":"c"
    }
  }
  ,
  {
    "name": "step2",
    "command_type": "cmd",
    "command_content": "set",
    "envs": {
      "env1":"a",
      "env2":"b",
      "env3":"c"
    }
  }
]'
```
Return content:  
Success:
```json
{
  "code": 0,
  "data":{
    "id": "87a25e00-5c35-453c-93c8-36109a24a104",
    "step_count": 3
  }
}
```
Parameter error:
```json
{
  "code": 1000,
  "message": "xxxxx, parameter error"
}
```

# Get results  
Request:
```shell
curl http://localhost:2376/87a25e00-5c35-453c-93c8-36109a24a104
```

Return content:  
success:
```json
{
  "code": 0,
  "data": [
    {
      "name": "xxxxxx",
      "step": 0,
      "exit_code": 0,
      "output": "xxxxxxxxxxxxxxxxx"
    } 
  ]
}
```
Parameter error:
```json
{
  "code": 1000,
  "message": "missing id parameter, parameter error"
}
```
Is executing:
```json
{
  "code": 1001,
  "message": "is executing"
}
```
:waiting to run:
```json
{
  "code": 1002,
  "message": "waiting to run"
}
```
Failed to execute:
```json
{
  "code": 1003,
  "message": "1, failed to execute",
  "data": [
    {
      "name": "xxxx",
      "step": 0,
      "exit_code": 100,
      "output": "xxxxxxxxxxxxxxxxx"
    }
  ]
}
```
No data:
```json
{
  "code": 1004,
  "message": "id does not exist, no data"
}
```

[Notes]
+ The return code has the following:
    - 0: success
    - 500: internal error
    - 1000: parameter error
    - 1001: is executing
    - 1002: waiting to run
    - 1003: failed to execute
    - 1004: no data
