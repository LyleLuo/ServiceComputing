import * as React from 'react';
import AppContext from '../../AppContext';
import { PrimaryButton, Stack, Text, TextField } from '@fluentui/react';
import useHttp from '../../hooks/http';

const Register: React.FunctionComponent = () => {
  const { setUser } = React.useContext(AppContext);
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();
  const [email, setEmail] = React.useState<string>();
  const registerRequest = useHttp<{ status: string }>('/api/user/register', 'POST');

  React.useEffect(() => {
    if (!registerRequest.loading) {
      if (registerRequest.data?.status === 'success') {
        console.log('注册成功');
        setUser!({
          id: 1,
          name: name!,
          email: email!
        })
      }
    }
  }, [registerRequest.loading]);

  const Register = () => {
    registerRequest.fire({
      username: name,
      password: password,
      email: email
    });
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
      <TextField label="邮箱" type="email" defaultValue={email} onChange={(_, v) => setEmail(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <PrimaryButton text="注册" onClick={Register} />
    </Stack.Item>
  </Stack >;
};

export default Register;