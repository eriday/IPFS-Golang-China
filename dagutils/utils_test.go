
//此源码被清华学神尹成大魔王专业翻译分析并修改
//尹成QQ77025077
//尹成微信18510341407
//尹成所在QQ群721929980
//尹成邮箱 yinc13@mails.tsinghua.edu.cn
//尹成毕业于清华大学,微软区块链领域全球最有价值专家
//https://mvp.microsoft.com/zh-cn/PublicProfile/4033620
package dagutils

import (
	"context"
	"testing"

	path "gx/ipfs/QmNYPETsdAu2uQ1k9q9S1jYEGURaLHV6cbYRSVFVRftpF8/go-path"
	dag "gx/ipfs/QmTQdH4848iTVCJmKXYyRiK72HufWTLYQQ8iN3JaQ8K1Hq/go-merkledag"
	mdtest "gx/ipfs/QmTQdH4848iTVCJmKXYyRiK72HufWTLYQQ8iN3JaQ8K1Hq/go-merkledag/test"

	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	ipld "gx/ipfs/QmcKKBwfz6FyQdHR2jsXrrF6XeSBXYL86anmWNewpFpoF5/go-ipld-format"
)

func TestAddLink(t *testing.T) {
	ctx, context := context.WithCancel(context.Background())
	defer context()

	ds := mdtest.Mock()
	fishnode := dag.NodeWithData([]byte("fishcakes!"))

	err := ds.Add(ctx, fishnode)
	if err != nil {
		t.Fatal(err)
	}

	nd := new(dag.ProtoNode)
	nnode, err := addLink(ctx, ds, nd, "fish", fishnode)
	if err != nil {
		t.Fatal(err)
	}

	fnprime, err := nnode.GetLinkedNode(ctx, ds, "fish")
	if err != nil {
		t.Fatal(err)
	}

	fnpkey := fnprime.Cid()
	if !fnpkey.Equals(fishnode.Cid()) {
		t.Fatal("wrong child node found!")
	}
}

func assertNodeAtPath(t *testing.T, ds ipld.DAGService, root *dag.ProtoNode, pth string, exp cid.Cid) {
	parts := path.SplitList(pth)
	cur := root
	for _, e := range parts {
		nxt, err := cur.GetLinkedProtoNode(context.Background(), ds, e)
		if err != nil {
			t.Fatal(err)
		}

		cur = nxt
	}

	curc := cur.Cid()
	if !curc.Equals(exp) {
		t.Fatal("node not as expected at end of path")
	}
}

func TestInsertNode(t *testing.T) {
	root := new(dag.ProtoNode)
	e := NewDagEditor(root, nil)

	testInsert(t, e, "a", "anodefortesting", false, "")
	testInsert(t, e, "a/b", "data", false, "")
	testInsert(t, e, "a/b/c/d/e", "blah", false, "no link by that name")
	testInsert(t, e, "a/b/c/d/e", "foo", true, "")
	testInsert(t, e, "a/b/c/d/f", "baz", true, "")
	testInsert(t, e, "a/b/c/d/f", "bar", true, "")

	testInsert(t, e, "", "bar", true, "cannot create link with no name")
testInsert(t, e, "///“，”斜线“，真，”不能创建没有名称的链接“）

	c := e.GetNode().Cid()

	if c.String() != "QmZ8yeT9uD6ouJPNAYt62XffYuXBT6b4mP4obRSE9cJrSt" {
		t.Fatal("output was different than expected: ", c)
	}
}

func testInsert(t *testing.T, e *Editor, path, data string, create bool, experr string) {
	child := dag.NodeWithData([]byte(data))
	err := e.tmp.Add(context.Background(), child)
	if err != nil {
		t.Fatal(err)
	}

	var c func() *dag.ProtoNode
	if create {
		c = func() *dag.ProtoNode {
			return &dag.ProtoNode{}
		}
	}

	err = e.InsertNodeAtPath(context.Background(), path, child, c)
	if experr != "" {
		var got string
		if err != nil {
			got = err.Error()
		}
		if got != experr {
			t.Fatalf("expected '%s' but got '%s'", experr, got)
		}
		return
	}

	if err != nil {
		t.Fatal(err, path, data, create, experr)
	}

	assertNodeAtPath(t, e.tmp, e.root, path, child.Cid())
}
