import { ActionButton, PrimaryButton, Stack, Image } from "@fluentui/react";
import * as React from "react";
import { NavLink, useParams } from "react-router-dom";
import AppContext from "../../AppContext";
import useHttp from "../../hooks/http";

interface ArticleModel {
  title: string;
  username: string;
  blog_id: string;
}

interface ListModel {
  status: string;
  data: ArticleModel[];
}

interface PageRouteParams {
  id?: string;
}

const Page: React.FunctionComponent = () => {
  const { user } = React.useContext(AppContext);
  const { id } = useParams<PageRouteParams>();

  const [list, setList] = React.useState<ArticleModel[]>();
  React.useEffect(() => {
    fetch("/api/page/" + id, {
      method: "GET",
      credentials: "include"
    })
      .then(res => {
        if (!res.ok) {
          throw "failed to fetch";
        }
        return res.json();
      })
      .then(data => {
        if (data.status === "success") {
          console.log(data);
          setList(data.data);
        } else {
          alert("已到列表的尽头");
          window.history.back();
        }
      })
      .catch(err => {
        alert(err);
      });
  }, [id]);

  return list ? <>

    <Stack>
      {
        list.map((v, i) => {
          return <Stack.Item key={i} styles={{ root: { paddingTop: 10 } }}>
            <p>Title: {v.title}</p>
            <p>Author: {v.username}</p>
            <PrimaryButton>
              <NavLink style={{ textDecoration: "none", color: "white" }} to={`/details/${v.blog_id}`}>Go to details</NavLink>
            </PrimaryButton>
            <hr />
          </Stack.Item>;
        })
      }
    </Stack>
    <PrimaryButton>
      <NavLink style={{ textDecoration: "none", color: "white" }} to={`/page/${parseInt(id ?? "1") - 1}`}>上一页</NavLink>
    </PrimaryButton>
    &nbsp;&nbsp;&nbsp;&nbsp;
    <PrimaryButton>
      <NavLink style={{ textDecoration: "none", color: "white" }} to={`/page/${parseInt(id ?? "1") + 1}`}>下一页</NavLink>
    </PrimaryButton>

  </> : <p>loading</p>;
};

export default Page;