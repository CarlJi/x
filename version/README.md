# version
package version defines the utility information for versioning binary

## Usage

Import this package in your application:

```go
import _ "github.com/qiniu/x/version"
```

Then build it with:

```shell
go build -ldflags "-X 'github.com/qiniu/x/version.BuildDate=$(date)'" .
```

