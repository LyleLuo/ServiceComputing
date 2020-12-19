import * as React from 'react';
import { PrimaryButton, Stack, Text, TextField } from '@fluentui/react';
import AppContext from '../../AppContext';
import Register from './Register';
import useHttp from '../../hooks/http';

const Login: React.FunctionComponent = () => {
  const { setUser } = React.useContext(AppContext);
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();
  const [type, setType] = React.useState<string>();
  const loginRequest = useHttp<{ status: string }>('/api/user/login', 'POST');

  React.useEffect(() => {
    if (!loginRequest.loading && type == "login") {
      if (loginRequest.data?.status === 'success') {
        setUser!({
          id: 1,
          name: name!,
          email: name + '@test.com'
        });
      }
      else {
        alert("密码或账号名错误");
      }
      setType("not defined");
    }
  }, [loginRequest.loading]);

  const login = () => {
    loginRequest.fire({
      username: name,
      password: password
    });
    setType("login");
  };

  const JumptoRegister = () => {
    setType("register");
  };

  const Content = <Stack>
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
      <Stack horizontal>
        <Stack.Item>
          <PrimaryButton text="登录" onClick={login} />
        </Stack.Item>
        <Stack.Item styles={{ root: { paddingLeft: 10 } }}>
          <PrimaryButton text="注册" onClick={(JumptoRegister)} />
        </Stack.Item>
      </Stack>
    </Stack.Item>
  </Stack >;


  if (type! == "register") {
    return <Register />;
  }
  else {
    return Content;
  }
};

export default Login;