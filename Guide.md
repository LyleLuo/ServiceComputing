# 前端开发指北

## 整体架构
项目采用 React 构建，UI 库使用 Fluent UI。

对于 UI 库的使用，可以参考 https://developer.microsoft.com/en-us/fluentui#/controls/web

## HTTP 请求
HTTP 请求封装在 `hooks/http.ts` 中，使用时只需要引入即可，例如：

```ts
import useHttp from "../hooks/http.ts";
```

`useHttp` 用法如下：

首先创建一个请求（注：因为做了 proxy 转发，所有的请求地址以 `/api` 开头，请求到后端会自动删除掉最前面的 `/api`）：

```ts
const myRequest = useHttp<响应类型>(请求相对地址, 请求方法);
// 例如
interface Response {
    status: string;
}
const myRequest = useHttp<Response>("/api/user/login", "POST");

// 或者不想声明一个 interface 也可以这样写：
const myRequest = useHttp<{ status: string }>("/api/user/login", "POST");
```

上述的例子表示：创建了一个到 `/user/login` 的 POST 请求，并且这个请求结束后返回的类型是 `Response` 类型的数据。

然后就可以发起请求了：

```ts
myRequest.fire({
    username: "abaabaaba",
    password: "1234567"
});
```

这样一来，请求就发了出去，那怎么样处理响应数据呢？

`myRequest` 中包含有 `data` 和 `loading` 两个成员，分别表示响应数据和是否正在加载。

我们可以通过 React Hooks 中的 `useEffect` 来处理，例如：

```ts
React.useEffect(() => {
    // 判断是否已经加载完了
    if (myRequest.loading) {
        console.log("加载中...");
    }
    // 判断是否有数据
    else if (myRequest.data) {
        console.log("请求完毕");
        console.log(myRequest.data);
    }
}, [myRequest.loading, myRequest.data]);
```

`useEffect` 像订阅一样，他的第二个参数是一个数组：`[myRequest.loading, myRequest.data]`，这表示，当这个数组里面的任何一个状态变化的时候就触发。当请求结束后，`loading` 会变成 `false`，`data` 会被设置，此时就可以触发回调函数，并拿 `data` 里面的数据了。

当然，如果你觉得我写的 `useHttp` 不好用的话，也可以直接用 `fetch` 函数做请求，例如：

```ts
// 一个 GET 请求的例子
fetch("/api/xxxx", {
      method: "GET",
      credentials: "include"
    })
      .then(res => res.json())
      .then(data => {
        console.log(data);
      });

// 一个 POST 请求的例子
fetch("/api/xxxx", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: "test", password: "123456" })
    })
      .then(res => res.json())
      .then(data => {
        console.log(data)
      });
```

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

注意，一个状态变化之后，任何引用了这个状态的 UI 会自动使用新的状态更新。例如如果我有：

```tsx
const XXX = () => {
  const [counter, setCounter] = React.useState(0);
  return <><p>现在数到了 {counter}</p><button onClick={() => setCounter(counter + 1)}>加 1</button></>;
}
```

如果点击按钮，那么 `onClick` 里会调用 `setCounter` 使得 `counter` 状态改变，界面也会自动随之改变。

这一点适用于状态、全局状态（见下文）、组件参数（见下文）和路由参数（见下文）。

### 全局状态
程序中可能需要一些全局状态，例如当前的用户信息等等，这些状态需要在多个组件中共享。

首先我们在 `AppContext.tsx` 里面定义全局状态所包含的东西。

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
    console.log("我被加载了");
    return /* B */ () => { console.log("我被卸载了") };
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
import * as React from "react";
import AppContext from "../../AppContext";
import { PrimaryButton, Stack, Text, TextField } from "@fluentui/react";
import useHttp from "../../hooks/http";

// 声明 Register 组件
const Register: React.FunctionComponent = () => {
  // 获取全局状态中的 setUser
  const { setUser } = React.useContext(AppContext);
  // 创建一系列的状态
  const [name, setName] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [email, setEmail] = React.useState("");
  // 创建一个 HTTP 请求
  const registerRequest = useHttp<{ status: string }>("/api/user/register", "POST");

  // 当 registerRequest.loading 和 registerRequest.data 发生改变时触发
  React.useEffect(() => {
    if (!registerRequest.loading) {
      if (registerRequest.data?.status === "success" && setUser /* 这里判断一下确保 setUser 不是空的 */) {
        console.log("注册成功");
        // 注册成功了于是设置全局状态中的用户信息
        setUser({
          id: 1,
          name: name,
          email: email
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
import * as React from "react";

export interface MyProps {
    content: string
}

const MyComponent: React.FunctionComponent<MyProps> = (props) => {
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

另外，通过 `props.children` 可以渲染传进来的子组件：

```tsx
const MyComponent: React.FunctionComponent<MyProps> = (props) => {
    return <div>{props.children}</div>;
}
```

上面这个例子，在 `<div></div>` 中渲染了父组件传进来的子组件，比如父组件这样使用 `MyComponent`：

```tsx
<MyComponent>
  <p>hello</p>
</MyComponent>
``` 

那这个 `<p>hello</p>` 会作为 `children` 传给 `MyComponent`，于是 `MyComponent` 就产生了 `<div><p>hello</p></div>`。

## 获取路由参数
假设你定义了如下路由：

```tsx
<Route path="/details/:id">
  <Details />
</Route>
```

这个路由有一个叫做 `id` 的参数，那如何在 `Details` 组件中获取这个参数的值呢？答案是使用 `useParams`。

首先定义路由参数的类型：

```ts
interface DetailsRouteParam {
  id: string
}
```

然后就可以使用了：

```ts
const { id } = useParams<DetailsRouteParam>();
```

这样当访问 `/details/123` 的时候，`id` 将是 123。

另外，可以通过 `:id?` 表示一个可以不指定的路由参数，比如用 `/details/:id?`，则直接访问 `/details` 也能命中路由，只不过这个时候 `id` 就是 `undefined` 了，对应的你就需要这么声明路由参数的类型：

```ts
interface DetailsRouteParam {
  id?: string
}
```

## 列表展开
假设你有一个数组，要怎么样把这个数组里面的东西渲染成 UI 呢？

可以用 `map` 方法：`map` 的作用是依次迭代数组中的每个元素，传递给你指定的函数，然后把它变成函数返回的东西。`map` 的参数是一个函数，这个函数有两个参数：当前的元素和当前的下标。

```tsx
const list = ["a", "b", "c", "d"];

const Content = <>
  {
    list.map((v, i) => {
      return <p key={i}>第 {i} 个元素是：{v}</p>
    })
  }
</>;
```

要注意 `map` 返回的每个最外层元素需要带一个 `key`，通过这个 `key` 可以告诉 React 哪些数据变化了需要重新渲染，而没有变化的数据就不会更新 UI，这样可以减少不必要的计算，上面的例子用了数组下标 `i` 作为 `key`。

上述代码最后 `Content` 会变成：

```tsx
<>
  <p>第 0 个元素是：a</p>
  <p>第 1 个元素是：b</p>
  <p>第 2 个元素是：c</p>
  <p>第 3 个元素是：d</p>
</>
```

## 页面导航
如果想要通过一个按钮把切换当前的地址到其他地址要怎么做呢？可以用 `NavLink`，`to` 指定跳转到的地址：

单独使用：
```tsx
<NavLink to="/abc">转到 abc</NavLink>
```

上面的代码会产生一个超链接，点击之后会跳转到 `/abc`。

```tsx
import { NavLink } from "react-router-dom";

const XXX = () => {
  const { page } = useParams<>();
  return <PrimaryButton>
    // 这里的 style 是为了消除下划线和字体的颜色
    <NavLink style={{ textDecoration: "none", color: "white" }} to="/abc">转到 abc</NavLink>
  </PrimaryButton>;
}
```

上面的代码会生成一个按钮，然后点击后会跳转到 `/abc`。

假设你有一个需要分页的页面，然后你的路由定义是 `/article/:page?`，也就是有一个 `page` 路由参数决定当前是哪一页（`/article/1` 就是第一页，`/article/2` 就是第二页，`:page?` 的 `?` 表示这个参数可以不指定，可以直接访问 `/article`，此时 `page` 将是 `undefined`）

那么你可以写：

```tsx
const XXX = () => {
  const { page } = useParams<{ page?: string }>();
  return <>
  <p>现在是第 {page} 页</p>
  <PrimaryButton>
    // `...` 括起来的字符串可以用 ${...} 嵌入表达式的值。
    <NavLink style={{ textDecoration: "none", color: "white" }} to={`/article/${(page ? parseInt(page) : 1) + 1}`}>下一页</NavLink>
  </PrimaryButton></>;
}
```

这样点击按钮之后就到下一页了。

另外，如果你用 `React.useEffect` 订阅 `page` 的更改的话，就可以做到翻页后加载新的数据了：

```tsx
// 新建一个状态用来保存当前文章列表
const [articles, setArticles] = React.useState<ArticleModel[]>();

React.useEffect(() => {
  // 如果 page 是空那默认第一页
  const currentPage = page ? page : "1";
  fetch(`/api/article/${currentPage}`, .....)
  .then(res => res.json())
  .then(data => {
    // 用请求回来的数据更新文章列表
    setArticles(data);
  })
}, [page])

return <>
  {
    articles.map((v, i) => <div key={i}><p>标题：{v.title} ......</p></div>)
  }
</>;
```
