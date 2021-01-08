import { ActionButton, INavButtonProps, INavLinkGroup, isRelativeUrl, Nav, Stack, Text } from "@fluentui/react";
import * as React from "react";
import { NavLink } from "react-router-dom";
import AppContext from "../../AppContext";

const navItems: INavLinkGroup[] = [{
  links: [
    {
      name: "文章",
      url: "/",
      key: "home",
      iconProps: {
        iconName: "List"
      }
    },
    {
      name: "标签",
      url: "/tags",
      key: "tags",
      iconProps: {
        iconName: "Tag"
      }
    },
    {
      name: "我的",
      url: "/portal",
      key: "portal",
      iconProps: {
        iconName: "Contact"
      }
    },
    {
      name: "发布",
      url: "/post",
      key: "post",
      iconProps: {
        iconName: "Add"
      }
    }
  ]
}];

const Layout: React.FunctionComponent = (props) => {
  const { selectedKey, setSelectedKey } = React.useContext(AppContext);
  const renderLink = (linkProps: INavButtonProps) => {
    const link = linkProps.link ?? {
      url: "", name: undefined, forceAnchor: undefined, target: undefined, key: undefined, iconProps: undefined, icon: undefined, title: undefined, disabled: false, ariaLabel: undefined
    };
    const rel = link.url && link.target && !isRelativeUrl(link.url) ? "noopener noreferrer" : undefined;
    return <ActionButton as="div" className={linkProps.className} rel={rel} styles={linkProps.styles} iconProps={link.iconProps || { iconName: link.icon }} title={link.title !== undefined ? link.title : link.name} target={link.target} disabled={link.disabled} ariaLabel={link.ariaLabel ? link.ariaLabel : undefined}>
      <NavLink
        to={link.url || (link.forceAnchor ? "#" : "")}
        style={{ textDecoration: "none", outline: "transparent", color: "ButtonText", width: "100%", textAlign: "left", paddingLeft: 10 }}
        onClick={() => setSelectedKey && setSelectedKey(link.key)}>
        {link.name}
      </NavLink>
    </ActionButton >;
  };

  return <>
    <Stack horizontal={true}>
      <Stack.Item styles={{ root: { width: "20%", minWidth: 200 } }}>
        <Text variant="xxLarge" block nowrap style={{ textAlign: "center" }}>Service Blog</Text>
        <div className="nav-menu" style={{ paddingTop: 20 }}>
          <Nav linkAs={renderLink} groups={navItems} selectedKey={selectedKey} />
        </div>
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingLeft: 20, width: "80%" } }} align="stretch">
        {props.children}
      </Stack.Item>
    </Stack>

  </>;
};

export default Layout;