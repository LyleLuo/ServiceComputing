import * as React from 'react';
import AppContext from '../../AppContext';
import { PrimaryButton, Stack, Text, TextField } from '@fluentui/react';

import axios from 'axios'
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';

const RegisterUrl = "http://localhost:8080/user/register"

const RegisterPost = (name:string, password:string, email:string) =>{
    console.log("username:%s",name);
    console.log("password:%s",password);
    console.log("email:%s",email);
    // let data = {"username":name, "password":password};
    let param = new URLSearchParams();
    param.append("username",name);
    param.append("password",password);
    param.append("email",email);
    
    axios({
      method:"post",
      url:RegisterUrl,
      data:param
    }).then(
      res=>{
       console.log("res=>",res);
      }
    );
  
    
    
    
  };

const Register: React.FunctionComponent = () => {
    const { setUser } = React.useContext(AppContext);
    const [name, setName] = React.useState<string>();
    const [password, setPassword] = React.useState<string>();
    const [email, setEmail] = React.useState<string>();

    const Register = () => {
        setUser({ id: -1, name: name!, email: email! })

        RegisterPost(name!,password!,email!);
    };

    return <Stack>
    <Stack.Item>
      <Text variant="xxLarge">注册账户</Text>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="用户名" defaultValue={name} onChange={(_, v) => setName(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="密码" canRevealPassword={true} type="password" defaultValue={password} onChange={(_, v) => setPassword(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="邮箱" canRevealPassword={true} type="email" defaultValue={email} onChange={(_, v) => setEmail(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <PrimaryButton text="注册" onClick={Register} />
    </Stack.Item>
  </Stack >;
};

export default Register;