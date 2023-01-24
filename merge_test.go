package go_cypherdsl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeSetConfig_ToString(t *testing.T) {
	t1 := MergeSetConfig{Name: "test", Member: "ttt", Target: 1}
	t2 := MergeSetConfig{Name: "test", Member: "ttt", TargetFunction: &FunctionConfig{
		Name: "test",
	}}

	t3 := MergeSetConfig{Name: "test"}
	t4 := MergeSetConfig{}
	t5 := MergeSetConfig{Name: "test", Member: "ttt"}
	t6 := MergeSetConfig{Name: "test", Member: "ttt", TargetFunction: &FunctionConfig{Name: "test"}, Target: 1}

	req := require.New(t)
	var err error
	var cypher string

	//name member normal target
	cypher, err = t1.ToString()
	req.Nil(err)
	req.EqualValues("test.ttt = 1", cypher)

	//name member target function
	cypher, err = t2.ToString()
	req.Nil(err)
	req.EqualValues("test.ttt = test()", cypher)

	//error - member not defined
	_, err = t3.ToString()
	req.NotNil(err)

	//error - member and name not defined
	_, err = t4.ToString()
	req.NotNil(err)

	//error - target and target function not defined
	_, err = t5.ToString()
	req.NotNil(err)

	//error - target and target function defined
	_, err = t6.ToString()
	req.NotNil(err)
}

func TestMergeSetConfigWithMembers_ToString(t *testing.T) {
	t1 := MultiMemberMergeSetConfig{Members: map[string]interface{}{"key": "value"}}
	t2 := MultiMemberMergeSetConfig{Name: "test"}
	t3 := MultiMemberMergeSetConfig{Name: "test", Members: map[string]interface{}{}}
	t4 := MultiMemberMergeSetConfig{Name: "test", Members: map[string]interface{}{"key1": 1, "key2": "value2", "key3": ParamString("$value3")}}

	req := require.New(t)
	var err error
	var cypher string

	//error - name not defined
	_, err = t1.ToString()
	req.NotNil(err)

	//error - members not defined
	_, err = t2.ToString()
	req.NotNil(err)

	//error - members empty
	_, err = t3.ToString()
	req.NotNil(err)

	//name members
	cypher, err = t4.ToString()
	req.Nil(err)
	req.Contains(cypher, "test.key1 = 1", "test.key2 = 'value2'", "test.key3 = $value3")
}

func TestMergeConfig_ToString(t *testing.T) {
	t1 := MergeConfig{Path: "test"}

	t2 := MergeConfig{Path: "test", OnCreate: &MergeSetConfig{
		Name:   "test",
		Member: "ttt",
		Target: 1,
	}}

	t3 := MergeConfig{Path: "test", OnMatch: &MergeSetConfig{
		Name:   "test",
		Member: "ttt",
		Target: 1,
	}}

	t4 := MergeConfig{Path: "test", OnMatch: &MergeSetConfig{
		Name:   "test",
		Member: "ttt",
		Target: 1,
	}, OnCreate: &MergeSetConfig{
		Name:   "test",
		Member: "tt1",
		Target: 2,
	}}

	t5 := MergeConfig{}

	t6 := MergeConfig{Path: "test", OnMatch: &MergeSetConfig{
		Name:   "test",
		Target: ParamString("$props"),
	}, OnCreate: &MergeSetConfig{
		Name:   "test",
		Target: ParamString("$props"),
	}}

	t7 := MergeConfig{Path: "test", OnMatchSetMembers: &MultiMemberMergeSetConfig{
		Name:    "test",
		Members: map[string]interface{}{"key1": 1, "key2": "value2", "key3": ParamString("$value3")},
	}, OnCreate: &MergeSetConfig{
		Name:   "test",
		Target: ParamString("$props"),
	}}

	req := require.New(t)
	var err error
	var cypher string

	//only merge
	cypher, err = t1.ToString()
	req.Nil(err)
	req.EqualValues("test", cypher)

	//merge with on create
	cypher, err = t2.ToString()
	req.Nil(err)
	req.EqualValues("test ON CREATE SET test.ttt = 1", cypher)

	//merge with on match
	cypher, err = t3.ToString()
	req.Nil(err)
	req.EqualValues("test ON MATCH SET test.ttt = 1", cypher)

	//merge with on create and on match
	cypher, err = t4.ToString()
	req.Nil(err)
	req.EqualValues("test ON CREATE SET test.tt1 = 2 ON MATCH SET test.ttt = 1", cypher)

	//error - path not defined
	_, err = t5.ToString()
	req.NotNil(err)

	//merge with on create and on match set to param string
	cypher, err = t6.ToString()
	req.Nil(err)
	req.EqualValues("test ON CREATE SET test = $props ON MATCH SET test = $props", cypher)

	//merge with on create and on match with members
	cypher, err = t7.ToString()
	req.Nil(err)
	req.Contains(cypher, "test ON CREATE SET test = $props")
	req.Contains(cypher, " ON MATCH SET ", "test.key1 = 1", "test.key2 = 'value2'", "test.key3 = $value3")
	req.Equal(102, len(cypher))
}
