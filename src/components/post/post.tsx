import * as React from 'react';
import AppContext from '../../AppContext';
import { PrimaryButton, DefaultButton, Stack, Text, TextField } from '@fluentui/react';
import BlogInfo from "../../models/BlogInfo"
import useHttp from '../../hooks/http';

const postUrl = "/api/user/post";

const Post: React.FunctionComponent = () => {
  const { user } = React.useContext(AppContext);
  const [title, setTitle] = React.useState<string>();
  const [currTags, setCurrTags] = React.useState<string>();
  const [tags,setTags] = React.useState<string[]>([]);
  const [text, setText] = React.useState<string>();
  const [type, setType] = React.useState<string>();
  
  

  React.useEffect(() =>{
    // console.log("title:",title);
    console.log("tags:",tags);
    // console.log("text:",text);
  },[tags]);

  const post = ()=>{
    fetch(postUrl,{
      method:"POST",
      credentials:"include",
      headers:{ "Content-Type": "application/json" },
      body:JSON.stringify({title:title,author_id:user ? user.id : -1,tags:tags,text:text})
    }).then(res=>res.json()).then(data => {
      console.log(data);
      setType(data["status"])
    })
  }

  const BeforePost = 
  <Stack>
    <Stack.Item>
      <Text variant="xxLarge">发布博客</Text>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="标题" onChange={(_, v) => setTitle(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <Stack horizontal>
        <Stack.Item>
        <TextField label="标签" onChange={(_, v) => setCurrTags(v)} />
        </Stack.Item>
        <Stack.Item styles={{ root: { paddingLeft: 10, paddingTop: 30} }}>
          <PrimaryButton text="添加标签" onClick={() => {
            if(currTags){
              setTags([...tags,currTags!]);
            }
            }} />
        </Stack.Item>
      </Stack>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 600, height:400 } }}>
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
  </Stack>

  const AfterPost = <>发布成功</>
  

  return type === "success" ? AfterPost : BeforePost;
};

export default Post;