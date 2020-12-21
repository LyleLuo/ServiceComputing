import * as React from "react";
import ReactMarkdown from 'react-markdown'
import { useParams } from "react-router-dom";
import useHttp from "../../hooks/http";
import BlogContent from "../../models/BlogInfo";

interface DetailsRouteParam {
  id?: string;
}

const Details: React.FunctionComponent = () => {
  const { id } = useParams<DetailsRouteParam>();
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
    detailsRequest.fire();
  }, []);

  return <>
    <p>文章 Id：{id}</p>
    {
      content && <>
        <p>作者：{content.author}</p>
        <p>标题：{content.title}</p>
        <p>内容：</p>
        <ReactMarkdown>{content.text}</ReactMarkdown>
      </>
    }
  </>;
};

export default Details;