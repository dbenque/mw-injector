rm $GOPATH/src/github.com/dbenque/mw-injector/pkg/api/mwinjector/v1/deepcopy_generated.go
rm -rf $GOPATH/src/github.com/dbenque/mw-injector/pkg/client/informers
rm -rf $GOPATH/src/github.com/dbenque/mw-injector/pkg/client/listers
rm -rf $GOPATH/src/github.com/dbenque/mw-injector/pkg/client/versioned


$GOPATH/bin/deepcopy-gen -i "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1" --output-file-base deepcopy_generated --go-header-file="$GOPATH/src/github.com/dbenque/mw-injector/hack/boilerplate.go.txt" --logtostderr
$GOPATH/bin/client-gen --input-base "github.com/dbenque/mw-injector/pkg/api" --input "mwinjector/v1" --go-header-file="$GOPATH/src/github.com/dbenque/mw-injector/hack/boilerplate.go.txt" --clientset-path "github.com/dbenque/mw-injector/pkg/client" --clientset-name "versioned"
$GOPATH/bin/lister-gen -i "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1" --output-package "github.com/dbenque/mw-injector/pkg/client/listers/" --go-header-file="$GOPATH/src/github.com/dbenque/mw-injector/hack/boilerplate.go.txt" --logtostderr
$GOPATH/bin/informer-gen -i "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1" --output-package "github.com/dbenque/mw-injector/pkg/client/informers/" --go-header-file="$GOPATH/src/github.com/dbenque/mw-injector/hack/boilerplate.go.txt" --versioned-clientset-package "github.com/dbenque/mw-injector/pkg/client/versioned" --listers-package "github.com/dbenque/mw-injector/pkg/client/listers"
