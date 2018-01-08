import React, { Component } from "react"
import { Switch, Route, withRouter } from "react-router-dom"
import { inject, observer } from "mobx-react"

import { Container, Label } from "semantic-ui-react"
import Account from "../Account"

@inject("commonStore")
@withRouter
@observer
export default class Home extends Component {
  componentDidMount() {
    this.props.commonStore.loadStatus()
  }
  componentWillReceiveProps(nextProps) {
    this.props.commonStore.loadStatus()
  }
  render() {
    const { status } = this.props.commonStore
    return (
      <Container style={{ marginTop: "6em" }}>
        <Container style={{ marginBottom: "1em" }}>
          <Label>
            LAST BLOCK:
            <Label.Detail>{status.latest_block_height}</Label.Detail>
          </Label>
        </Container>
        <Switch>
          <Route path="/account/:address" component={Account} />
        </Switch>
      </Container>
    )
  }
}
