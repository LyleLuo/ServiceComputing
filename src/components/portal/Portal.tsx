import { PrimaryButton, ProgressIndicator, Stack } from "@fluentui/react";
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
        setList(data.result ?? []);

      })
      .catch(err => {
        alert(err);
      });
  }, [user]);

  return user ? <>
    <p style={{fontSize: 40}}>欢迎，{user.name}！</p>
    <p style={{fontSize: 20}}>你的邮箱：{user.email}</p>
    <p style={{fontSize: 20}}>你的id：{user.id}</p>
    <p style={{fontSize: 20}}>你的博客：</p>
    <Stack>
      {
        list?.map((v, i) => {
          return <Stack.Item key={i} styles={{ root: { paddingLeft: 20, height: 80, borderRadius: 20, marginBottom: 10, borderStyle: "groove",  backgroundColor : "cornsilk"} }}>
            <p style={{fontSize: 15}}>博客号: {v.id}</p>
            <p style={{fontWeight: "bold"}}>标题: {v.title}</p>
            <PrimaryButton style={{position: "relative", top: -55, right: 20, borderRadius: 10, float: "right"}}>
              <NavLink style={{ textDecoration: "none", color: "white" }} to={`/details/${v.id}`}>查看详情</NavLink>
            </PrimaryButton>
          </Stack.Item>;
        }) ?? <ProgressIndicator label="请稍等" description="正在加载文章列表" />
      }
    </Stack>
    <PrimaryButton text="退出" onClick={logout} />

  </> : <Login />;
};

export default Portal;
