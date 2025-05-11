// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package krn

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKrnParser(t *testing.T) {
	var krn *KRN

	var err error

	fullResourceName := "//service.example.com/key1/:value1:/key2/:value2:/key3/:value3:"
	krn, err = New(fullResourceName)
	assert.Nil(t, err)
	assert.Equal(t, "service.example.com", krn.ServiceName)
	assert.Equal(t, "service", ServiceID(krn.ServiceName, ".example.com"))
	assert.Equal(t, "key1/:value1:/key2/:value2:/key3/:value3:", krn.RelativeResourceName)
	assert.Equal(t, fullResourceName, krn.String())
}

func TestCollectionID(t *testing.T) {
	fullResourceName := "//service.example.com/key1/:value1:/key2/:value2:/key3/:value3:"
	krn, err := New(fullResourceName)
	assert.Nil(t, err)

	space, err := krn.ResourceID("key1")
	require.NoError(t, err)
	assert.Equal(t, ":value1:", space)

	region, err := krn.ResourceID("key2")
	require.NoError(t, err)
	assert.Equal(t, ":value2:", region)

	asset, err := krn.ResourceID("key3")
	require.NoError(t, err)
	assert.Equal(t, ":value3:", asset)
}

func TestEquals(t *testing.T) {
	fullResourceName := "//service.example.com/key1/:value1:/key2/:value2:/key3/:value3:"
	krn, err := New(fullResourceName)
	require.NoError(t, err)
	assert.True(t, krn.Equals(fullResourceName))

	krn, err = New("//service.example.com/key1/:value1:/key2/:value2:/key3/:value4:")
	require.NoError(t, err)
	assert.False(t, krn.Equals(fullResourceName))
}

func TestKrnToPropertyName(t *testing.T) {
	krn, err := New("//policy.api.kopexa.com/bundle/M%2FVSHZaChL8=/queries/sshdCiphers")
	require.NoError(t, err)

	res, err := krn.ResourceID("queries")
	require.NoError(t, err)
	assert.Equal(t, "sshdCiphers", res)

	res, err = krn.ResourceID("nothere")
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

// NewDataCategoryKRN generates a KRN for a data category using only spaceID and a UUID.
func newDataCategoryKRN(spaceID string) KRN {
	id := uuid.New().String()
	resourceID := fmt.Sprintf("%s-%s", spaceID, id)

	return KRN{
		ServiceName:          "kopexa.com",
		RelativeResourceName: fmt.Sprintf("data-categories/%s", resourceID),
	}
}

func TestDataCategoryKRN(t *testing.T) {
	spaceID := "space-123"
	krn := newDataCategoryKRN(spaceID)

	assert.Equal(t, "kopexa.com", krn.ServiceName)
	assert.Contains(t, krn.RelativeResourceName, "data-categories/")
	assert.Contains(t, krn.RelativeResourceName, spaceID)

	// Prüfe Format: "data-categories/space-123-..."
	parts := strings.Split(krn.RelativeResourceName, "/")
	assert.Equal(t, "data-categories", parts[0])
	assert.True(t, strings.HasPrefix(parts[1], spaceID+"-"))

	// JSON roundtrip
	jsonData, err := json.Marshal(krn)
	require.NoError(t, err)

	var parsed KRN
	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	assert.True(t, parsed.Equals(krn.String()))
}

func TestKRN_UnmarshalYAML(t *testing.T) {
	var target struct {
		Ref struct {
			KRN KRN `yaml:"krn"`
		} `yaml:"controlRefs"`
	}

	yamlData := []byte(`controlRefs:
  krn: "//kopexa.com/frameworks/iso-27001-2022"`)

	err := yaml.Unmarshal(yamlData, &target)
	require.NoError(t, err)
	assert.Equal(t, "kopexa.com", target.Ref.KRN.ServiceName)
	assert.Equal(t, "frameworks/iso-27001-2022", target.Ref.KRN.RelativeResourceName)
}

func TestNewChildKRN(t *testing.T) {
	type args struct {
		ownerKRN   string
		resource   string
		resourceID string
	}

	tests := []struct {
		name    string
		args    args
		want    *KRN
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				ownerKRN:   "//kopexa.com",
				resource:   "checks",
				resourceID: "123e4567-e89b-12d3-a456-426614174000",
			},
			want: &KRN{
				ServiceName:          "kopexa.com",
				RelativeResourceName: "checks/123e4567-e89b-12d3-a456-426614174000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChildKRN(tt.args.ownerKRN, tt.args.resource, tt.args.resourceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChildKRN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChildKRN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	assert.True(t, IsValid("//kopexa.com/frameworks/iso-27001-2022"))
	assert.False(t, IsValid("http://kopexa.com/frameworks/iso-27001-2022"))
	assert.False(t, IsValid(":invalid:krn"))
}

func TestMustNewAndMustParse(t *testing.T) {
	// MustNew should not panic for valid input
	assert.NotPanics(t, func() {
		krn := MustNew("//kopexa.com/frameworks/iso-27001-2022")
		assert.Equal(t, "kopexa.com", krn.ServiceName)
	})
	// MustNew should panic for invalid input
	assert.Panics(t, func() {
		_ = MustNew(":invalid:krn")
	})

	// MustParse should not panic for valid input
	assert.NotPanics(t, func() {
		krn := MustParse("//kopexa.com/frameworks/iso-27001-2022")
		assert.Equal(t, "kopexa.com", krn.ServiceName)
	})
	// MustParse should panic for invalid input
	assert.Panics(t, func() {
		_ = MustParse(":invalid:krn")
	})
}

func TestParseLegacy(t *testing.T) {
	krn, err := ParseLegacy("kopexa.com/frameworks/iso-27001-2022")
	assert.NoError(t, err)
	assert.Equal(t, "kopexa.com", krn.ServiceName)
	assert.Equal(t, "frameworks/iso-27001-2022", krn.RelativeResourceName)

	krn2, err := ParseLegacy("//kopexa.com/frameworks/iso-27001-2022")
	assert.NoError(t, err)
	assert.Equal(t, krn, krn2)

	_, err = ParseLegacy(":invalid")
	assert.Error(t, err)
}

func TestIsZero(t *testing.T) {
	assert.True(t, KRN{}.IsZero())
	assert.True(t, KRN{ServiceName: ""}.IsZero())
	assert.True(t, KRN{RelativeResourceName: ""}.IsZero())
	assert.False(t, KRN{ServiceName: "foo", RelativeResourceName: "bar"}.IsZero())
}

func TestBasenameAndCollectionName(t *testing.T) {
	krn := MustNew("//kopexa.com/frameworks/iso-27001-2022")
	assert.Equal(t, "iso-27001-2022", krn.Basename())
	assert.Equal(t, "frameworks", krn.CollectionName())

	krn2 := MustNew("//kopexa.com/")
	assert.Equal(t, "", krn2.Basename())
	assert.Equal(t, "", krn2.CollectionName())
}

func TestScanAndValue(t *testing.T) {
	krn := MustNew("//kopexa.com/frameworks/iso-27001-2022")
	var scanKRN KRN
	// Scan from string
	err := scanKRN.Scan("//kopexa.com/frameworks/iso-27001-2022")
	assert.NoError(t, err)
	assert.Equal(t, krn, &scanKRN)
	// Scan from []byte
	err = scanKRN.Scan([]byte("//kopexa.com/frameworks/iso-27001-2022"))
	assert.NoError(t, err)
	assert.Equal(t, krn, &scanKRN)
	// Scan from legacy string
	err = scanKRN.Scan("kopexa.com/frameworks/iso-27001-2022")
	assert.NoError(t, err)
	assert.Equal(t, krn, &scanKRN)
	// Scan from unsupported type
	err = scanKRN.Scan(123)
	assert.Error(t, err)
	// Value
	val, err := krn.Value()
	assert.NoError(t, err)
	assert.Equal(t, "//kopexa.com/frameworks/iso-27001-2022", val)
}
