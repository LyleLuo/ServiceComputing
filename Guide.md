# 前端开发指北

## 整体架构
项目采用 React 构建，UI 库使用 Fluent UI。

对于 UI 库的使用，可以参考 https://developer.microsoft.com/en-us/fluentui#/controls/web

## HTTP 请求
HTTP 请求封装在 `hooks/http.ts` 中，使用时只需要引入即可，例如：

```ts
import useHttp from '../hooks/http.ts';
```

`useHttp` 用法如下：

首先创建一个请求（注：因为做了 proxy 转发，所有的请求地址以 `/api` 开头，请求到后端会自动删除掉最前面的 `/api`）：

```ts
const myRequest = useHttp<响应类型>(请求相对地址, 请求方法);
// 例如
interface Response {
    status: string;
}
const myRequest = useHttp<Response>('/api/user/login', 'POST');

// 或者不想声明一个 interface 也可以这样写：
const myRequest = useHttp<{ status: string }>('/api/user/login', 'POST');
```

上述的例子表示：创建了一个到 `/user/login` 的 POST 请求，并且这个请求结束后返回的类型是 `Response` 类型的数据。

然后就可以发起请求了：

```ts
myRequest.fire({
    username: 'abaabaaba',
    password: '1234567'
});
```

这样一来，请求就发了出去，那怎么样处理响应数据呢？

`myRequest` 中包含有 `data` 和 `loading` 两个成员，分别表示响应数据和是否正在加载。

我们可以通过 React Hooks 中的 `useEffect` 来处理，例如：

```ts
React.useEffect(() => {
    // 判断是否已经加载完了
    if (myRequest.loading) {
        console.log('加载中...');
    }
    // 判断是否有数据
    else if (myRequest.data) {
        console.log('请求完毕');
        console.log(myRequest.data);
    }
}, [myRequest.loading, myRequest.data]);
```

`useEffect` 像订阅一样，他的第二个参数是一个数组：`[myRequest.loading, myRequest.data]`，这表示，当这个数组里面的任何一个状态变化的时候就触发。当请求结束后，`loading` 会变成 `false`，`data` 会被设置，此时就可以触发回调函数，并拿 `data` 里面的数据了。

## React Hooks
### useState 和 useEffect
页面中可能会有很多的状态，例如一个计数器页面，则需要一个 counter 变量保存当前数到的数字。

那么在 React 中怎么创建这样的状态呢？只需要：

```ts
const [counter, setCounter] = React.useState(0);
```

这表示我创建了一个初始状态是 0 的 `counter`，并且 `setCounter` 函数用来更新状态，比如调用 `setCounter(2)` 可以把 `counter` 更新成 2。

然后我就可以使用这个状态了，但是如果我想让状态改变时触发一个函数该怎么办呢？此时可以用 `React.useEffect`。

```ts
React.useEffect(() => {
    console.log(`counter 改变成了 ${counter}`);
}, [counter]);
```

`[counter]` 表示在 `counter` 变更的时候触发（如果有多个触发条件可以接着添加，因为这个参数是个数组，比如：`[state1, state2, state3...]`），于是当利用 `setCounter` 函数改变了 `counter` 时，就会触发定义的函数：

```ts
() => {
    console.log(`counter 改变成了 ${counter}`);
}
```
。

### 全局状态
程序中可能需要一些全局状态，例如当前的用户信息等等，这些状态需要在多个组件中共享。

首先我们在 `AppContext.tx` 里面定义全局状态所包含的东西。

例如：
```ts
export interface Context {
  user?: UserInfo;
  setUser?: (userInfo?: UserInfo) => void;
};
```

此时定义的 `Context` 里面有 2 个成员，分别是：用户信息、设置用户信息的函数，因为全局状态可能是空的，因此加了 `?` 表示可空。

然后我们在 `App.tsx` 中，创建了 `user` 这个状态，然后把这个状态提供出去即可：

```tsx
const App: React.FunctionComponent = () => {
  // 创建状态
  const [user, setUser] = React.useState<UserInfo>();

  return (
    // 提供状态
    <AppContext.Provider value={{ user, setUser }}>
        // 这里的东西将能拿到全局状态
    </AppContext.Provider>
  );
};
```

这样，我们在子组件中就能够使用定义的全局状态了：

```ts
const { user, setUser } = React.useContext(AppContext);
```

### 组件生命周期
如果想要让一个函数在它所在的组件加载或者卸载的时候执行怎么办？依然可以通过 `React.useEffect` 来实现：

```ts
React.useEffect(/* A */() => {
    console.log('我被加载了');
    return /* B */ () => { console.log('我被卸载了') };
}, []);
```

注意到此时不需要给 `React.useEffect` 的第二个参数的数组内放任何东西，只需要定义一个函数 A，然后这个函数 A 最后返回一个函数 B，那么函数 A 就会在组件加载的时候执行，B 就会在组件卸载的时候执行。

## 数据不变性
只要不是变量，一律使用 `const` 声明成常量，而不是用 `let`。

## 路由
路由的作用是，根据 URL 导航页面，路由定义在 `App.tsx` 中：

```tsx
<Router>
  <Layout>
    <Switch>
      <Route exact path="/">
        <Home />
      </Route>
      <Route path="/tags">
        <Tags />
      </Route>
      <Route path="/portal">
        <Portal />
      </Route>
    </Switch>
  </Layout>
</Router>
```

上述路由表示，当用户导航到地址 `/` 时（因为有 exact 所以必须是恰好为 `/`），向用户展示 `Home` 组件，而当地址前缀为 `/tags` 时则展示 `Tags` 组件，同理，当地址前缀为 `/portal` 时则展示 `Portal` 组件。

## 组件示例
以 Register 为例子：

```tsx
import * as React from 'react';
import AppContext from '../../AppContext';
import { PrimaryButton, Stack, Text, TextField } from '@fluentui/react';
import useHttp from '../../hooks/http';

// 声明 Register 组件
const Register: React.FunctionComponent = () => {
  // 获取全局状态中的 setUser
  const { setUser } = React.useContext(AppContext);
  // 创建一系列的状态
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();
  const [email, setEmail] = React.useState<string>();
  // 创建一个 HTTP 请求
  const registerRequest = useHttp<{ status: string }>('/api/user/register', 'POST');

  // 当 registerRequest.loading 和 registerRequest.data 发生改变时触发
  React.useEffect(() => {
    if (registerRequest.data && !registerRequest.loading) {
      if (registerRequest.data?.status === 'success') {
        console.log('注册成功');
        // 注册成功了于是设置全局状态中的用户信息
        setUser!({
        // 你可能会好奇这里的 '!' 是干什么用的：当你在访问一个可为空的东西的时候，如果你确定它一定不是空的，可以直接用 ! 来取消掉编译期的空安全检查，而不必在使用前通过类似 if (setUser) { setUser(...) } 判断到底是不是为空
          id: 1,
          name: name!,
          email: email!
        })
      }
    }
  }, [registerRequest.loading, registerRequest.data]);

  // 声明一个函数，用来发起注册的请求
  const Register = () => {
    // 发起请求
    registerRequest.fire({
      username: name,
      password: password,
      email: email
    });
  };

  // 你甚至可以把组件的一部分拿出来当作一个单独的组件
  const Button = <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <PrimaryButton text="注册" onClick={Register} />
    </Stack.Item>;

  // 每一个组件最后都要返回这个组件的内容布局
  return <Stack>
    <Stack.Item>
      <Text variant="xxLarge">注册账户</Text>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      // 这个编辑框默认值为 name 状态的值，如果用户输入发生改变了则调用 setName 更新 name 状态
      <TextField label="用户名" defaultValue={name} onChange={(_, v) => setName(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="密码" canRevealPassword={true} type="password" defaultValue={password} onChange={(_, v) => setPassword(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="邮箱" type="email" defaultValue={email} onChange={(_, v) => setEmail(v)} />
    </Stack.Item>
    // 这里使用上面定义的组件 Button
    { Button }
    // 甚至还能用条件来决定什么时候显示什么东西（null 表示空组件）
    {
      name.length > 3 ? Button : null
    }
  </Stack>;
};

// 最后把这个组件导出，以便其他地方可以引用这个组件
export default Register;
```

## React Fragment
注意到，每个组件最后 `return` 的时候只能返回一个组件，如果有两个组件需要返回的话那就需要把这两个组件用其他组件套起来，例如：

```tsx
// 不行
return <Component1 /><Component2 />
// 可以
return <div>
  <Component1 />
  <Component2 />
</div>
```

但如果不想让最上面有一层 `div` 的话，可以用 `<></>` 来表示这是个 `Fragment`，它可以用来嵌套组件，但是本身并不会产生一个 HTML 节点：

```tsx
return <>
  <Component1 />
  <Component2 />
</>;
```

## 组件参数
有的组件可能需要一些参数，例如一个显示文本的组件，需要提供显示的文本作为参数，那么可以通过 `props` 来做到：

```tsx
import * as React from 'react';

export interface MyProps {
    content: string
}

const MyComponent: React.FunctionComponent = (props: MyProps) => {
    return <div><p>{props.content}</p></div>;
}
```

这样的话，这个组件 `MyComponent` 需要提供一个参数 `props` 才能创建，使用的时候只需要提供 `content` 即可：

```tsx
<MyComponent content="12345"></MyComponent>
```

也可以传递状态，这样当父组件的状态变化之后，子组件也能感知到并更新界面：

```tsx
const [state, setState] = React.useState("hello");

return <MyComponent content={state}></MyComponent>
```

当 state 改变，这个改变也能自动传到 `MyComponent` 里面并自动更新界面，如果在 `MyComponent` 通过 `React.useEffect` 订阅了 `props.content` 的话，该状态改变的时候还能触发你定义的函数。

