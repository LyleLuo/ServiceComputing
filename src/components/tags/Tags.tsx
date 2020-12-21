import { PrimaryButton, DefaultButton, Stack } from "@fluentui/react";
import * as React from "react";
import useHttp from "../../hooks/http";

const Tags: React.FunctionComponent = () => {
  const [tag, setTag] = React.useState<string>();



  const Content =
    <Stack>
      <Stack.Item>
        <DefaultButton text="全部" onClick={() => setTag("register")} />
      </Stack.Item>
    </Stack>;

  return Content;
};

export default Tags;