import * as React from "react";
import AppContext from "../../AppContext";
import { PrimaryButton, Stack, Text, TextField } from "@fluentui/react";

const postUrl = "/api/user/post";

const Post: React.FunctionComponent = () => {
  const { user, setSelectedKey } = React.useContext(AppContext);
  const [title, setTitle] = React.useState<string>();
  const [currTags, setCurrTags] = React.useState<string>();
  const [tags, setTags] = React.useState<string[]>([]);
  const [text, setText] = React.useState<string>();
  const [type, setType] = React.useState<string>();

  React.useEffect(() => {
    setSelectedKey && setSelectedKey("post");
  }, [])

  React.useEffect(() => {
    // console.log("title:",title);
    console.log("tags:", tags);
    // console.log("text:",text);
  }, [tags]);

  const post = () => {
    fetch(postUrl, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title: title, author_id: user ? user.id : -1, tags: tags, text: text })
    }).then(res => res.json()).then(data => {
      console.log(data);
      setType(data["status"]);
    });
  };

  const BeforePost =
    <Stack>
      <Stack.Item>
        <Text variant="xxLarge">发布博客</Text>
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="标题" onChange={(_, v) => setTitle(v)} />
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 1000 } }}>
        <Stack horizontal>
          <Stack.Item>
            <TextField value={currTags} label="标签" onChange={(_, v) => setCurrTags(v)} />
          </Stack.Item>
          <Stack.Item styles={{ root: { paddingLeft: 10, paddingTop: 30 } }}>
            <PrimaryButton text="添加标签" onClick={() => {
              if (currTags) {
                setTags([...tags, currTags]);
                setCurrTags("");
              }
            }} />
          </Stack.Item>
          <Stack.Item styles={{ root: { paddingTop: 40, paddingLeft: 20 } }}>
            <i>已加标签：</i>
            {
              tags.map(
                (tag, index) => {
                  return <i key={index}>{" " + tag}</i>;
                }
              )
            }
          </Stack.Item>
        </Stack>
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 600, height: 400 } }}>
        <TextField label="正文" multiline rows={20} onChange={(_, v) => setText(v)} />
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <PrimaryButton text="发布" onClick={post} />
      </Stack.Item>
      {
        (type === "not login" || type === "failure") && <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
          <p style={{ color: "red" }}>{type}</p>
        </Stack.Item>
      }
    </Stack>;

  const AfterPost =
    <Stack>
      <Stack.Item styles={{ root: { paddingTop: 20 } }}>
        <Text variant="xLarge">发布成功</Text>
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <PrimaryButton text="再写一篇" onClick={() => {
          setType("initial");
          setCurrTags("");
          setTags([]);
        }} />
      </Stack.Item>
    </Stack>;


  return user ? type === "success" ? AfterPost : BeforePost : <p>请先登录</p>;
};

export default Post;
