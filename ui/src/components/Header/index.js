import React from "react"
import { Link } from "react-router-dom"
import { inject } from "mobx-react"
import { Container, Image, Menu } from "semantic-ui-react"

import SearchBox from "./SearchBox"

@inject("commonStore")
class Header extends React.Component {
  render() {
    return (
      <div>
        <Menu fixed="top" inverted stackable>
          <Container>
            <Menu.Item as={Link} header to="/">
              <Image size="mini" src="/logo.png" style={{ marginRight: "1.5em" }} />
              {this.props.commonStore.appName}
            </Menu.Item>
            <Menu.Item position="right">
              <SearchBox />
            </Menu.Item>
          </Container>
        </Menu>
      </div>
    )
  }
}

export default Header
