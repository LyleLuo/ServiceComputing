import * as React from "react";
import ReactMarkdown from "react-markdown";
import gfm from "remark-gfm";
import { useParams } from "react-router-dom";
import AppContext from "../../AppContext";
import useHttp from "../../hooks/http";
import BlogContent from "../../models/BlogInfo";
import { ProgressIndicator } from "@fluentui/react";

interface DetailsRouteParam {
  id?: string;
}

const Details: React.FunctionComponent = () => {
  const { id } = useParams<DetailsRouteParam>();
  const { setSelectedKey } = React.useContext(AppContext);
  const [content, setContent] = React.useState<BlogContent>();
  const detailsRequest = useHttp<{ status: string, data: BlogContent; }>(`/api/details/${id}`, "GET");

  React.useEffect(() => {
    if (!detailsRequest.loading) {
      if (detailsRequest.data) {
        if (detailsRequest.data.status === "success") {
          setContent(detailsRequest.data.data);
        }
      }
      if (detailsRequest.error) {
        alert("加载失败");
      }
    }
  }, [detailsRequest.data, detailsRequest.ok, detailsRequest.loading]);

  React.useEffect(() => {
    setSelectedKey && setSelectedKey("home");
    detailsRequest.fire();
  }, []);

  return <>
    {
      content ? <>
        <h1>{content.title}</h1>
        <p>标签：{content.tags?.map((v, i) => <span key={i}>{v}&nbsp;</span>)}</p>
        <p>作者：{content.author}</p>
        <ReactMarkdown plugins={[gfm]}>{content.text}</ReactMarkdown>
      </> :
        <ProgressIndicator label="请稍等" description={`正在加载文章 ${id}`} />
    }
  </>;
};

export default Details;