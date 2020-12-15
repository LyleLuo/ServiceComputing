import * as React from 'react';

const emptyFunction = () => { };

export interface UserInfo {
  id: number;
  name: string;
  email: string;
};

export interface Context {
  user?: UserInfo;
  setUser: (userInfo?: UserInfo) => void;
  selectedKey?: string;
  setSelectedKey: (key?: string) => void;
};

const AppContext = React.createContext<Context>({ setUser: emptyFunction, setSelectedKey: emptyFunction });

export default AppContext;