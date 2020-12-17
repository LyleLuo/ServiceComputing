import { PrimaryButton } from '@fluentui/react';
import * as React from 'react';
import AppContext from '../../AppContext';
import useHttp from '../../hooks/http';
import Login from './Login';

const Portal: React.FunctionComponent = () => {
  const { user, setUser } = React.useContext(AppContext);
  const logoutRequest = useHttp<{ status: string }>('/api/user/logout', 'POST');

  React.useEffect(() => {
    if (logoutRequest.data && !logoutRequest.loading) {
      if (logoutRequest.data.status === "success") {
        setUser!(undefined);
      }
    }
  }, [logoutRequest.data, logoutRequest.loading])

  const logout = () => {
    logoutRequest.fire();
  }

  return user ? <>
    <p>欢迎，{user.name}！</p>
    <p>邮箱：{user.email}</p>
    <p>Id：{user.id}</p>
    <PrimaryButton text="退出" onClick={logout} />
  </> : <Login />;
};

export default Portal;
