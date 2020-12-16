import * as React from 'react';
import AppContext from '../../AppContext';
import Login from './Login';
import Register from './Register';

const Portal: React.FunctionComponent = () => {
  const { user } = React.useContext(AppContext);
  
  return <>
    { user ? <p>欢迎，{user.name}！</p> : <Login />}
  </>;
};

export default Portal;
