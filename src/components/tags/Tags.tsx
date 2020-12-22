
import { IStackTokens, Dropdown, DropdownMenuItemType, IDropdownStyles, PrimaryButton, DefaultButton, Stack, Checkbox, IDropdownOption, ProgressIndicator } from "@fluentui/react";
import * as React from "react";
import { NavLink, useParams } from "react-router-dom";
import AppContext from "../../AppContext";
import useHttp from "../../hooks/http";
// const options: IDropdownOption[] = tags?.map(v => ({ key: v.tagname, text: v.tagname }))
const stackTokens: IStackTokens = { childrenGap: 20 };
const dropdownStyles: Partial<IDropdownStyles> = {
  dropdown: { width: 300 },
};


interface ArticleModelT {
  id: string;
  author: string;
  title: string;
  tags: string[];
}

interface TagsModel {
  id: string;
  tagname: string;
}
const Tags: React.FunctionComponent = () => {
  const { setSelectedKey } = React.useContext(AppContext);
  const [selectedKeys, setSelectedKeys] = React.useState<string[]>([]);
  const [list, setList] = React.useState<ArticleModelT[]>();
  const [listorigin, setListO] = React.useState<ArticleModelT[]>();
  const [tags, setTags] = React.useState<TagsModel[]>();

  const isWant = (element: ArticleModelT): boolean => {  //筛选函数
    let good = 0;
    console.log("listorigin is", listorigin);
    console.log("list is", list);
    if (selectedKeys.length == 0) {
      return true;
    }
    if (element.tags == null) {
      element.tags = ["no tags"];
    }
    for (let i = 0; i <= selectedKeys.length - 1; i++) {
      if (element.tags.indexOf(selectedKeys[i].toString()) == -1) {
        good = -1;
      }
    }
    console.log("good is", good);
    return good != -1;
  };
  const onChange = (event: React.FormEvent<HTMLDivElement>, item?: IDropdownOption): void => {  //checkbox的函数
    if (item) {
      setSelectedKeys(
        item.selected ? [...selectedKeys, item.key as string] : selectedKeys.filter(key => key !== item.key),
      );
    }
    console.log(`The option has been changed to ${selectedKeys}.`);
  };

  React.useEffect(() => {
    setSelectedKey && setSelectedKey("tags");
    fetch("/api/user/tags", {
      method: "GET",
      credentials: "include"
    })
      .then(res => {
        if (!res.ok) {
          throw "failed to fetch";
        }
        return res.json();
      })
      .then(data => {
        if (data.status === "success") {
          console.log(data);
          setListO(data.blogs);
          setTags(data.tags);
        } else {
          alert("failed to load");
        }
      })
      .catch(err => {
        alert(err);
      });
  }, []);

  React.useEffect(() => {
    //setList((selectedKeys ?? []).length === 0 ? listorigin : listorigin?.filter(v => selectedKeys?.find(x => v.tags?.find(u => u === x) !== undefined) !== undefined));
    setList((selectedKeys ?? []).length === 0 ? listorigin : listorigin?.filter(isWant));
  }, [listorigin, selectedKeys]);

  const options: IDropdownOption[] = tags?.map(v => ({ key: v.tagname, text: v.tagname })) ?? [];

  return list ? <>

    <Stack tokens={stackTokens}>

      <Dropdown
        placeholder="选择标签"
        label="选择标签以筛选文章"
        selectedKeys={selectedKeys}
        // defaultSelectedKeys={['apple', 'banana', 'grape']}
        multiSelect
        options={options}
        styles={dropdownStyles}
        onChange={onChange}
      />
    </Stack>
    <hr />
    <Stack>
      {
        list.map((v, i) => {
          return <Stack.Item key={i} styles={{ root: { paddingTop: 0 } }}>
            <p>标题: {v.title}</p>
            <p>作者: {v.author}</p>
            <p>标签:{` ${v.tags ?? "无"}`}</p>
            <PrimaryButton>
              <NavLink style={{ textDecoration: "none", color: "white" }} to={`/details/${v.id}`}>查看详情</NavLink>
            </PrimaryButton>
            <hr />
          </Stack.Item>;
        })
      }
    </Stack>
  </> : <ProgressIndicator label="请稍等" description="正在加载标签列表" />;
};

export default Tags;