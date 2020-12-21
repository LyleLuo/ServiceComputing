import { PrimaryButton, DefaultButton, Stack } from "@fluentui/react";
import * as React from "react";
import AppContext from "../../AppContext";
import useHttp from "../../hooks/http";

const Tags: React.FunctionComponent = () => {
  const { setSelectedKey } = React.useContext(AppContext);
  const [tag, setTag] = React.useState<string>();

  React.useEffect(() => {
    setSelectedKey && setSelectedKey("tags");
  }, []);

  const Content =
    <Stack>
      <Stack.Item>
        <DefaultButton text="全部" onClick={() => setTag("register")} />
      </Stack.Item>
    </Stack>;

  return Content;
};

export default Tags;