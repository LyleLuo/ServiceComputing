import { PrimaryButton } from '@fluentui/react'
import * as React from 'react'
import { NavLink } from 'react-router-dom';
import AppContext from '../../AppContext';

const Home: React.FunctionComponent = () => {
  const { user, setUser } = React.useContext(AppContext);

  const login = () => {
    setUser(user ? undefined : {
      id: 1,
      name: 'user',
      email: 'someone@example.com'
    });
  };

  return <>
    <PrimaryButton text="hello" onClick={login} />
    <p>{user?.name}</p>
    <p>{user?.email}</p>
    <NavLink to="/tags">tags</NavLink>
  </>;
};

export default Home;