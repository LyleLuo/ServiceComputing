import * as React from 'react';
import { Route, BrowserRouter as Router, Switch } from 'react-router-dom';
import { IStackTokens, PrimaryButton, Stack, Dropdown, DropdownMenuItemType, IDropdownStyles, IDropdownOption } from '@fluentui/react';

const dropdownStyles: Partial<IDropdownStyles> = {
  dropdown: { width: 300 },
};

const options: IDropdownOption[] = [
  { key: 'fruitsHeader', text: 'Fruits', itemType: DropdownMenuItemType.Header },
  { key: 'apple', text: 'Apple' },
  { key: 'banana', text: 'Banana' },
  { key: 'orange', text: 'Orange', disabled: true },
  { key: 'grape', text: 'Grape' },
  { key: 'divider_1', text: '-', itemType: DropdownMenuItemType.Divider },
  { key: 'vegetablesHeader', text: 'Vegetables', itemType: DropdownMenuItemType.Header },
  { key: 'broccoli', text: 'Broccoli' },
  { key: 'carrot', text: 'Carrot' },
  { key: 'lettuce', text: 'Lettuce' },
];

const stackTokens: IStackTokens = { childrenGap: 20 };

const App: React.FunctionComponent = () => {
  return (
    <Router>
      <Switch>
        <Route exact path="/">
          <>
            <p>hello</p>
            <PrimaryButton text="hello" />
            <Stack tokens={stackTokens}>
              <Dropdown
                placeholder="Select an option"
                label="Basic uncontrolled example"
                options={options}
                styles={dropdownStyles}
              />

              <Dropdown
                label="Disabled example with defaultSelectedKey"
                defaultSelectedKey="broccoli"
                options={options}
                disabled={true}
                styles={dropdownStyles}
              />

              <Dropdown
                placeholder="Select options"
                label="Multi-select uncontrolled example"
                defaultSelectedKeys={['apple', 'banana', 'grape']}
                multiSelect
                options={options}
                styles={dropdownStyles}
              />
            </Stack>
          </>
        </Route>
      </Switch>
    </Router>
  );
};

export default App;