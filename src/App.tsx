import * as React from "react";
import { Route, BrowserRouter as Router, Switch } from "react-router-dom";
import AppContext from "./AppContext";
import Home from "./components/home/Home";
import Layout from "./components/layout/Layout";
import Portal from "./components/portal/Portal";
import Tags from "./components/tags/Tags";
import useHttp from "./hooks/http";
import UserInfo from "./models/UserInfo";

const App: React.FunctionComponent = () => {
  const [user, setUser] = React.useState<UserInfo>();
  const [selectedKey, setSelectedKey] = React.useState<string>();
  const userInfoRequest = useHttp<UserInfo>("/api/user/self", "GET");

  React.useEffect(() => {
    if (userInfoRequest.data && !userInfoRequest.loading) {
      setUser(userInfoRequest.data);
    }
  }, [userInfoRequest.data, userInfoRequest.loading]);

  React.useEffect(() => {
    userInfoRequest.fire();
  }, []);

  return (
    <AppContext.Provider value={{ user, setUser, selectedKey, setSelectedKey }}>
      <Router>
        <Layout>
          <Switch>
            <Route exact path="/">
              <Home />
            </Route>
            <Route path="/tags">
              <Tags />
            </Route>
            <Route path="/portal">
              <Portal />
            </Route>
          </Switch>
        </Layout>
      </Router>
    </AppContext.Provider>
  );
};

export default App;