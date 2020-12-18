import * as React from "react";
import UserInfo from "./models/UserInfo";

export interface Context {
  user?: UserInfo;
  setUser?: (userInfo?: UserInfo) => void;
  selectedKey?: string;
  setSelectedKey?: (key?: string) => void;
}

const AppContext = React.createContext<Context>({});

export default AppContext;