import { PrimaryButton } from '@fluentui/react';
import * as React from 'react';

const Tags: React.FunctionComponent = () => {
  const [test, setTest] = React.useState(1);

  React.useEffect(() => {
    alert('test 被改成 ' + test + ' 了');
  }, [test]);

  return <>
    <p>啊吧啊吧</p>
    <p>现在的 test 是：{test}</p>
    <PrimaryButton text="修改 test 的值" onClick={() => setTest(test + 1)}></PrimaryButton>
  </>;
};

export default Tags;