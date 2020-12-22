
import { IStackTokens, Dropdown, DropdownMenuItemType, IDropdownStyles, PrimaryButton, DefaultButton, Stack, Checkbox, IDropdownOption } from "@fluentui/react";
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

  const isWant = (element:ArticleModelT):boolean=> {  //筛选函数
    var good=0;
    console.log("listorigin is",listorigin)
    console.log("list is",list)
    if(selectedKeys.length ==0){
      return true;
    }
    if(element.tags ==null){
      element.tags=["no tags"]
    }
    for (let i = 0; i <= selectedKeys.length-1; i++) {
      if(element.tags.indexOf(selectedKeys[i].toString())==-1){
        good = -1;
      }
    }
    console.log("good is",good)
    return good != -1;
  }
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
        placeholder="Select tags"
        label="Select Tags to search blog"
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
            <p>Title: {v.title}</p>
            <p>Author: {v.author}</p>
            <p>Tags:{` ${v.tags}`}</p>
            <PrimaryButton>
              <NavLink style={{ textDecoration: "none", color: "white" }} to={`/details/${v.id}`}>Go to details</NavLink>
            </PrimaryButton>
            <hr />
          </Stack.Item>;
        })
      }
    </Stack>
  </> : <p>loading...</p>;
};

export default Tags;