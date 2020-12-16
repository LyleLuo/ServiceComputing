import * as React from 'react';
import { PrimaryButton, Stack, Text, TextField } from '@fluentui/react';
import AppContext from '../../AppContext';
import axios from 'axios'
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';


const LoginUrl = "http://localhost:8080/user/login"

const axiosInstance = axios.create({
  baseURL: 'http://localhost:8080/user/',
  responseType: 'json',
  timeout: 50000,
  headers: {'X-Custom-Header': 'foobar'}
});

// const RequestType = '"get" | "GET" | "delete" | "DELETE" | "head" | "HEAD" | "options" | "OPTIONS" | "post" | "POST" | "put" | "PUT" | "patch" | "PATCH" | "purge" | "PURGE" | "link" | "LINK" | "unlink" | "UNLINK" | undefined'

export function fetch(url:string, params:any,type:string) {
  return new Promise((resolve, reject) => {
      axiosInstance({
        method: 'post',
        url,
        data:params,
        responseType: type?'blob':'json',
    }).then(res => {
       resolve(res)
    }).catch(err => {
       reject(err)
    })
  })
}


const LoginPost = (name:string, password:string) =>{
  console.log("username:%s",name);
  console.log("password:%s",password);
  // let data = {"username":name, "password":password};
  let param = new URLSearchParams();
  param.append("username",name);
  param.append("password",password);
  // axios.post(LoginUrl,data).then(
  //   res=>{
  //     console.log("res=>",res);
  //   }
  // ); 
  axios({
    method:"post",
    url:LoginUrl,
    data:param
  });

  // fetch("/login",data,"")
  
  
};

const Login: React.FunctionComponent = () => {
  const { setUser } = React.useContext(AppContext);
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();

  const login = () => {
    setUser({ id: 1, name: name!, email: name! + "@test.com" })
    LoginPost(name!,password!);
  };

  return <Stack>
    <Stack.Item>
      <Text variant="xxLarge">登录账户</Text>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="用户名" defaultValue={name} onChange={(_, v) => setName(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="密码" canRevealPassword={true} type="password" defaultValue={password} onChange={(_, v) => setPassword(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <PrimaryButton text="登录" onClick={login} />
    </Stack.Item>
  </Stack >;
};

export default Login;