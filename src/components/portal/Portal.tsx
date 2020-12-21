import { PrimaryButton, Stack } from "@fluentui/react";
import * as React from "react";
import AppContext from "../../AppContext";
import { NavLink } from "react-router-dom";
import useHttp from "../../hooks/http";
import Login from "./Login";

interface ArticleModel {
  id: string;
  title: string;
}

const Portal: React.FunctionComponent = () => {
  const { user, setUser, setSelectedKey } = React.useContext(AppContext);
  const logoutRequest = useHttp<{ status: string; }>("/api/user/logout", "POST");
  const [list, setList] = React.useState<ArticleModel[]>();

  React.useEffect(() => {
    if (logoutRequest.data && !logoutRequest.loading) {
      if (logoutRequest.data.status === "success" && setUser) {
        setUser(undefined);
      }
    }
  }, [logoutRequest.data, logoutRequest.loading]);

  const logout = () => {
    logoutRequest.fire();
  };

  React.useEffect(() => {
    setSelectedKey && setSelectedKey("portal");
    fetch("/api/user/portal", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ author_id: user?.id })
    })
      .then(res => {
        if (!res.ok) {
          throw "failed to fetch";
        }
        return res.json();
      })
      .then(data => {

        console.log(data);
        setList(data.result);

      })
      .catch(err => {
        alert(err);
      });
  }, []);

  return user ? <>
    <p>欢迎，{user.name}！</p>
    <p>邮箱：{user.email}</p>
    <p>Id：{user.id}</p>
    <p>你所发布的全部博客：</p>
    <Stack>
      {
        list?.map((v, i) => {
          return <Stack.Item key={i} styles={{ root: { paddingTop: 10 } }}>
            <p>Blog_id: {v.id}</p>
            <p>Aitle: {v.title}</p>
            <PrimaryButton>
              <NavLink style={{ textDecoration: "none", color: "white" }} to={`/details/${v.id}`}>Go to details</NavLink>
            </PrimaryButton>
            <hr />
          </Stack.Item>;
        })
      }
    </Stack>
    <PrimaryButton text="退出" onClick={logout} />

  </> : <Login />;
};

export default Portal;
