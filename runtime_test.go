package runtime_test

import (
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	r "runtime"
	"strings"
	"testing"
	"time"

	"trpc.group/trpc-go/trpc-go/errs"
	runtime "trpc.group/trpc-go/trpc-metrics-runtime"

	"github.com/agiledragon/gomonkey"
)

func TestRuntime(t *testing.T) {

	patch2 := gomonkey.ApplyFunc(http.Post, func(url, contentType string,
		body io.Reader) (resp *http.Response, err error) {
		return &http.Response{Body: &net.TCPConn{}}, errs.New(1, "test")
	})
	defer patch2.Reset()

	patch3 := gomonkey.ApplyFunc(os.Open, func(name string) (*os.File, error) {
		return os.NewFile(uintptr(3), "test"), nil
	})
	defer patch3.Reset()

	patch5 := gomonkey.ApplyFunc(strings.Split, func(s, sep string) []string {
		return []string{"1", "2", "3"}
	})
	defer patch5.Reset()

	mockRet := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{
				1, nil,
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, nil,
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		}, {
			Values: gomonkey.Params{
				1, errs.New(1, "test"),
			},
		},
	}
	a := &os.File{}
	patch4 := gomonkey.ApplyMethodSeq(reflect.TypeOf(a), "Read", mockRet)
	defer patch4.Reset()

	patch6 := gomonkey.ApplyFunc(r.ReadMemStats, func(m *r.MemStats) {
		*m = r.MemStats{PauseNs: [256]uint64{1, 100e3, 500e3, 1e6, 10e6, 50e6, 100e6, 500e6, 1e9}}
	})
	defer patch6.Reset()

	runtime.RuntimeMetrics()
	time.Sleep(time.Second * 6)
}

func TestRuntime2(t *testing.T) {
	runtime.RuntimeMetrics()
	time.Sleep(time.Second * 6)
}
