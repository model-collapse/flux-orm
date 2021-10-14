package florm

import (
	"testing"
	"time"
)

func TestCheckSchemaStrict(t *testing.T) {
	fillMockDataM3Once.Do(fillMockDataM3)

	if err := CheckSchema(&CarTypeS{}, defaultAPIManager, CheckModeStrict); err != nil {
		t.Fatalf("expect no error from CarTypeS, but %v", err)
	}

	if err := CheckSchema(&CarTypeSS{}, defaultAPIManager, CheckModeStrict); err == nil {
		t.Fatalf("expect an error from CarTypeSS")
	} else {
		t.Log("CarTypeSS hit err = ", err)
	}

	if err := checkSchema(&CarTypeD{}, defaultAPIManager, CheckModeStrict, time.Unix(0, 0), time.Unix(10, 0)); err == nil {
		t.Fatalf("expect an error from CarTypeD")
	} else {
		if err == ErrNoSchemaData {
			t.Fatal(err)
		}
		t.Log("CarTypeD hit err = ", err)
	}
}

func TestCheckSchemaCompatible(t *testing.T) {
	fillMockDataM3Once.Do(fillMockDataM3)

	if err := checkSchema(&CarTypeD{}, defaultAPIManager, CheckModeCompatible, time.Unix(0, 0), time.Unix(5, 0)); err != nil {
		t.Fatalf("expect no error from CarTypeD, but %v", err)
	}

	if err := checkSchema(&CarTypeD{}, defaultAPIManager, CheckModeCompatible, time.Unix(0, 0), time.Unix(6, 0)); err == nil {
		t.Fatalf("expect an error from CarTypeD")
	} else {
		t.Log("CarTypeD hit err = ", err)
	}
}
