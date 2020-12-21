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

interface HomeRouteParams {
  page?: string;
}

const Home: React.FunctionComponent = () => {
  const { user, setSelectedKey } = React.useContext(AppContext);
  const { page } = useParams<HomeRouteParams>();

  const [list, setList] = React.useState<ArticleModel[]>();
  React.useEffect(() => {
    setSelectedKey && setSelectedKey("home");
    fetch("/api/page/1", {
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
          alert("failed to load");
        }
      })
      .catch(err => {
        alert(err);
      });
  }, []);

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
      <NavLink style={{ textDecoration: "none", color: "white" }} to={`/page/${(page ? parseInt(page) : 1) + 1}`}>下一页</NavLink>
    </PrimaryButton>
  </> : <p>loading...</p>;
};

export default Home;