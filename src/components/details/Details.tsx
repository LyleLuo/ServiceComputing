import * as React from "react";

export interface DetailsModel {
  title: string;
  content: string;
  arthor: string;
  time: Date;
  tags: string[]
}

const Details: React.FunctionComponent<DetailsModel> = (props) => {
  return <>
    <p>{props.title}</p>
  </>;
};

export default Details;