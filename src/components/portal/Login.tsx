import * as React from 'react';
import { PrimaryButton, Stack, Text, TextField } from '@fluentui/react';
import AppContext from '../../AppContext';

const Login: React.FunctionComponent = () => {
  const { setUser } = React.useContext(AppContext);
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();

  const login = () => {
    setUser({ id: 1, name: name!, email: name! + "@test.com" })
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