//go:build !binary_log
// +build !binary_log

package zerolog_test

import (
	"bytes"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/treavorj/zerolog"
)

func ExampleNew() {
	log := zerolog.New(os.Stdout)

	log.Info().Msg("hello world")
	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_With() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","foo":"bar","message":"hello world"}
}

func ExampleLogger_Level() {
	log := zerolog.New(os.Stdout).Level(zerolog.WarnLevel)

	log.Info().Msg("filtered out message")
	log.Error().Msg("kept message")

	// Output: {"level":"error","message":"kept message"}
}

func ExampleLogger_Sample() {
	log := zerolog.New(os.Stdout).Sample(&zerolog.BasicSampler{N: 2})

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")

	// Output: {"level":"info","message":"message 1"}
	// {"level":"info","message":"message 3"}
}

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	if l != zerolog.NoLevel {
		e.Str("level_name", l.String())
	} else {
		e.Str("level_name", "NoLevel")
	}
}

type MessageHook string

func (h MessageHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Str("the_message", msg)
}

func ExampleLogger_Hook() {
	var levelNameHook LevelNameHook
	var messageHook MessageHook = "The message"

	log := zerolog.New(os.Stdout).Hook(levelNameHook, messageHook)

	log.Info().Msg("hello world")

	// Output: {"level":"info","level_name":"info","the_message":"hello world","message":"hello world"}
}

func ExampleLogger_Print() {
	log := zerolog.New(os.Stdout)

	log.Print("hello world")

	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Printf() {
	log := zerolog.New(os.Stdout)

	log.Printf("hello %s", "world")

	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Println() {
	log := zerolog.New(os.Stdout)

	log.Println("hello world")

	// Output: {"level":"debug","message":"hello world\n"}
}

func ExampleLogger_Trace() {
	log := zerolog.New(os.Stdout)

	log.Trace().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: {"level":"trace","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Debug() {
	log := zerolog.New(os.Stdout)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: {"level":"debug","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Info() {
	log := zerolog.New(os.Stdout)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: {"level":"info","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Warn() {
	log := zerolog.New(os.Stdout)

	log.Warn().
		Str("foo", "bar").
		Msg("a warning message")

	// Output: {"level":"warn","foo":"bar","message":"a warning message"}
}

func ExampleLogger_Error() {
	log := zerolog.New(os.Stdout)

	log.Error().
		Err(errors.New("some error")).
		Msg("error doing something")

	// Output: {"level":"error","error":"some error","message":"error doing something"}
}

func ExampleLogger_WithLevel() {
	log := zerolog.New(os.Stdout)

	log.WithLevel(zerolog.InfoLevel).
		Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_Write() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")

	// Output: {"foo":"bar","message":"hello world"}
}

func ExampleLogger_Log() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")

	// Output: {"foo":"bar","bar":"baz"}
}

func ExampleEvent_Dict() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).
		Msg("hello world")

	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

type User struct {
	Name    string
	Age     int
	Created time.Time
}

func (u User) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", u.Name).
		Int("age", u.Age).
		Time("created", u.Created)
}

type Price struct {
	val  uint64
	prec int
	unit string
}

func (p Price) MarshalZerologObject(e *zerolog.Event) {
	denom := uint64(1)
	for i := 0; i < p.prec; i++ {
		denom *= 10
	}
	result := []byte(p.unit)
	result = append(result, fmt.Sprintf("%d.%d", p.val/denom, p.val%denom)...)
	e.Str("price", string(result))
}

type Users []User

func (uu Users) MarshalZerologArray(a *zerolog.Array) {
	for _, u := range uu {
		a.Object(u)
	}
}

func ExampleEvent_Array() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1).
			Dict(zerolog.Dict().
				Str("bar", "baz").
				Int("n", 1),
			),
		).
		Msg("hello world")

	// Output: {"foo":"bar","array":["baz",1,{"bar":"baz","n":1}],"message":"hello world"}
}

func ExampleEvent_Array_object() {
	log := zerolog.New(os.Stdout)

	// Users implements zerolog.LogArrayMarshaler
	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	log.Log().
		Str("foo", "bar").
		Array("users", u).
		Msg("hello world")

	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
}

func ExampleEvent_Object() {
	log := zerolog.New(os.Stdout)

	// User implements zerolog.LogObjectMarshaler
	u := User{"John", 35, time.Time{}}

	log.Log().
		Str("foo", "bar").
		Object("user", u).
		Msg("hello world")

	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
}

func ExampleEvent_EmbedObject() {
	log := zerolog.New(os.Stdout)

	price := Price{val: 6449, prec: 2, unit: "$"}

	log.Log().
		Str("foo", "bar").
		EmbedObject(price).
		Msg("hello world")

	// Output: {"foo":"bar","price":"$64.49","message":"hello world"}
}

func ExampleEvent_Interface() {
	log := zerolog.New(os.Stdout)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().
		Str("foo", "bar").
		Interface("obj", obj).
		Msg("hello world")

	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
}

func ExampleEvent_Dur() {
	d := 10 * time.Second

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dur("dur", d).
		Msg("hello world")

	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
}

func ExampleEvent_Durs() {
	d := []time.Duration{
		10 * time.Second,
		20 * time.Second,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Durs("durs", d).
		Msg("hello world")

	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}

func ExampleEvent_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleEvent_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleContext_Dict() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

func ExampleContext_Array() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1),
		).Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","array":["baz",1],"message":"hello world"}
}

func ExampleContext_Array_object() {
	// Users implements zerolog.LogArrayMarshaler
	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Array("users", u).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
}

func ExampleContext_Object() {
	// User implements zerolog.LogObjectMarshaler
	u := User{"John", 35, time.Time{}}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Object("user", u).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
}

func ExampleContext_EmbedObject() {

	price := Price{val: 6449, prec: 2, unit: "$"}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		EmbedObject(price).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","price":"$64.49","message":"hello world"}
}

func ExampleContext_Interface() {
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Interface("obj", obj).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
}

func ExampleContext_Dur() {
	d := 10 * time.Second

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dur("dur", d).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
}

func ExampleContext_Durs() {
	d := []time.Duration{
		10 * time.Second,
		20 * time.Second,
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Durs("durs", d).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}

func ExampleContext_IPAddr() {
	hostIP := net.IP{192, 168, 0, 100}
	log := zerolog.New(os.Stdout).With().
		IPAddr("HostIP", hostIP).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"HostIP":"192.168.0.100","message":"hello world"}
}

func ExampleContext_IPPrefix() {
	route := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
	log := zerolog.New(os.Stdout).With().
		IPPrefix("Route", route).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"Route":"192.168.0.0/24","message":"hello world"}
}

func ExampleContext_MACAddr() {
	mac := net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
	log := zerolog.New(os.Stdout).With().
		MACAddr("hostMAC", mac).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"hostMAC":"00:14:22:01:23:45","message":"hello world"}
}

func ExampleContext_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Fields(fields).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleContext_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Fields(fields).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleContext_DeDup() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Str("foo", "baz").
		DeDup().
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","foo":"baz","message":"hello world"}
}

func ExampleContext_DeDup_unused() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Str("foo", "baz").
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","foo":"bar","foo":"baz","message":"hello world"}
}

func ExampleContext_DeDup_empty() {
	log := zerolog.New(os.Stdout).
		With().
		DeDup().
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleContext_DeDup_dictionary() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Dict("dict", zerolog.Dict().
			Str("foo", "bar").
			Int("n", 1),
		).
		Dict("dict", zerolog.Dict().
			Str("foo", "baz").
			Int("n", 2),
		).
		DeDup().
		Logger()

	log.Info().Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"level":"info","dict":{"foo":"baz","n":2},"message":"hello world"}`
	expectedFormat2 := `{"dict":{"foo":"bam","n":3},"level":"info","message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}

func ExampleContext_DeDup_array() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Array("array", zerolog.Arr().
			Str("bar").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "bar").
				Int("n", 1),
			),
		).
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "baz").
				Int("n", 2),
			),
		).
		DeDupDeep().
		Logger()

	log.Info().Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"array":["baz",1,{"foo":"baz","n":2}],"level":"info","message":"hello world"}`
	expectedFormat2 := `{"level":"info","array":["baz",1,{"foo":"baz","n":2}],"message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}

func ExampleEvent_DeDup() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Str("foo", "baz").
		Logger()

	log.Info().Str("foo", "bam").DeDup().Msg("hello world")

	// Output: {"level":"info","foo":"bam","message":"hello world"}
}

func ExampleEvent_DeDup_unused() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Str("foo", "baz").Msg("hello world")

	// Output: {"level":"info","foo":"bar","foo":"baz","message":"hello world"}
}

func ExampleEvent_DeDup_empty() {
	log := zerolog.New(os.Stdout).
		With().
		Logger()

	log.Info().DeDup().Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleEvent_DeDup_dictionary() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Dict("dict", zerolog.Dict().
			Str("foo", "bar").
			Int("n", 1),
		).
		Dict("dict", zerolog.Dict().
			Str("foo", "baz").
			Int("n", 2),
		).
		Logger()

	log.Info().
		Dict("dict", zerolog.Dict().
			Str("foo", "bam").
			Int("n", 3),
		).
		DeDup().
		Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"level":"info","dict":{"foo":"bam","n":3},"message":"hello world"}`
	expectedFormat2 := `{"dict":{"foo":"bam","n":3},"level":"info","message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}

func ExampleEvent_DeDup_array() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Array("array", zerolog.Arr().
			Str("bar").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "bar").
				Int("n", 1),
			),
		).
		DeDupDeep().
		Logger()

	log.Info().
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "baz").
				Int("n", 2),
			),
		).
		DeDupDeep().
		Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"array":["baz",1,{"foo":"baz","n":2}],"level":"info","message":"hello world"}`
	expectedFormat2 := `{"level":"info","array":["baz",1,{"foo":"baz","n":2}],"message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}

func ExampleContext_DeDupDeep() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Str("foo", "baz").
		DeDupDeep().
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","foo":"baz","message":"hello world"}
}

func ExampleContext_DeDupDeep_unused() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Str("foo", "baz").
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","foo":"bar","foo":"baz","message":"hello world"}
}

func ExampleContext_DeDupDeep_empty() {
	log := zerolog.New(os.Stdout).
		With().
		DeDupDeep().
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleContext_DeDupDeep_dictionary() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Dict("dict", zerolog.Dict().
			Str("foo", "bar").
			Int("n", 1),
		).
		Dict("dict", zerolog.Dict().
			Str("foo", "baz").
			Int("n", 2),
		).
		DeDupDeep().
		Logger()

	log.Info().Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"level":"info","dict":{"foo":"baz","n":2},"message":"hello world"}`
	expectedFormat2 := `{"dict":{"foo":"bam","n":3},"level":"info","message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}

func ExampleContext_DeDupDeep_array() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Array("array", zerolog.Arr().
			Str("bar").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "bar").
				Int("n", 1),
			),
		).
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "baz").
				Int("n", 2),
			),
		).
		DeDupDeep().
		Logger()

	log.Info().Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"array":["baz",1,{"foo":"baz","n":2}],"level":"info","message":"hello world"}`
	expectedFormat2 := `{"level":"info","array":["baz",1,{"foo":"baz","n":2}],"message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}

func ExampleEvent_DeDupDeep() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Str("foo", "baz").
		Logger()

	log.Info().Str("foo", "bam").DeDupDeep().Msg("hello world")

	// Output: {"foo":"bam","level":"info","message":"hello world"}
}

func ExampleEvent_DeDupDeep_unused() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Str("foo", "baz").Msg("hello world")

	// Output: {"level":"info","foo":"bar","foo":"baz","message":"hello world"}
}

func ExampleEvent_DeDupDeep_empty() {
	log := zerolog.New(os.Stdout).
		With().
		Logger()

	log.Info().DeDupDeep().Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleEvent_DeDupDeep_dictionary() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Dict("dict", zerolog.Dict().
			Str("foo", "bar").
			Int("n", 1),
		).
		Dict("dict", zerolog.Dict().
			Str("foo", "baz").
			Int("n", 2),
		).
		Logger()

	log.Info().
		Dict("dict", zerolog.Dict().
			Str("foo", "bam").
			Int("n", 3),
		).
		DeDupDeep().
		Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"level":"info","dict":{"foo":"bam","n":3},"message":"hello world"}`
	expectedFormat2 := `{"dict":{"foo":"bam","n":3},"level":"info","message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}
func ExampleEvent_DeDupDeep_array() {
	var buf bytes.Buffer
	log := zerolog.New(&buf).
		With().
		Array("array", zerolog.Arr().
			Str("bar").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "bar").
				Int("n", 1),
			),
		).
		DeDupDeep().
		Logger()

	log.Info().
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1).
			Dict(zerolog.Dict().
				Str("foo", "baz").
				Int("n", 2),
			),
		).
		DeDupDeep().
		Msg("hello world")

	output := strings.TrimSpace(buf.String())
	expectedFormat1 := `{"array":["baz",1,{"foo":"baz","n":2}],"level":"info","message":"hello world"}`
	expectedFormat2 := `{"level":"info","array":["baz",1,{"foo":"baz","n":2}],"message":"hello world"}`

	if output != expectedFormat1 && output != expectedFormat2 {
		fmt.Printf("Output did not match either expected format.\n  Got:      %s\n  Expected: %s\n  or:       %s",
			output, expectedFormat1, expectedFormat2)
	}

	// Output:
}
