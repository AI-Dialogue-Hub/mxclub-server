package event

import (
	"encoding/json"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"mxclub/pkg/utils"
	"testing"
)

func TestEvent(t *testing.T) {
	hashMap := maps.NewLinkedHashMap[int, []*EventBO]()
	eventBOS := make([]*EventBO, 0)
	eventBOS = append(eventBOS, &EventBO{EventCode: 111})
	hashMap.Put(0, eventBOS)
	bos := hashMap.MustGet(0)
	bos = append(bos, &EventBO{EventCode: 222})
	hashMap.Put(0, bos)
	t.Logf("%+v", hashMap)
	t.Logf("%+v", bos)
	data, _ := json.Marshal(bos)
	t.Logf("%+v", utils.ObjToJsonStr(bos))
	t.Logf("%+v", string(data))
}
