package florm

import (
	"testing"
	"time"
)

func TestGetTagKeysPerm(t *testing.T) {
	initialAPIOnce.Do(initializeAPI)

	shp := SchemaHelper{ss: NewFluxSession()}
	kps, err := shp.GetTagKeysPerm("misc2", "student", time.Unix(0, 0), time.Unix(0, 1))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(kps)
}

func TestGetTagKeysDistinct(t *testing.T) {
	initialAPIOnce.Do(initializeAPI)

	shp := SchemaHelper{ss: NewFluxSession()}
	kps, err := shp.GetTagKeysDistinct("misc2", "student", time.Unix(0, 0), time.Unix(0, 1))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(kps)
}

func TestGetFieldsPerm(t *testing.T) {
	initialAPIOnce.Do(initializeAPI)

	shp := SchemaHelper{ss: NewFluxSession()}
	kps, err := shp.GetFieldsPerm("misc2", "student", time.Unix(0, 0), time.Unix(0, 1))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(kps)
}

func TestGetFieldsDistinct(t *testing.T) {
	initialAPIOnce.Do(initializeAPI)

	shp := SchemaHelper{ss: NewFluxSession()}
	kps, err := shp.GetFieldsDistinct("misc2", "student", time.Unix(0, 0), time.Unix(0, 1))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(kps)
}
