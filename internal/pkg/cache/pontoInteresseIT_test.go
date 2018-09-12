package cache

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestMapKeyReferenc(t *testing.T) {

	m := make(map[bson.ObjectId]string)

	id1 := bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")
	id2 := bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")

	k1 := id1
	k2 := id2

	ori := "t1"
	exp := "t2"

	m[k1] = ori
	m[k2] = exp

	if len(m) != 1 {
		t.Errorf("Mapa deveria ter 1, mas possui %d elementos", len(m))
	}

	if m[k1] != exp {
		t.Errorf("Valor deveria ser %q, mas é %q", exp, m[k1])
	}

}
