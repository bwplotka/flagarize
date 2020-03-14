package flagarize

import (
	"net/url"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/thanos-io/thanos/pkg/extflag"
	thanosmodel "github.com/thanos-io/thanos/pkg/model"
)

type testConfig struct {
	boolField bool `flagarize:"bool_field,<help>"`
	// How to make help better?
	boolFieldHidden  bool `flagarize:"bool_field_hidden,<help fd s sdf lol>,hidden"`
	stringField      string
	intField         int
	int64Field       int64
	uint64Field      uint64
	stringSliceField []string
	stringKVMapField map[string]string
	durationField    time.Duration
	timeField        time.Time

	customDurationField       model.Duration
	customDurationOrTimeField thanosmodel.TimeOrDurationValue
	customComplexTypeField    *extflag.PathOrContent
	urlField                  *url.URL

	wrongURLField url.URL

	// TODO rest
}

type testEmbeddedConfig testConfig

type testParseConfig struct {
	testEmbeddedConfig

	config1 testConfig
	config2 testConfig
	configs []testConfig
}

func TestParse(t *testing.T) {

}
