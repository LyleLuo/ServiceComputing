import { PrimaryButton } from "@fluentui/react";
import * as React from "react";
import { NavLink } from "react-router-dom";
import AppContext from "../../AppContext";

const Home: React.FunctionComponent = () => {
  const { user } = React.useContext(AppContext);

  return <>
    <PrimaryButton text="hello" />
    <p>{user?.name}</p>
    <p>{user?.email}</p>
    <NavLink to="/tags">tags</NavLink>
  </>;
};

export default Home;