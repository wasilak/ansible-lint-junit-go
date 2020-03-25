package main

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
)

func TestAppVersion(t *testing.T) {
	want := AppVersion
	if got := appVersion(); got != want {
		t.Errorf("appVersion() = %q, want %q", got, want)
	}
}

func TestDeleteEmpty(t *testing.T) {
	want := []string{"aaa"}
	given := []string{"", "aaa", ""}
	got := deleteEmpty(given)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("deleteEmpty() = %q, want %q", got, want)
	}
}

func TestCreateXML(t *testing.T) {
	got := createXML("../test.yml:21: [E201] Trailing whitespace")
	want := `<testsuites>
    <testsuite errors="1" failures="0" tests="1" time="0">
        <testcase name="../test.yml">
            <failure file="../test.yml" line="21" message="[E201] Trailing whitespace" type="Ansible Lint">[E201] Trailing whitespace</failure>
        </testcase>
    </testsuite>
</testsuites>`

	xmlContent, err := xml.MarshalIndent(got, "", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	if string(xmlContent) != want {
		t.Errorf("createXML() = %q, want %q", string(xmlContent), want)
	}
}
