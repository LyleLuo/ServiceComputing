import * as React from "react";
import { useParams } from "react-router-dom";

interface DetailsRouteParam {
  id?: string
}

const Details: React.FunctionComponent = () => {
  const { id } = useParams<DetailsRouteParam>();
  return <>
    <p>文章 Id：{id}</p>
  </>;
};

export default Details;